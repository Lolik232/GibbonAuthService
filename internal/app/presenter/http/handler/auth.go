package handler

import (
	"auth-server/internal/app/service/services"
	"github.com/gorilla/mux"
	"net/http"
)

type AuthHandler struct {
	Handler
	serviceManager *services.Manager
}

func NewAuthHandler(serviceManager *services.Manager) *AuthHandler {
	return &AuthHandler{
		Handler:        Handler{},
		serviceManager: serviceManager,
	}
}

func (a *AuthHandler) ConfigureRoutes(router *mux.Router) {

	//auth := router.PathPrefix("/auth").Subrouter()

}

func (a *AuthHandler) authenticate() http.HandlerFunc {

}
