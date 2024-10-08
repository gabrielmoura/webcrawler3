package crawler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"sync"
	"time"
	"webCrawler3/config"
	"webCrawler3/infra/db"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/log"
	"webCrawler3/infra/pb"
)

var Wg sync.WaitGroup
var (
	mimeNotAllow        = errors.New("mime: not allowed")
	ErrUnexpectedStatus = errors.New("unexpected status")
)

// processPage processa uma página, extrai links e dados
func processPage(queueId pgtype.UUID, pageUrl string, depth int) {
	dataPageCtx, cancel := context.WithCancel(context.Background())
	log.Logger.Debug("Looping queue, depth: ", zap.Int("depth", depth))
	if depth > config.Conf.MaxDepth {
		log.Logger.Info("Reached max depth", zap.Int("depth", depth), zap.Int("maxDepth", config.Conf.MaxDepth))
		return
	}

	// Verifica se a página já foi visitada
	visited, err := repoDb.PageHasVisited(dataPageCtx, pageUrl)
	if err != nil {
		log.Logger.Error("Error checking if page has been visited", zap.Error(err))
		return
	}
	if visited {
		cancel()
		return
	}

	log.Logger.Info("Visiting", zap.String("url", pageUrl))
	plainText, htmlDoc, err := visitLink(pageUrl, queueId)
	if err != nil {
		if errors.Is(err, mimeNotAllow) {
			//log.Logger.Info(fmt.Sprintf("MIME not allowed: %s", pageUrl))
			return
		}
		log.Logger.Debug(fmt.Sprintf("Error checking link: %s", err))
		return
	}
	log.Logger.Info("visitado", zap.String("url", pageUrl))

	links, err := extractLinks(pageUrl, htmlDoc)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Error extracting links: %s", err))
		return
	}
	dataPage := &pb.Page{}
	err = extractData(htmlDoc, dataPage)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Error extracting data: %s", err))
		return
	}
	words, _ := countWordsInText(plainText)
	dataPage.Words = words
	dataPage.Url = pageUrl
	dataPage.Links = links
	dataPage.Timestamp = timestamppb.Now()
	dataPage.Visited = true

	if err := repoDb.SavePageAndMarkVisited(dataPageCtx, dataPage, queueId); err != nil {
		log.Logger.Error("Error saving page", zap.Error(err))
		return
	}

	handleAddToQueue(links, depth+1)
}
func HandleQueue() {

	log.Logger.Info("Starting queue handler")
	for {

		links, err := repoDb.RetrieveFromQueue(context.Background(), config.Conf.MaxConcurrency, 0) // Get a batch of links
		if err != nil {
			log.Logger.Error("Error retrieving links from queue", zap.Error(err))
			return
		}
		if len(links) == 0 {
			log.Logger.Info("No links in queue, sleeping for 5 seconds and retrying", zap.Any("links", links))
			time.Sleep(5 * time.Second)
			continue
		}

		for _, link := range links {
			if link.Url.String == "" {
				continue
			}
			Wg.Add(1) // Para cada link, incrementa o WaitGroup

			go func(link db.Queue) {
				defer Wg.Done()
				processPage(link.ID, link.Url.String, int(link.Depth.Int32))
			}(link)
		}
		Wg.Wait()
	}
}

func visitLink(pageUrl string, queueId pgtype.UUID) ([]byte, *html.Node, error) {
	resp, err := httpRequest(pageUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching URL %s: %w", pageUrl, err)
	}
	defer resp.Body.Close()

	if isStatusErr(resp.StatusCode, resp.Request.URL) {
		err := repoDb.MarkAsError(context.Background(), queueId)
		if err != nil {
			log.Logger.Error("Error marking as error", zap.Error(err))
		}
		log.Logger.Info("Status Error", zap.String("URL", pageUrl), zap.String("Status", resp.Status))
		return nil, nil, ErrUnexpectedStatus
	}

	// Streamlined MIME type check and early return
	if !isAllowedMIME(resp.Header.Get("Content-Type")) {
		return nil, nil, mimeNotAllow
	}

	// Efficiently read the response body into a buffer
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse HTML from the buffered content
	htmlDoc, err := html.Parse(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	return bodyBytes, htmlDoc, nil
}
