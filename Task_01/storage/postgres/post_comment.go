package postgres

import (
	"context"

	"database/sql"
	"fmt"
	"time"
	"user/models"
	"user/pkg/helper"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postCommentRepo struct {
	db *pgxpool.Pool
}

func NewPostCommentRepo(db *pgxpool.Pool) *postCommentRepo {
	return &postCommentRepo{
		db: db,
	}
}

func (r *postCommentRepo) CreatePostComment(ctx context.Context, req *models.CreatePostComment) (string, error) {
	userInfo := ctx.Value("user_info").(helper.TokenInfo)
	id := uuid.NewString()

	query := `
		INSERT INTO "post_comments" (
			"id",
			"post_id",
			"comment",
			"created_by",
			"created_at"
		)
	VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(context.Background(), query,
		id,
		req.Post_id,
		req.Comment,
		userInfo.UserID,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create comment: %v", err)
	}

	return id, nil
}
func (u postCommentRepo) GetPostComment(ctx context.Context, req *models.IdRequest) (rep *models.PostComment, err error) {
	var (
		createdAt sql.NullTime
		createdBy sql.NullString
	)

	query := `SELECT 
                    id,
                    post_id,
                    comment,
                    created_at,
					created_by
				FROM
				    post_comments
            	WHERE
	               deleted_at IS  NULL
	               AND id = $1;`

	post := models.PostComment{}

	err = u.db.QueryRow(context.Background(), query, req.Id).Scan(
		&post.ID,
		&post.Post_id,
		&post.Comment,
		&createdAt,
		&createdBy,
	)
	if err != nil {
		return nil, fmt.Errorf("PostComment not found")
	}
	post.CreatedAt = createdAt.Time.Format(time.RFC3339)
	post.CreatedBy = createdBy.String

	return &post, nil
}

func (u *postCommentRepo) GetAllPostComment(ctx context.Context, req *models.GetAllPostCommentRequest) (*models.GetAllPostComment, error) {
	params := make(map[string]interface{})
	filter := ""
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	createdBy := sql.NullString{}

	s := `
	SELECT 
		id,
		post_id,
		comment,
		created_at,
		created_by
	FROM
		post_comments
	WHERE
		deleted_at IS  NULL
	`

	limit := fmt.Sprintf(" LIMIT %d", req.Limit)
	offsetQ := fmt.Sprintf(" OFFSET %d", offset)
	query := s + filter + limit + offsetQ

	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, q, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.GetAllPostComment{}
	count := 0
	for rows.Next() {
		var post models.PostComment
		count++
		err := rows.Scan(
			&post.ID,
			&post.Post_id,
			&post.Comment,
			&createdAt,
			&createdBy,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.Format(time.RFC3339)

		if createdBy.Valid {
			post.CreatedBy = createdBy.String
		}

		resp.Comments = append(resp.Comments, post)
	}

	resp.Count = count
	return resp, nil
}
func (u *postCommentRepo) GetMyPostComment(ctx context.Context, req *models.GetAllPostCommentRequest) (*models.GetAllPostComment, error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	params := make(map[string]interface{})
	filter := fmt.Sprintf(`deleted_at IS NULL AND created_by = '%s'`, TokenUser.UserID)
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	fmt.Println(req)
	s := `
	SELECT 
	       id,
	       post_id,
	       comment,
	       created_at,
	       created_by
	FROM
		post_comments`

	if req.Post_id != "" {
		filter += ` AND post_id = :search`
		params["search"] = req.Post_id
	}

	fmt.Println(s)

	limit := fmt.Sprintf(" LIMIT %d", req.Limit)
	offsetQ := fmt.Sprintf(" OFFSET %d", offset)
	query := s + " WHERE " + filter + limit + " " + offsetQ
	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, q, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.GetAllPostComment{}
	count := 0
	for rows.Next() {
		var post models.PostComment
		count++
		err := rows.Scan(
			&post.ID,
			&post.Post_id,
			&post.Comment,
			&createdAt,
			&TokenUser.UserID,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.Format(time.RFC3339)

		resp.Comments = append(resp.Comments, post)
	}

	resp.Count = count
	return resp, nil
}

func (u *postCommentRepo) GetAllDeletedPostComment(ctx context.Context, req *models.GetAllPostCommentRequest) (*models.GetAllPostComment, error) {
	params := make(map[string]interface{})
	filter := ` `
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	createdBy := sql.NullString{}

	s := `
	SELECT 
		id,
		post_id,
		comment,
		created_at,
		created_by
	FROM
		post_comments
		where deleted_at IS  NOT NULL`

	limit := fmt.Sprintf(" LIMIT %d", req.Limit)
	offsetQ := fmt.Sprintf(" OFFSET %d", offset)
	query := s + filter + limit + offsetQ

	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, q, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.GetAllPostComment{}
	count := 0
	for rows.Next() {
		var post models.PostComment
		count++
		err := rows.Scan(
			&post.ID,
			&post.Post_id,
			&post.Comment,
			&createdAt,
			&createdBy,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.CreatedBy = createdBy.String

		resp.Comments = append(resp.Comments, post)
	}

	resp.Count = count
	return resp, nil
}

func (u *postCommentRepo) UpdatePostComment(ctx context.Context, req *models.UpdatePostComment) (string, error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	query := `
			UPDATE "post_comments" 
				SET 
				"comment" = $1,
				"updated_by"=$2,
				"updated_at" = NOW()
				WHERE "id" = $3 and deleted_at is null `

	result, err := u.db.Exec(
		context.Background(),
		query,
		req.Comment,
		TokenUser.UserID,
		req.ID,
		TokenUser.UserID,
	)

	if err != nil {
		return "Error Update PostComment", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("post not found")
	}

	return req.ID, nil
}

func (b *postCommentRepo) DeletePostComment(ctx context.Context, req *models.DeletePostCommentRequest) (resp string, err error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "post_comments" 
		SET 
		     "deleted_at" = NOW(),
			 "deleted_by"=$1
		WHERE 
		     "deleted_at" IS  NULL AND
		     "id" = $2
`
	result, err := b.db.Exec(
		context.Background(),
		query,
		TokenUser.UserID,
		req.Id,
	)
	if err != nil {
		return "Failed to delete post", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("post with ID %s not found", req.Id)

	}

	return req.Id, nil
}
