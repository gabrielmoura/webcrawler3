package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webCrawler3/infra/db/repoDb"
)

// HandlerSearch godoc
// @Summary Search for a URL in the database
// @Description Search for a URL in the database
// @Tags manager
// @url /api/v1/manager/search
func HandlerSearch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
	// search for a URL in the database
	case http.MethodGet:
		searchGet(w, r)
		// search for a URL in the database
	default:
		responseWithMethodNotAllowed(w)
	}
}
func searchGet(w http.ResponseWriter, r *http.Request) {
	limit := GetIntByQuery(r, "limit", 10)
	offset := GetIntByQuery(r, "offset", 0)
	searchTerm := GetStrByQuery(r, "search", "")
	if searchTerm == "" {
		responseWithBadRequest(w, fmt.Errorf("search term is required"))
		return
	}

	pages, err := repoDb.SearchByWeightWithPagination(r.Context(), searchTerm, int32(limit), int32(offset))
	if err != nil {
		responseWithInternalServerError(w, err)
		return
	}

	data, err := json.Marshal(pages)
	if err != nil {
		responseWithInternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
