package user_service

import (
	errors "auth-server/internal/app/errors/types"
	"auth-server/internal/app/model"
	"auth-server/internal/app/store"
	"auth-server/internal/app/utils"
	"context"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

type UserService struct {
	store       store.Store
	emailSender emailsender.IEmailSender
}

func New(emailSender emailsender.IEmailSender, store store.Store) (*UserService, error) {
	return &UserService{
		store:       store,
		emailSender: emailSender,
	}, nil
}

func (u UserService) FindUserByID(ctx context.Context, userID string, fields *store.UserFields) (*model.User, error) {

	usr, err := u.store.User().FindById(ctx, userID, fields)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (u UserService) FindUserByLogin(ctx context.Context, login string, fields *store.UserFields) (*model.User, error) {
	var usr *model.User
	var err error

	if strings.Contains(login, "@") {
		usr, err = u.store.User().FindByEmail(ctx, login, fields)
	} else {
		usr, err = u.store.User().FindByName(ctx, login, fields)
	}
	return usr, err
}

func (u UserService) FindUserByName(ctx context.Context, username string, fields *store.UserFields) (*model.User, error) {
	if len(username) > 0 {
		usr, err := u.store.User().FindByName(ctx, username, fields)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
	err := errors.ErrInvalidArgument.New("Username not be null!")
	return nil, err
}

func (u UserService) FindUserByEmail(ctx context.Context, email string, fields *store.UserFields) (*model.User, error) {
	if len(email) > 0 {
		usr, err := u.store.User().FindByEmail(ctx, email, fields)
		if err != nil {
			return nil, err
		}
		return usr, nil
	}
	err := errors.ErrInvalidArgument.New("email not be null!")
	return nil, err
}

func (u UserService) UpdateUserInfo(ctx context.Context, userID string, userinfo map[string]string) error {
	panic("implement me")
}

func (u UserService) Registration(ctx context.Context, user *model.User, password string) error {
	panic("implement me")
}

func (u UserService) ConfirmEmail(ctx context.Context, user *model.User, token string) error {
	key := viper.GetString("keys.email_conf_key")
	decodedID, tokenDeadTime, err := decodeEmailConfToken(token, key)
	if err != nil {
		switch errors.GetType(err) {
		case errors.ErrInvalidArgument:
			return err
		default:
			log.Printf("Err in conf email %s", err.Error())
			return errors.NoType.New("")
		}
	}
	if tokenDeadTime > time.Now().Unix() {
		return errors.ErrInvalidArgument.New("Token is dead.")
	}
	if decodedID != user.ID {
		return errors.ErrInvalidArgument.New("Invalid token.")
	}
	err = u.store.User().Update(ctx, user.ID, &model.User{EmailConfirmed: true})
	if err != nil {
		log.Printf("Err in conf email %s", err.Error())
	}
	return nil
}

func (u UserService) FindUserSessions(ctx context.Context, userID string) (*[]model.UserSession, error) {
	panic("implement me")
}

func (u UserService) Authenticate(ctx context.Context, login, password, clientID string) (*model.Identity, error) {
	panic("implement me")
}

func (u UserService) UpdateRefToken(ctx context.Context, userID, clientID, refToken string) (*model.Identity, error) {
	panic("implement me")
}

func (u UserService) SignOut(ctx context.Context, userID, sessionID string) error {
	panic("implement me")
}

func (u UserService) GenerateEmailConfToken(ctx context.Context, userID string) (string, error) {
	key := viper.GetString("keys.email_conf_key")
	token, err := generateEmailConfToken(userID, key)
	if err != nil {
		log.Printf("Error in generation token %s", err.Error())
	}
	return token, nil
}
