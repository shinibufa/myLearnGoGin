package repository

type CommonResponse struct {
	Code int
	Msg  string
}

type LoginResponse struct {
	UserId int64 `json:"user_id"`
	Token string `json:"token"`
}