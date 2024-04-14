package custom_errors

import "errors"

var (
	ErrAlreadyRegistered = errors.New("such a user is already registered")
	ErrBannerNotFound = errors.New("banner was not found")
	ErrNoTokenProvided   = errors.New("token was not provided")
	ErrTokenIsInvalid    = errors.New("invalid token provided")
)