package user_service

import (
	cfg "auth-server/internal/app/config"
	"auth-server/internal/app/model"
	"auth-server/internal/app/store"
	"auth-server/internal/app/utils/validators"
	errors "auth-server/pkg/errors/types"
	"context"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

type UserService struct {
	store         store.Store
	userValidator validators.IUserValidator
}

func (u *UserService) hashUserPassword(passBytes []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.MinCost)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "Err generate hash password.")
	}
	return bytes, nil
}
func (u *UserService) compareHashAndPassword(hash, pass []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, pass)
	return err
}

func New(store store.Store, uvalidator validators.IUserValidator) (*UserService, error) {
	us := UserService{
		store:         store,
		userValidator: uvalidator,
	}
	return &us, nil
}

func (u *UserService) FindUserByID(ctx context.Context, userID string, fields *store.UserFields) (*model.User, error) {
	usr, err := u.store.User().FindById(ctx, userID, fields)
	if err != nil {
		return nil, err
	}
	usr.Sanitize()
	return usr, nil
}

func (u *UserService) FindUserByLogin(ctx context.Context, login string, fields *store.UserFields) (*model.User, error) {
	var usr *model.User
	var err error

	if strings.Contains(login, "@") {
		usr, err = u.FindUserByEmail(ctx, login, fields)
	} else {
		usr, err = u.FindUserByName(ctx, login, fields)
	}
	if usr != nil {
		usr.Sanitize()
	}
	return usr, err
}

func (u *UserService) FindUserByName(ctx context.Context, username string, fields *store.UserFields) (*model.User, error) {
	if len(username) > 0 {
		usr, err := u.store.User().FindByName(ctx, username, fields)
		if err != nil {
			return nil, err
		}
		usr.Sanitize()
		return usr, nil
	}
	err := errors.ErrInvalidArgument.New("Username not be null!")
	return nil, err
}

func (u *UserService) FindUserByEmail(ctx context.Context, email string, fields *store.UserFields) (*model.User, error) {
	if len(email) == 0 || !strings.Contains(email, "@") {
		err := errors.ErrInvalidArgument.New("Email not be null!")
		return nil, err
	}
	usr, err := u.store.User().FindByEmail(ctx, email, fields)
	if err != nil {
		return nil, err
	}
	usr.Sanitize()
	return usr, nil
}

func (u UserService) UpdateUserInfo(ctx context.Context, userID string, userinfo map[string]string) error {
	panic("implement me")
}

func mapUserInfo(info map[string]string) map[string]string {
	userInfo := map[string]string{
		"first_name": info[store.UserInfoFirstName],
		"last_name":  info[store.UserInfoLastName],
		"mid_name":   info[store.UserInfoMidName],
	}
	return userInfo
}

func (u *UserService) Registration(ctx context.Context, user *model.User) (string, error) {
	userInfo := mapUserInfo(user.UserInfo)

	err := u.userValidator.Validate(ctx, u, user)
	if err != nil {
		switch errors.GetType(err) {
		case errors.ErrInvalidArgument:
			return "", err
		default:
			log.Printf("Err in registration user. Err: %s", err.Error())
			return "", errors.NoType.Newf("")
		}
	}
	passBytes := []byte(user.Password)
	passHash, err := u.hashUserPassword(passBytes)
	if err != nil {
		return "", err
	}
	user.PasswordHash = string(passHash)
	user.SanitizeForRegistration()

	user.UserInfo = userInfo

	id, err := u.store.User().Create(ctx, user)
	if err != nil {
		return "", err
	}
	token, err := u.GenerateEmailConfToken(ctx, id)
	if err != nil {
		log.Printf("Err in registration user. Err: %s", err.Error())
		err = u.store.User().DeleteById(ctx, id)
		if err != nil {
			log.Printf("Err in registration user. Err: %s", err.Error())
		}
		return "", errors.NoType.Newf("")
	}
	return token, nil
}
func (u *UserService) ConfirmEmail(ctx context.Context, userID, token string) error {
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
	if decodedID != userID {
		return errors.ErrInvalidArgument.New("Invalid token.")
	}
	err = u.store.User().Update(ctx, userID, &model.User{EmailConfirmed: true})
	if err != nil {
		log.Printf("Err in conf email %s", err.Error())
		return errors.NoType.Newf("")
	}
	return nil
}

//TODO: Implement FindUserSessions method
func (u *UserService) FindUserSessions(ctx context.Context, userID string) (*[]model.UserSession, error) {
	panic("implement me")
}

func (u *UserService) Authenticate(ctx context.Context, login, password, clientID string) (*model.User, *model.ClientRefToken, error) {
	fields := store.UserFields{
		UserName:         true,
		Email:            true,
		CreatedAt:        false,
		UserInfo:         false,
		UserSessions:     false,
		UserRoles:        false,
		UserPasswordHash: true,
	}
	user, err := u.FindUserByLogin(ctx, login, &fields)
	if err != nil {
		switch errors.GetType(err) {
		case errors.ErrInvalidArgument:
			return nil, nil, errors.ErrInvalidPasswordOrUsername.New("")
		}
		return nil, nil, err
	}
	err = u.compareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, errors.ErrInvalidPasswordOrUsername.New("")
	}
	user.Sanitize()
	panic("Implement me")
	//return user, nil, nil
}

func (u *UserService) UpdateRefToken(ctx context.Context, userID, clientID, refToken string) (*model.ClientRefToken, error) {
	panic("implement me")
}

func (u *UserService) SignOut(ctx context.Context, userID, sessionID string) error {
	panic("implement me")
}

func (u *UserService) GenerateEmailConfToken(ctx context.Context, userID string) (string, error) {
	key := cfg.Cfg.EmailConfKey
	token, err := generateEmailConfToken(userID, key)
	if err != nil {
		log.Printf("Error in generation token %s", err.Error())
		return "", err
	}
	return token, nil
}

func (u *UserService) DeleteById(ctx context.Context, userID string) error {
	err := u.store.User().DeleteById(ctx, userID)
	return err
}
func (u *UserService) DeleteByName(ctx context.Context, username string) error {
	err := u.store.User().DeleteByName(ctx, username)
	return err
}
