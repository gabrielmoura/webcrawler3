package handler

/**
HandlerQueue /api/v1/manager/queue - POST - Add a new URL to the queue
HandlerQueue /api/v1/manager/queue - GET - Get the list of URLs in the queue
HandlerSearch /api/v1/manager/search - POST - Search for a URL in the database
HandlerCrawl /api/v1/manager/crawl - POST - Start and Stop the crawling process
HandlerCrawl /api/v1/manager/crawl - GET - Get the status of the crawling process
*/
import (
	"net/http"
)

type Endpoint struct {
	Description string `json:"description"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Params      string `json:"params"`
}

// HandlerMaganer godoc
// @Summary Add a new URL to the queue
// @Description Add a new URL to the queue
// @Tags manager
// @Accept json
// @Produce json
// @url /api/v1/manager/queue
func HandlerMaganer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// return a list of endpoints

	default:
		responseWithMethodNotAllowed(w)
	}
}

// HandlerCrawl godoc
// @Summary Start and Stop the crawling process
// @Description Start and Stop the crawling process
// @Tags manager
// @url /api/v1/manager/crawl
func HandlerCrawl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// start and stop the crawling process
	case http.MethodGet:
	// get the status of the crawling process
	default:
		responseWithMethodNotAllowed(w)
	}
}
