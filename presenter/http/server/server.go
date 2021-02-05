package server

import (
	errors "auth-server/internal/app/errors/types"
	"auth-server/internal/app/service"
	"auth-server/internal/app/store"
	"context"
	"github.com/gorilla/mux"
)

type Server struct {
	router         *mux.Router
	serviceManager *service.Manager
	store          *store.Store
}

func newServer(ctx context.Context, sm *service.Manager, store *store.Store) (*Server, error) {
	if sm == nil {
		return nil, errors.ErrInvalidArgument.New("No service manager provided.")
	}
	if store == nil {
		return nil, errors.ErrInvalidArgument.New("No store provided.")
	}

}
