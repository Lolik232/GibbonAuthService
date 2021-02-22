package service

import (
	errors "auth-server/pkg/errors/types"
)

type Manager struct {
	User   UserService
	Client ClientService
}

func NewManager(u UserService, c ClientService) (*Manager, error) {
	if u == nil {
		return nil, errors.ErrInvalidArgument.New("No user service provided.")
	}
	//if c == nil {
	//	return nil, errors.ErrInvalidArgument.New("No client service provided.")
	//}
	return &Manager{
		User: u,
		//Client: c,
	}, nil
}
