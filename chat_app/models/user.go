package models

type CreateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type CreateUserResp struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResp struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	AccessToken string
}
