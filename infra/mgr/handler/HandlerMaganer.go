package handler

/**
HandlerQueue /api/v1/manager/queue - POST - Add a new URL to the queue
HandlerQueue /api/v1/manager/queue - GET - Get the list of URLs in the queue
HandlerSearch /api/v1/manager/search - POST - Search for a URL in the database
HandlerCrawl /api/v1/manager/crawl - POST - Start and Stop the crawling process
HandlerCrawl /api/v1/manager/crawl - GET - Get the status of the crawling process
*/
import (
	"encoding/json"
	"net/http"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/pb"
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
		respondWithMethodNotAllowed(w)
	}
}

// HandlerQueue godoc
// @Summary Get the list of URLs in the queue
// @Description Get the list of URLs in the queue
// @Tags manager
// @url /api/v1/manager/queue
func HandlerQueue(w http.ResponseWriter, r *http.Request) {
	// check if the request is a POST
	switch r.Method {
	case http.MethodPost:
		decode := json.NewDecoder(r.Body)
		var queue pb.QueueType
		err := decode.Decode(&queue)
		// add a new URL to the queue

		_, err = repoDb.AddToQueueIfNotExists(r.Context(), queue.Url, int(queue.Depth))

		if err == nil {
			msg := pb.StatusResponse{Message: "URL added to the queue"}
			data, _ := json.Marshal(msg)
			w.WriteHeader(http.StatusCreated)
			w.Write(data)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	case http.MethodGet:
	// get the list of URLs in the queue

	default:
		respondWithMethodNotAllowed(w)
	}
}

// HandlerSearch godoc
// @Summary Search for a URL in the database
// @Description Search for a URL in the database
// @Tags manager
// @url /api/v1/manager/search
func HandlerSearch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// search for a URL in the database
	default:
		respondWithMethodNotAllowed(w)
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
		respondWithMethodNotAllowed(w)
	}
}
