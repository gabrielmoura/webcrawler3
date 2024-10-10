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

	// Graph
	mux.HandleFunc("/api/v1/manager/graph", handler.HandlerGraph)
	//mux.HandleFunc("/api/v1/manager/graph", cacheMiddleware(10)(handler.HandlerGraph))

	mux.HandleFunc("/api/v1/manager", handler.HandlerMaganer)

	// load static files
	//fs := http.FileServer(http.Dir("./web"))
	//mux.Handle("/", fs)

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

//func cacheMiddleware(ttl int) func(http.HandlerFunc) http.HandlerFunc {
//	var cache = make(map[string][]byte)
//	return func(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			// Check if the request method is GET
//			if r.Method != http.MethodGet {
//				next(w, r)
//				return
//			}
//
//			// Check if the request URL is in the cache
//			if data, found := cache[r.RequestURI]; found {
//				w.Write(data)
//				return
//			}
//
//			// Call the next handler
//			next(w, r)
//		}
//	}
//}
