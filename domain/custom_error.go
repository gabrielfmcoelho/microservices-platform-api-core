package domain

// custom errors to represent the possible errors in the application

import (
	"errors"
)

var (
	ErrUserEmailNotFound     = errors.New("user email not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserWebUnauthorized   = errors.New("user unauthorized to login via web")
	ErrUserPasswordNotMatch  = errors.New("password does not match")
	ErrUnauthorized          = errors.New("unauthorized by the system")
	ErrNotFound              = errors.New("not found")
	ErrBadRequest            = errors.New("bad request")
	ErrInternalServerError   = errors.New("internal server error")
	ErrDataBaseInternalError = errors.New("database internal error")
	ErrInvalidIdentifier     = errors.New("invalid identifier (email or id)")
	ErrInvalidNumberToParse  = errors.New("invalid number to parse")
	ErrCategoryAlreadyExists = errors.New("category already exists")
)
