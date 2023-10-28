package storage

import (
	"context"
	"time"
	"user/models"
)

type StorageI interface {
	User() UsersI
	Post() PostI
	PostComment() PostCommentI
	PostLike() PostLikeI
}
type CacheI interface {
	Cache() RedisI
}

type UsersI interface {
	CreateUser(context.Context, *models.CreateUser) (string, error)
	GetUser(context.Context, *models.IdRequest) (*models.User, error)
	GetAllUser(context.Context, *models.GetAllUserRequest) (*models.GetAllUser, error)
	UpdateUser(context.Context, *models.User) (string, error)
	DeleteUser(context.Context, *models.IdRequest) (string, error)
	GetByLogin(context.Context, *models.LoginRequest) (*models.LoginDataRespond, error)
	GetAllDeletedUser(context.Context, *models.GetAllUserRequest) (*models.GetAllUser, error)
}

type PostI interface {
	CreatePost(context.Context, *models.CreatePost) (string, error)
	GetPost(context.Context, *models.IdRequest) (*models.Post, error)
	GetAllPost(context.Context, *models.GetAllPostRequest) (*models.GetAllPost, error)
	UpdatePost(context.Context, *models.UpdatePost) (string, error)
	DeletePost(context.Context, *models.DeletePostRequest) (string, error)

	GetMyPost(ctx context.Context, req *models.GetAllPostRequest) (*models.GetAllPost, error)
	GetAllDeletedPost(context.Context, *models.GetAllPostRequest) (*models.GetAllPost, error)
}

type PostLikeI interface {
	CreatePostLike(context.Context, *models.PostLikes) (string, error)
	GetPostLikes(context.Context, *models.PostLikes) ([]models.PostLike, error)
	DeletePostLike(context.Context, *models.DeletePostLikeRequest) (string, error)
}

type PostCommentI interface {
	CreatePostComment(context.Context, *models.CreatePostComment) (string, error)
	GetPostComment(context.Context, *models.IdRequest) (*models.PostComment, error)
	GetAllPostComment(context.Context, *models.GetAllPostCommentRequest) (*models.GetAllPostComment, error)
	UpdatePostComment(context.Context, *models.UpdatePostComment) (string, error)
	DeletePostComment(context.Context, *models.DeletePostCommentRequest) (string, error)

	GetMyPostComment(ctx context.Context, req *models.GetAllPostCommentRequest) (*models.GetAllPostComment, error)
	GetAllDeletedPostComment(context.Context, *models.GetAllPostCommentRequest) (*models.GetAllPostComment, error)
}

type RedisI interface {
	Create(ctx context.Context, id string, obj interface{}, ttl time.Duration) error
	Get(ctx context.Context, id string, response interface{}) (bool, error)
	Delete(ctx context.Context, id string) error
}
