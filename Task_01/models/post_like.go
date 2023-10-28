package models

type PostLikes struct {
	Post_Id string `json:"id"`
}
type PostLike struct {
	ID        string `json:"id"`
	UserId    string `json:"user_id"`
	Post_Id   string `json:"post_id"`
	CreatedAt string `json:"created_at"`
}

type DeletePostLikeRequest struct {
	Id string `json:"id"`
}
