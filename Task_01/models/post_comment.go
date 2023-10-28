package models

type CreatePostComment struct {
	Post_id string `json:"post_id"`
	Comment string `json:"comment"`
}

type PostComment struct {
	ID        string `json:"id"`
	Post_id   string `json:"post_id"`
	Comment   string `json:"comment"`
	Like      string `json:"like"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type DeletePostCommentRequest struct {
	Id string `json:"id"`
}

type UpdatePostComment struct {
	ID      string `json:"id"`
	Post_id string `json:"post_id"`
	Comment string `json:"comment"`
}

type GetAllPostCommentRequest struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Post_id string `json:"post_id"`
}

type GetAllPostComment struct {
	Comments []PostComment `json:"PostComments"`
	Count    int           `json:"count"`
}
