package webrest

import (
	"encoding/json"
	"net/http"

	"github.com/colibri-project-io/colibri-sdk-go/pkg/base/logging"
)

type Error struct {
	Error string `json:"error"`
}

func JsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logging.Error(err.Error())
	}
}

func ErrorResponse(r *http.Request, w http.ResponseWriter, statusCode int, err error) {
	logging.Error("[%s] %s (%d): %v", r.Method, r.RequestURI, statusCode, err)
	JsonResponse(w, statusCode, Error{err.Error()})
}
