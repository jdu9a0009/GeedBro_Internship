package storage

import (
	"context"
	"time"
	"user/models"
)

type StorageI interface {
	User() UsersI
	Post() PostI
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
}

type PostI interface {
	CreatePost(context.Context, *models.CreatePost) (string, error)
	GetPost(context.Context, *models.IdRequest) (*models.Post, error)
	GetAllPost(context.Context, *models.GetAllPostRequest) (*models.GetAllPost, error)
	UpdatePost(context.Context, *models.Post) (string, error)
	DeletePost(context.Context, *models.IdRequest) (string, error)
}
type RedisI interface {
	Create(ctx context.Context, id string, obj interface{}, ttl time.Duration) error
	Get(ctx context.Context, id string, response interface{}) (bool, error)
	Delete(ctx context.Context, id string) error
}