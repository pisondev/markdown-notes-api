package exception

import "errors"

var (
	ErrConflictUser = errors.New("username already exist")
)
