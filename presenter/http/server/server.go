package server

import (
	errors "auth-server/internal/app/errors/types"
	"auth-server/internal/app/service"
	"auth-server/internal/app/store"
	"auth-server/presenter/http/handler"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type server struct {
	server         *http.Server
	router         *mux.Router
	serviceManager *service.Manager
	store          *store.Store
}

func NewServer(ctx context.Context, sm *service.Manager, store *store.Store, handlers ...handler.IHandler) (*server, error) {
	if sm == nil {
		return nil, errors.ErrInvalidArgument.New("No service manager provided.")
	}
	if store == nil {
		return nil, errors.ErrInvalidArgument.New("No store provided.")
	}
	router := &mux.Router{}
	for _, h := range handlers {
		h.ConfigureRotes(router)
	}
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	return &server{
		server:         srv,
		router:         router,
		serviceManager: sm,
		store:          store,
	}, nil
}
func (s *server) Run() error {
	return s.server.ListenAndServe()
}

//func (s *server) Shutdown() error{
//	s.server.Shutdown(ctx)
//}
