package auth

import "errors"

var (
	ErrPasswordsNotEqual = errors.New("passwords are not equal")
	ErrBadCredentials    = errors.New("bad credentials")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoUsers           = errors.New("no users registered")
)
