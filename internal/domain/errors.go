package domain

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmptyToken           = errors.New("empty token")
	ErrNoMacFound           = errors.New("no mac found")
	ErrInvalidSystemContext = errors.New("invalid system context")
	ErrUnauthorized         = errors.New("unauthorized")
)
