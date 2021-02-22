package handler

import (
	"auth-server/internal/app/model"
	he "auth-server/pkg/errors/error"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type IHandler interface {
	ConfigureRoutes(router *mux.Router)
}

type Handler struct {
}

func (h Handler) respondJson(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
	return
}

func (h Handler) error(w http.ResponseWriter, r *http.Request, err error) {
	httpErr, code := he.New(err)
	resp := model.CreateBadResponce(httpErr)
	h.respondJson(w, r, code, resp)
	return
}

func (h Handler) respondHtml(w http.ResponseWriter, r *http.Request, data interface{}) {
	panic("Implement me.")
}
