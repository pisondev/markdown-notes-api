package web

type UserAuthRequest struct {
	Username string `validate:"required,min=3" json:"username"`
	Password string `validate:"required,min=6" json:"password"`
}
