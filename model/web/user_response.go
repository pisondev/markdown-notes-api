package web

type UserRegisterResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
