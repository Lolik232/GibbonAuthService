package store

import (
	"auth-server/internal/app/model"
	"context"
)

//UserRepository interface
type (
	UserRepository interface {
		//CRUD methods
		FindById(ctx context.Context, id string, params *UserFields) (*model.User, error)
		Update(ctx context.Context, userID string, user *model.User) error
		Create(ctx context.Context, user *model.User) (string, error)
		//
		UserSessionsFinder
		//
		FindUserClientRoles(ctx context.Context, userID, clientID string) ([]model.UserRole, error)
		//Utility functions
		UserAuthenticator
		//
		UserPassChecker
	}
	UserAuthenticator interface {
		Authenticate(ctx context.Context, login, passwordHash string) (*model.User, error)
		AuthenticateByName(ctx context.Context, username, passwordHash string) (*model.User, error)
		AuthenticateByEmail(ctx context.Context, email, passwordHash string) (*model.User, error)
		AuthenticateByID(ctx context.Context, userID, passwordHash string) (*model.User, error)
	}
	UserPassChecker interface {
		Authenticate(ctx context.Context, login, passwordHash string) (*model.User, error)
		AuthenticateByName(ctx context.Context, username, passwordHash string) (*model.User, error)
		AuthenticateByEmail(ctx context.Context, email, passwordHash string) (*model.User, error)
		AuthenticateByID(ctx context.Context, userID, passwordHash string) (*model.User, error)
	}
	UserSessionsFinder interface {
		FindSessions(ctx context.Context, id string) (*[]model.UserSession, error)
		CheckSession(ctx context.Context, id string) error
	}
)

type UserFields struct {
	UserName     bool `json:"username,omitempty"`
	Email        bool `json:"email,omitempty"`
	CreatedAt    bool `json:"created_at,omitempty"`
	UserInfo     bool `json:"user_info,omitempty"`
	UserSessions bool `json:"user_sessions,omitempty"`
	UserRoles    bool `json:"user_roles,omitempty"`
}

//ClientRepository interface
type ClientRepository interface {
	//CRUD methods
	FindById(ctx context.Context, id string) (*model.Client, error)
	FindRefToken(ctx context.Context, clientID, sessionID, refToken string) (*model.ClientRefToken, error)
	CheckRefToken(ctx context.Context, clientID, sessionID, refToken string) error
}
