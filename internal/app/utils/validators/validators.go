package validators

import (
	"auth-server/internal/app/model"
	"auth-server/internal/app/service"
	"auth-server/internal/app/store"
	errors "auth-server/pkg/errors/types"
	"context"
	"fmt"
	"regexp"
)

const (
	UsernameAllowedSymbols = "^[a-zA-Z0-9_-]{%d,%d}$"
	UsernameMinLength      = 4
	UsernameMaxLength      = 15
	UniqueEmail            = true
	UniqueUsername         = true
	PassMinLength          = 8
	PassRequiredDigits     = true
)

type (
	userValidatorConfiguration struct {
		UserNameAllowedSymbols *regexp.Regexp
		UsernameMinLength      int
		UsernameMaxLength      int
		PassAllowedSymbols     *regexp.Regexp
		PassRequiredSymbols    *regexp.Regexp
		UniqueEmail            bool
		UniqueUsername         bool
		PassMinLength          int
		PassRequiredDigits     bool
	}
	IUserValidator interface {
		Validate(ctx context.Context, service service.UserFinder, user *model.User) error
	}
)

type (
	UserValidator struct {
		params userValidatorConfiguration
	}
)

//New is constructor for UserValidator
func NewUserValidator() (*UserValidator, error) {
	usernameAllowedSymbols, err := regexp.Compile(fmt.Sprintf(UsernameAllowedSymbols, UsernameMinLength, UsernameMaxLength))
	if err != nil {
		return nil, err
	}
	//usernameRequiredSymbols, err := regexp.Compile(UsernameRequiredSymbols)
	//if err != nil {
	//	return nil, err
	//}
	//passRequiredSymbols, err := regexp.Compile()

	params := userValidatorConfiguration{
		UserNameAllowedSymbols: usernameAllowedSymbols,
		UsernameMinLength:      UsernameMinLength,
		UsernameMaxLength:      UsernameMaxLength,
		PassAllowedSymbols:     nil,
		PassRequiredSymbols:    nil,
		UniqueEmail:            UniqueEmail,
		UniqueUsername:         UniqueUsername,
		PassMinLength:          PassMinLength,
		PassRequiredDigits:     PassRequiredDigits,
	}

	return &UserValidator{
		params: params,
	}, nil
}
func (u UserValidator) Validate(ctx context.Context, service service.UserFinder, user *model.User) error {
	if UniqueEmail {
		u, err := (service).FindUserByEmail(ctx, user.Email, nil)
		if u != nil {
			return errors.ErrInvalidArgument.Newf("Error in validation. Email already taken.")
		}
		if err != nil {
			switch errors.GetType(err) {
			case errors.ErrInvalidArgument:
			default:
				return err
			}
		}
	}
	if UniqueUsername {
		u, err := (service).FindUserByName(ctx, user.UserName, nil)
		if u != nil {
			return errors.ErrInvalidArgument.Newf("Error in validation. Username already taken.")
		}
		if err != nil {
			switch errors.GetType(err) {
			case errors.ErrInvalidArgument:
			default:
				return err
			}
		}
	}
	if len(user.Password) < u.params.PassMinLength {
		return errors.ErrInvalidArgument.Newf("Error in validation. Password too short, min length is %d.", u.params.PassMinLength)
	}
	if len(user.UserName) < u.params.UsernameMinLength {
		return errors.ErrInvalidArgument.Newf("Error in validation. Username too short, min length is %d", u.params.UsernameMinLength)
	}
	if len(user.UserName) > u.params.UsernameMaxLength {
		return errors.ErrInvalidArgument.Newf("Error in validation. Username too long, max length is %d", u.params.UsernameMaxLength)
	}
	if ok := u.params.UserNameAllowedSymbols.MatchString(user.UserName); !ok {
		return errors.ErrInvalidArgument.Newf("Error in validation. Username must contains only \"A-Z,a-z,0-9,_,-\".")
	}
	pattern := "^[A-Za-z]{1,22}$"
	if ok, _ := regexp.MatchString(pattern, user.UserInfo[store.UserInfoFirstName]); !ok {
		return errors.ErrInvalidArgument.Newf("Error in validation. First name may contains only A-Z,a-z.")
	}
	if ok, _ := regexp.MatchString(pattern, user.UserInfo[store.UserInfoLastName]); !ok {
		return errors.ErrInvalidArgument.Newf("Error in validation. Last name may contains only A-Z,a-z.")
	}
	if ok, _ := regexp.MatchString(pattern, user.UserInfo[store.UserInfoMidName]); !ok {
		return errors.ErrInvalidArgument.Newf("Error in validation. Mid name may contains only A-Z,a-z.")
	}
	return nil
}
