package store

import (
	"auth-server/internal/app/model"
	"context"
)

//UserRepository interface
type (
	UserRepository interface {
		UserCrud
		UserSessionsFinder
		UserPassChecker
	}
	UserCrud interface {
		FindById(ctx context.Context, id string, params *UserFields) (*model.User, error)
		FindByName(ctx context.Context, name string, params *UserFields) (*model.User, error)
		FindByEmail(ctx context.Context, email string, params *UserFields) (*model.User, error)
		FindUserClientRoles(ctx context.Context, userID, clientID string) ([]model.UserRole, error)
		Update(ctx context.Context, userID string, user *model.User) error
		Create(ctx context.Context, user *model.User) (string, error)
	}
	UserPassChecker interface {
		CheckPassByID(ctx context.Context, userID, passwordHash string) error
		CheckPassByName(ctx context.Context, username, passwordHash string) error
		CheckPassByEmail(ctx context.Context, email, passwordHash string) error
	}
	UserSessionsFinder interface {
		FindSessions(ctx context.Context, id string) (*[]model.UserSession, error)
		CheckSession(ctx context.Context, id string) error
	}

	UserFields struct {
		UserName     bool `json:"username,omitempty"`
		Email        bool `json:"email,omitempty"`
		CreatedAt    bool `json:"created_at,omitempty"`
		UserInfo     bool `json:"user_info,omitempty"`
		UserSessions bool `json:"user_sessions,omitempty"`
		UserRoles    bool `json:"user_roles,omitempty"`
	}

	//ClientRepository interface
	ClientRepository interface {
		//CRUD methods
		FindById(ctx context.Context, id string) (*model.Client, error)
		FindRefToken(ctx context.Context, clientID, sessionID, refToken string) (*model.ClientRefToken, error)
		CheckRefToken(ctx context.Context, clientID, sessionID, refToken string) error
	}
)
