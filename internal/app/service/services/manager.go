package services

import (
	"auth-server/internal/app/service"
	"auth-server/internal/app/service/services/user_service"
	"auth-server/internal/app/store"
	"auth-server/internal/app/utils/validators"
	errors "auth-server/pkg/errors/types"
)

type Manager struct {
	User   service.UserService
	Client service.ClientService
}

//NewManager created a service manager and create services.
func NewManager(store store.Store, uv validators.IUserValidator) (*Manager, error) {
	if store == nil {
		return nil, errors.ErrInvalidArgument.New("Store is nill.")
	}
	if uv == nil {
		return nil, errors.ErrInvalidArgument.New("User validator is nill.")
	}
	//Create services
	userService, _ := user_service.New(store, uv)

	return &Manager{
		User: userService,
		//Client: c,
	}, nil
}
