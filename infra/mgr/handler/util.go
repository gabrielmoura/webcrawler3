package handler

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"webCrawler3/infra/log"
)

// responseWithMethodNotAllowed sends a 405 Method Not Allowed response.
func responseWithMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed"))
}
func responseWithBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	errStr := "Bad Request"
	if err != nil {
		errStr = err.Error()
	}
	w.Write([]byte(errStr))
}

func responseWithInternalServerError(w http.ResponseWriter, err error) {
	log.Logger.Error("Internal Server Error: ", zap.Error(err))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error"))
}
func responseWithNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found"))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetIntByQuery(r *http.Request, name string, defaultValue int) int {
	if queryValue := r.URL.Query().Get(name); queryValue != "" {
		if parsedValue, err := strconv.Atoi(queryValue); err == nil {
			return parsedValue
		}
	}
	return defaultValue
}
func GetStrByQuery(r *http.Request, name string, defaultValue string) string {
	if queryValue := r.URL.Query().Get(name); queryValue != "" {
		return queryValue
	}
	return defaultValue
}
