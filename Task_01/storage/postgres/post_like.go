package postgres

import (
	"context"

	"fmt"
	"user/models"
	"user/pkg/helper"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postLikeRepo struct {
	db *pgxpool.Pool
}

func NewPostLikeRepo(db *pgxpool.Pool) *postLikeRepo {
	return &postLikeRepo{
		db: db,
	}
}

func (u *postLikeRepo) CreatePostLike(ctx context.Context, req *models.PostLikes) (string, error) {
	id := uuid.NewString()
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	var exists bool
	err := u.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)",
		req.Post_Id,
		TokenUser.UserID).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("error checking for existing like: %w", err)
	}

	if !exists {
		query := `INSERT INTO post_likes(
			id,
			post_id,
			created_by,
			created_at)
			VALUES($1, $2, $3,now())`

		_, err := u.db.Exec(ctx, query,
			id,
			req.Post_Id,
			TokenUser.UserID,
		)

		if err != nil {
			return "Error CreatePostLike", err
		}
	}

	return id, nil
}

func (u *postLikeRepo) GetPostLikes(ctx context.Context, req *models.PostLikes) ([]models.PostLike, error) {

	var likes []models.PostLike

	query := `SELECT id, user_id, created_at 
			  FROM post_likes 
			  WHERE post_id = $1`

	rows, err := u.db.Query(ctx, query, req.Post_Id)
	if err != nil {
		return nil, fmt.Errorf("error querying likes: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var like models.PostLike

		err := rows.Scan(&like.ID, &like.UserId, &like.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning like: %w", err)
		}

		likes = append(likes, like)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting rows: %w", err)
	}

	return likes, nil
}
func (u *postLikeRepo) DeletePostLike(ctx context.Context, req *models.DeletePostLikeRequest) (resp string, err error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "post_likes" 
		SET 
		     "deleted_at" = NOW(),
			 "deleted_by"=$1
		WHERE 
		     "deleted_at" IS  NULL AND
		     "id" = $2
`
	result, err := u.db.Exec(
		context.Background(),
		query,
		TokenUser.UserID,
		req.Id,
	)
	if err != nil {
		return "Failed to delete like", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("like with ID %s not found", req.Id)

	}

	return req.Id, nil
}
