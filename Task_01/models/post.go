package models

type CreatePost struct {
	UserId      string   `json:"created_by"`
	Description string   `json:"description"`
	Photos      []string `json:"photos" form:"photos"`
}

type Post struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	CreatedAt   string   `json:"created_at"`
	CreatedBy   string   `json:"created_by"`
	UpdatedAt   string   `json:"updated_at"`
	UpdatedBy   string   `json:"updated_by"`
	DeletedAt   string   `json:"deleted_at"`
	DeletedBy   string   `json:"deleted_by"`
}

type DeletePostRequest struct {
	DeletedBy string `json:"deleted_by"`
	Id        string `json:"id"`
}
type MyPostRequest struct {
	CreatedBy string `json:"created_by"`
}

type UpdatePost struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Photos      []string `json:"photos"`
	UpdatedBy   string   `json:"updated_by"`
}

type GetAllPostRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"description"`
}

type GetAllPost struct {
	Posts []Post `json:"Posts"`
	Count int    `json:"count"`
}
