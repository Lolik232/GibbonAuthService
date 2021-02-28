package server

import (
	"auth-server/internal/app/presenter/http/handler"
	"auth-server/internal/app/service/services"
	"auth-server/internal/app/store"
	errors "auth-server/pkg/errors/types"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

type server struct {
	server         *http.Server
	router         *mux.Router
	serviceManager *services.Manager
	store          store.Store
}

func NewServer(sm *services.Manager, store store.Store, handlers ...handler.IHandler) (*server, error) {
	if sm == nil {
		return nil, errors.ErrInvalidArgument.New("No service manager provided.")
	}
	if store == nil {
		return nil, errors.ErrInvalidArgument.New("No store provided.")
	}
	router := &mux.Router{}
	for _, h := range handlers {
		h.ConfigureRoutes(router)
	}
	port := os.Getenv("application_port")
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
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
