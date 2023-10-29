package postgres

import (
	"context"

	"fmt"
	"user/models"
	"user/pkg/helper"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

type commentLikeRepo struct {
	db *pgxpool.Pool
}

func NewCommentLikeRepo(db *pgxpool.Pool) *commentLikeRepo {
	return &commentLikeRepo{
		db: db,
	}
}

func (u *commentLikeRepo) CreateCommentLike(ctx context.Context, req *models.CommentLikes) (string, error) {
	id := uuid.NewString()
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	var exists bool
	err := u.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)",
		req.Comment_Id,
		TokenUser.UserID).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("error checking for existing like: %w", err)
	}

	if exists {
		return "", fmt.Errorf("like already exists for the given comment and user")
	}

	query := `INSERT INTO comment_likes(
		id,
		comment_id,
		user_id,
		created_at)
		VALUES($1, $2, $3, now())`

	_, err = u.db.Exec(ctx, query,
		id,
		req.Comment_Id,
		TokenUser.UserID,
	)
	if err != nil {
		return "", fmt.Errorf("error creating comment like: %w", err)
	}

	return id, nil
}

func (u *commentLikeRepo) GetCommentLikes(ctx context.Context, req *models.CommentLikes) ([]models.CommentLike, error) {
	var likes []models.CommentLike

	query := `SELECT id, user_id, created_at 
			  FROM comment_likes 
			  WHERE comment_id = $1`

	rows, err := u.db.Query(ctx, query, req.Comment_Id)
	if err != nil {
		return nil, fmt.Errorf("error querying likes: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var like models.CommentLike

		err := rows.Scan(&like.ID, &like.UserId, &like.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning like: %w", err)
		}

		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting rows: %w", err)
	}

	if len(likes) == 0 {
		return nil, fmt.Errorf("no likes found for the specified comment_id")
	}

	return likes, nil
}

func (u *commentLikeRepo) DeleteCommentLike(ctx context.Context, req *models.DeleteCommentLikeRequest) (resp string, err error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	query := `
		UPDATE "comment_likes" 
		SET 
			"deleted_at" = NOW()
		WHERE 
			"user_id" = $1 AND "comment_id" = $2
	`
	result, err := u.db.Exec(
		context.Background(),
		query,
		TokenUser.UserID,
		req.Comment_Id,
	)
	if err != nil {
		return "Failed to delete like", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("like with ID %s not found", req.Comment_Id)
	}

	return req.Comment_Id, nil
}
