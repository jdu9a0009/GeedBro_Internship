package storage

import (
	"context"
	"user/models"
)

type StorageI interface {
	User() UsersI
	Message() MessageI
}

type UsersI interface {
	CreateUser(context.Context, *models.CreateUserReq) (*models.CreateUserResp, error)
	GetUserByEmail(context.Context, *models.LoginUserReq) (*models.User, error)
}
type MessageI interface {
	CreateMessage(context.Context, *models.Message) (string, error)
}
