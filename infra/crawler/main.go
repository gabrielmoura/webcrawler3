package crawler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
func processPage(queueId uuid.UUID, pageUrl string, depth, delay int) {
	dataPageCtx, cancel := context.WithCancel(context.Background())
	if depth > config.Conf.MaxDepth {
		log.Logger.Info("Reached max depth", zap.Int("depth", depth), zap.Int("maxDepth", config.Conf.MaxDepth))
		err := repoDb.MarkAsSkipped(context.Background(), queueId)
		if err != nil {
			log.Logger.Error("Error marking as skipped", zap.Error(err))
			return
		}
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
		err := repoDb.MarkAsSkipped(context.Background(), queueId)
		if err != nil {
			log.Logger.Error("Error marking as skipped", zap.Error(err))
		}
		return
	}

	// Caso haja um delay, aguarda o tempo especificado
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}

	log.Logger.Info("Visiting", zap.String("url", pageUrl))
	plainText, htmlDoc, err := visitLink(dataPageCtx, pageUrl, queueId)
	if err != nil {
		if errors.Is(err, mimeNotAllow) {
			//log.Logger.Info(fmt.Sprintf("MIME not allowed: %s", pageUrl))
			return
		}
		log.Logger.Debug(fmt.Sprintf("Error checking link: %s", err))
		err := repoDb.MarkAsError(context.Background(), queueId)
		if err != nil {
			log.Logger.Error("Error marking as error", zap.Error(err))
			return
		}
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

	// Salva a página e marca como visitada
	// Se houver um erro, marca item na fila como erro
	if err := repoDb.SavePageAndMarkVisited(dataPageCtx, dataPage, queueId); err != nil {
		log.Logger.Error("Error saving page", zap.Error(err), zap.String("url", pageUrl))
		err := repoDb.MarkAsError(context.Background(), queueId)
		if err != nil {
			return
		}
		return
	}

	handleAddToQueue(links, depth+1, delay)
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
			time.Sleep(config.Conf.GetCrlIdleTime())
			continue
		}

		for _, link := range links {
			if link.Url.String == "" {
				continue
			}
			Wg.Add(1) // Para cada link, incrementa o WaitGroup

			go func(link db.Queue) {
				defer Wg.Done()
				processPage(link.ID, link.Url.String, int(link.Depth.Int32), int(link.Delay.Int32))
			}(link)
		}
		Wg.Wait()
	}
}

func visitLink(ctx context.Context, pageUrl string, queueId uuid.UUID) ([]byte, *html.Node, error) {
	resp, err := httpRequest(ctx, pageUrl)
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
