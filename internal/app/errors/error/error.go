package error

import "auth-server/internal/app/errors/types"

type HTTPError struct {
	Code    uint              `json:"code"`
	Message string            `json:"err_msg"`
	Params  map[string]string `json:"params"`
}

var (
	codes = map[types.ErrorType]uint{
		types.NoType:                       500,
		types.ErrInvalidArgument:           1,
		types.ErrDatabaseDown:              500,
		types.ErrInvalidArgument:           500,
		types.ErrDuplicateEntry:            5,
		types.ErrInvalidPassword:           403,
		types.ErrInvalidPasswordOrUsername: 403,
	}
)

func New(err error) *HTTPError {
	errtype := types.GetType(err)
	return &HTTPError{
		Code:    codes[errtype],
		Message: "",
		Params:  nil,
	}
}
