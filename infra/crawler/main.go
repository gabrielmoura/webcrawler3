package crawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync"
	"time"
	"webCrawler3/config"
	"webCrawler3/infra/db"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/log"
)

var Wg sync.WaitGroup
var (
	mimeNotAllow        = errors.New("mime: not allowed")
	ErrUnexpectedStatus = errors.New("unexpected status")
)

func processPage(queueId uuid.UUID, pageUrl string, depth, delay int) {
	urlParsed, err := url.Parse(pageUrl)
	if err != nil {
		log.Logger.Error("Error parsing URL", zap.Error(err))
		return
	}
	switch urlParsed.Scheme {
	case "http", "https":
		processHttpScheme(queueId, pageUrl, depth, delay)
		return
	default:
		schemaNotAllow(urlParsed)
		return
	}
}
func schemaNotAllow(urlParsed *url.URL) {
	switch urlParsed.Scheme {
	case "ftp", "mailto", "tel", "irc", "xmpp", "magnet":
		log.Logger.Info(fmt.Sprintf("Schema %s encountered", strings.ToUpper(urlParsed.Scheme)), zap.String("Link", urlParsed.String()))
		return
	default:
		log.Logger.Error("Schema not allowed", zap.String("Link", urlParsed.String()), zap.Error(invalidSchemaErr))
	}

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
