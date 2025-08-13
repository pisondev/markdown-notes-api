package exception

import "errors"

var (
	ErrConflictUser      = errors.New("username already exist")
	ErrUnauthorizedLogin = errors.New("invalid username or password")
)
