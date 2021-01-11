package service

import (
	"auth-server/internal/app/model"
	"auth-server/internal/app/store"
	"context"
)

type (
	UserService interface {
		//Crud methods
		FindUserByID(ctx context.Context, userID string, fields *store.UserFields) (*model.User, error)
		//FindUserByName(ctx context.Context, username string, fields *store.UserFields) (*model.User, error)
		//FindUserByEmail(ctx context.Context, email string, fields *store.UserFields) (*model.User, error)
		//
		UserSessionsFinder
		//
		UpdateUserInfo(ctx context.Context, userID string, userinfo map[string]string) error
		//
		Registration(ctx context.Context, username, password string) error
		ConfirmEmail(ctx context.Context, user *model.User, token string) error
		//
		//
		UserAuthenticator
	}
	UserAuthenticator interface {
		Authenticate(ctx context.Context, login, password, clientID string) (*model.Identity, error)
		UpdateToken(ctx context.Context, userID, clientID, refToken string) (*model.Identity, error)
		SignOut(ctx context.Context, userID, sessionID string) error
	}
	UserSessionsFinder interface {
		FindUserSessions(ctx context.Context, userID string) (*[]model.UserSession, error)
	}
	//
	ClientService interface {
		FindClientByID(ctx context.Context, clientID string) (*model.Client, error)
	}
)
