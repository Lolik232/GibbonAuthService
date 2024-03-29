package store

import (
	"auth-server/internal/app/model"
	"context"
)

const (
	ParamUserName         = "username"
	ParamEmail            = "email"
	ParamCreatedAt        = "created_at"
	ParamUserInfo         = "user_info"
	ParamUserSessions     = "user_sessions"
	ParamUserRoles        = "user_roles"
	ParamUserPasswordHash = "password_hash"
)

const (
	UserInfoFirstName = "first_name"
	UserInfoLastName  = "last_name"
	UserInfoMidName   = "mid_name"
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
		DeleteById(ctx context.Context, userID string) error
		DeleteByName(ctx context.Context, userID string) error
		CreateSession(ctx context.Context, userID string) (string, error)
		DeleteSession(ctx context.Context, userID, sessionID string) error
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
		UserName         bool `json:"username,omitempty"`
		Email            bool `json:"email,omitempty"`
		CreatedAt        bool `json:"created_at,omitempty"`
		UserInfo         bool `json:"user_info,omitempty"`
		UserSessions     bool `json:"user_sessions,omitempty"`
		UserRoles        bool `json:"user_roles,omitempty"`
		UserPasswordHash bool `json:"-"`
	}

	//ClientRepository interface
	ClientRepository interface {
		//CRUD methods
		FindById(ctx context.Context, id string) (*model.Client, error)
		CreateRefToken(ctx context.Context, clientID string, refToken *model.ClientRefToken) error
		FindRefToken(ctx context.Context, clientID, sessionID, refToken string) (*model.ClientRefToken, error)
		CheckRefToken(ctx context.Context, clientID, sessionID, refToken string) (bool, error)
		DeleteRefToken(ctx context.Context, clientID, sessionID, refToken string) error
	}
)
