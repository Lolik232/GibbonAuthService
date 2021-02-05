package handler

import (
	he "auth-server/internal/app/errors/error"
	"encoding/json"
	"net/http"
)

type Handler struct {
}

func (h Handler) respondJson(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (h Handler) error(w http.ResponseWriter, r *http.Request, err error) {
	httpErr, code := he.New(err)
	h.respondJson(w, r, code, httpErr)
}

func (h Handler) respondHtml(w http.ResponseWriter, r *http.Request, data interface{}) {
	panic("Implement me.")
}
