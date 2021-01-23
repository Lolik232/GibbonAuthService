package error

import (
	"auth-server/internal/app/errors/types"
	"net/http"
)

type HTTPError struct {
	Code    uint              `json:"code"`
	Message string            `json:"err_msg"`
	Params  map[string]string `json:"params"`
}

var (
	codes = map[types.ErrorType]uint{
		types.NoType:                       http.StatusInternalServerError,
		types.ErrDatabaseDown:              http.StatusInternalServerError,
		types.ErrInvalidArgument:           1,
		types.ErrDuplicateEntry:            5,
		types.ErrInvalidPassword:           http.StatusUnauthorized,
		types.ErrInvalidPasswordOrUsername: http.StatusUnauthorized,
	}
)

func New(err error) *HTTPError {
	if err == nil {
		return &HTTPError{
			Code:    codes[types.NoType],
			Message: "Internal server error.",
			Params:  nil,
		}
	}
	errtype := types.GetType(err)
	msg := ""
	switch errtype {
	//Users should not be aware of internal problems
	case types.NoType, types.ErrDatabaseDown:
		msg = "Internal server error,"
	default:
		msg = err.Error()
	}
	return &HTTPError{
		Code:    codes[errtype],
		Message: msg,
		Params:  nil,
	}
}
