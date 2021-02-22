package types

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	NoType = ErrorType(iota)
	ErrInvalidArgument
	ErrInvalidPasswordOrUsername
	ErrInvalidPassword
	ErrDatabaseDown
	ErrDuplicateEntry
)

type ErrorType uint

var messages = map[ErrorType]string{
	ErrInvalidArgument:           "Invalid argument. ",
	ErrInvalidPassword:           "Invalid password. ",
	ErrInvalidPasswordOrUsername: "Invalid password or username. Check and try again. ",
	ErrDatabaseDown:              "Database down. ",
	ErrDuplicateEntry:            "Duplicate entry. ",
}

type customError struct {
	errorType     ErrorType
	standartError error
	contextInfo   errorContext
}

type errorContext struct {
	Field   string
	Message string
}

func (err customError) Error() string {
	return err.standartError.Error()
}

//New creates a no type error
func New(msg string) error {
	return customError{
		errorType:     NoType,
		standartError: errors.New(msg),
	}
}

//Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return customError{
		errorType:     NoType,
		standartError: fmt.Errorf(msg, args...),
	}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

//Wrap wraps an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

//Wrapf wraps an error with a format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			standartError: wrappedError,
			contextInfo:   customErr.contextInfo,
		}
	}
	return customError{
		errorType:     NoType,
		standartError: wrappedError,
	}
}

//AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{
		Field:   field,
		Message: message,
	}
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			standartError: customErr.standartError,
			contextInfo:   context,
		}
	}
	return customError{
		errorType:     NoType,
		standartError: err,
		contextInfo:   context,
	}
}

//GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok || customErr.contextInfo != emptyContext {
		return map[string]string{
			"field":   customErr.contextInfo.Field,
			"message": customErr.contextInfo.Message,
		}
	}
	return nil
}

//GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}
	return NoType
}

//New creates a new custom error
func (errorType ErrorType) New(msg string) error {
	return customError{
		errorType:     errorType,
		standartError: errors.New(messages[errorType] + msg),
	}
}

//Newf creates a new custom error with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(messages[errorType]+msg, args...)
	return customError{
		errorType:     errorType,
		standartError: err,
	}
}

//Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

//Wrapf creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return customError{
		errorType:     errorType,
		standartError: errors.Wrapf(err, messages[errorType]+msg, args...),
	}
}
