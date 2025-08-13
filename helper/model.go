package helper

import (
	"pisondev/markdown-notes-api/model/domain"
	"pisondev/markdown-notes-api/model/web"
)

func ToUserRegisterResponse(user domain.User) web.UserRegisterResponse {
	return web.UserRegisterResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}
