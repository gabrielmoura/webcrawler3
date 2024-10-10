package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/pb"
)

// HandlerQueue godoc
// @Summary Get the list of URLs in the queue
// @Description Get the list of URLs in the queue
// @Tags manager
// @url /api/v1/manager/queue
func HandlerQueue(w http.ResponseWriter, r *http.Request) {
	// check if the request is a POST
	switch r.Method {
	case http.MethodPost:
		addToQueueMethod(w, r)
	case http.MethodGet:
		getQueueMethod(w, r)
	case http.MethodPatch:
		updateQueueMethod(w, r)
	default:
		responseWithMethodNotAllowed(w)
	}
}
func updateQueueMethod(w http.ResponseWriter, r *http.Request) {
	//TODO: implement update queue method
	responseWithInternalServerError(w, fmt.Errorf("Not implemented"))
}
func addToQueueMethod(w http.ResponseWriter, r *http.Request) {
	decode := json.NewDecoder(r.Body)
	var queue pb.QueueType
	err := decode.Decode(&queue)
	// add a new URL to the queue
	// validate url , depth and delay
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	_, err = url.Parse(queue.Url)

	if err != nil {
		responseWithBadRequest(w, fmt.Errorf("Invalid URL: %s", err.Error()))
		return
	}

	if queue.Depth < 0 {
		responseWithBadRequest(w, fmt.Errorf("Depth must be greater than or equal to 0"))
		return
	}

	if queue.Delay < 0 {
		responseWithBadRequest(w, fmt.Errorf("Delay must be greater than or equal to 0"))
		return
	}
	_, err = repoDb.AddToQueueIfNotExists(r.Context(), queue.Url, int(queue.Depth), int(queue.Delay))

	if err == nil {
		msg := pb.StatusResponse{Message: "URL added to the queue"}
		data, _ := json.Marshal(msg)
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
		return
	}
	responseWithInternalServerError(w, err)
}

func getQueueMethod(w http.ResponseWriter, r *http.Request) {
	limit := GetIntByQuery(r, "limit", 10)
	offset := GetIntByQuery(r, "offset", 0)
	status := GetStrByQuery(r, "status", "pending")

	queues, err := repoDb.GetQueue(r.Context(), int32(limit), int32(offset), status)
	if err != nil {
		responseWithInternalServerError(w, err)
		return
	}

	data, _ := json.Marshal(queues)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
