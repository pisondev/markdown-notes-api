package service

import (
	"context"
	"pisondev/markdown-notes-api/model/web"
)

type UserService interface {
	Register(ctx context.Context, req web.UserAuthRequest) (web.UserRegisterResponse, error)
	Login(ctx context.Context, req web.UserAuthRequest) (web.UserLoginResponse, error)
}
