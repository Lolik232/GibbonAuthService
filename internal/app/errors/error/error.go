package error

import (
	"auth-server/internal/app/errors/types"
	"net/http"
)

type HTTPError struct {
	Code    uint              `json:"code,omitempty"`
	Message string            `json:"err_msg,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
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

func New(err error) (*HTTPError, int) {
	if err == nil {
		return &HTTPError{
			Code:    codes[types.NoType],
			Message: "Internal server error.",
			Params:  nil,
		}, http.StatusInternalServerError
	}
	errtype := types.GetType(err)
	msg := ""
	httpCode := 400

	switch errtype {
	//Users should not be aware of internal problems
	case types.NoType, types.ErrDatabaseDown:
		msg = "Internal server error,"
		httpCode = http.StatusConflict
	case types.ErrDuplicateEntry:
		msg = err.Error()
		httpCode = http.StatusConflict
	default:
		msg = err.Error()
	}

	return &HTTPError{
		Code:    codes[errtype],
		Message: msg,
		Params:  nil,
	}, httpCode
}
