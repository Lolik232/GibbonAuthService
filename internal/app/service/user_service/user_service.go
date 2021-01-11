package user_service

import (
	"auth-server/internal/app/model"
	"auth-server/internal/app/store"
	"context"
)

type UserService struct {
	store store.Store
}

func (u UserService) FindUserByID(ctx context.Context, userID string, fields *store.UserFields) (*model.User, error) {
	user, err := u.store.User().FindById(ctx, userID, fields)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserService) FindUserSessions(ctx context.Context, userID string) (*[]model.UserSession, error) {
	sessions, err := u.store.User().FindSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (u UserService) UpdateUserInfo(ctx context.Context, userID string, userinfo map[string]string) error {
	panic("implement me")
}

func (u UserService) Registration(ctx context.Context, username, password string) error {
	panic("implement me")
}

func (u UserService) ConfirmEmail(ctx context.Context, user *model.User, token string) error {
	panic("implement me")
}

func (u UserService) Authenticate(ctx context.Context, login, password, clientID string) (*model.Identity, error) {
	panic("implement me")
}

func (u UserService) UpdateToken(ctx context.Context, userID, clientID, refToken string) (*model.Identity, error) {
	panic("implement me")
}

func (u UserService) SignOut(ctx context.Context, userID, sessionID string) error {
	panic("implement me")
}
