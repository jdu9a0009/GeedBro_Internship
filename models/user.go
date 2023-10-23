package models

type CreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	IsActive  bool   `json:"isactive"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type IdRequest struct {
	Id string `json:"id"`
}

type LoginRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUp struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

type LoginRespond struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginRes struct {
	Token string `json:"token"`
}
type GetAllUserRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	UserName string `json:"username"`
}

type GetAllUser struct {
	Users []User `json:"user"`
	Count int    `json:"count"`
}
type ChangePassword struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}
type ReqNewPassword struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
