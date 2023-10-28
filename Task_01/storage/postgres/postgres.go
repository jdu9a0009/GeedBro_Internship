package postgres

import (
	"context"
	"fmt"
	"user/config"
	"user/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	db          *pgxpool.Pool
	users       *userRepo
	posts       *postRepo
	postLikes   *postLikeRepo
	postComment *postCommentRepo
}

func NewStorage(ctx context.Context, cfg config.Config) (storage.StorageI, error) {
	connect, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	))

	if err != nil {
		return nil, err
	}
	connect.MaxConns = cfg.PostgresMaxConnections

	pgxpool, err := pgxpool.ConnectConfig(ctx, connect)
	if err != nil {
		return nil, err
	}

	return &store{
		db: pgxpool,
	}, nil
}

func (b *store) User() storage.UsersI {
	if b.users == nil {
		b.users = NewUserRepo(b.db)
	}
	return b.users
}
func (b *store) Post() storage.PostI {
	if b.posts == nil {
		b.posts = NewPostRepo(b.db)
	}
	return b.posts
}

func (b *store) PostLike() storage.PostLikeI {
	if b.postLikes == nil {
		b.postLikes = NewPostLikeRepo(b.db)
	}
	return b.postLikes
}

func (b *store) PostComment() storage.PostCommentI {
	if b.postComment == nil {
		b.postComment = NewPostCommentRepo(b.db)
	}
	return b.PostComment()
}

func (s *store) Close() {
	s.db.Close()
}
