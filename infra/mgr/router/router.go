package router

import (
	"context"
	"net/http"
	"time"
	"webCrawler3/config"
	"webCrawler3/infra/mgr/handler"
)

/**
HandlerQueue /api/v1/manager/queue - POST - Add a new URL to the queue
HandlerQueue /api/v1/manager/queue - GET - Get the list of URLs in the queue
HandlerSearch /api/v1/manager/search - POST - Search for a URL in the database
HandlerCrawl /api/v1/manager/crawl - POST - Start and Stop the crawling process
HandlerCrawl /api/v1/manager/crawl - GET - Get the status of the crawling process
*/

func InitRouters(ctx context.Context, conf *config.Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/manager/queue", handler.HandlerQueue)
	mux.HandleFunc("/api/v1/manager/search", handler.HandlerSearch)
	mux.HandleFunc("/api/v1/manager/crawl", handler.HandlerCrawl)
	mux.HandleFunc("/api/v1/manager", handler.HandlerMaganer)

	s := &http.Server{
		Addr:           config.Conf.GetAddrManager(),
		Handler:        mux,
		ReadTimeout:    10 * time.Second, // Tempo máximo para leitura
		WriteTimeout:   10 * time.Second, // Tempo máximo para escrita
		MaxHeaderBytes: 1 << 10,          // 1 KB
	}

	go func() {
		<-(ctx).Done()
		s.Shutdown(context.Background())
	}()
	err := s.ListenAndServe()
	if err != nil {
		return
	}
}
