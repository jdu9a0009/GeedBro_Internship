package models

import "time"

type CommentLikes struct {
	Comment_Id string `json:"comment_id"`
}
type CommentLike struct {
	ID         string    `json:"id"`
	Comment_Id string    `json:"comment_id"`
	UserId     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type DeleteCommentLikeRequest struct {
	Comment_Id string `json:"post_id"`
}
