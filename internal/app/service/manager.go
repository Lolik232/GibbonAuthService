package service

import (
	errors "auth-server/internal/app/errors/types"
	"context"
)

type Manager struct {
	User   *UserService
	Client *ClientService
}

func NewManager(ctx context.Context, u *UserService, c *ClientService) (*Manager, error) {
	if u == nil {
		return nil, errors.ErrInvalidArgument.New("No user service provided.")
	}
	if c == nil {
		return nil, errors.ErrInvalidArgument.New("No client service provided.")
	}
	return &Manager{
		User:   u,
		Client: c,
	}, nil
}
