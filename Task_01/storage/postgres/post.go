package postgres

import (
	"context"
	"strconv"

	"database/sql"
	"fmt"
	"time"
	"user/models"
	"user/pkg/helper"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postRepo struct {
	db *pgxpool.Pool
}

func NewPostRepo(db *pgxpool.Pool) *postRepo {
	return &postRepo{
		db: db,
	}
}

func (u *postRepo) CreatePost(ctx context.Context, req *models.CreatePost) (string, error) {
	id := uuid.NewString()
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	query := `INSERT INTO posts(
	                        id,
							description,
							photos,
							created_by,
							created_at)
							VALUES($1, $2, $3,$4,now())`

	_, err := u.db.Exec(ctx, query,
		id,
		req.Description,
		req.Photos,
		TokenUser.UserID,
	)

	if err != nil {
		return "Error CreatePost", err
	}

	return id, nil
}

func (u postRepo) GetPost(ctx context.Context, req *models.IdRequest) (rep *models.Post, err error) {
	var (
		createdAt sql.NullTime
		createdBy sql.NullString
		likeCount sql.NullInt32
	)

	query := `
	SELECT 
             p.id, 
			 p.description,
			 p.photos,
			 p.created_at,
			 p.created_by,
    (SELECT COUNT(*) FROM post_likes l 
	        WHERE deleted_at IS NULL AND
			l.post_id = p.id) AS like_counts
FROM posts p 
WHERE
			p."deleted_at" IS NULL
			AND p."id" = $1
  `

	post := models.Post{}

	err = u.db.QueryRow(context.Background(), query, req.Id).Scan(
		&post.ID,
		&post.Description,
		&post.Photos,
		&createdAt,
		&createdBy,
		&likeCount,
	)
	if err != nil {
		return nil, fmt.Errorf("Post not found")
	}
	post.CreatedAt = createdAt.Time.Format(time.RFC3339)
	post.CreatedBy = createdBy.String
	if likeCount.Valid {
		post.LikeCount = strconv.Itoa(int(likeCount.Int32))
	}
	return &post, nil
}

func (u *postRepo) GetAllPost(ctx context.Context, req *models.GetAllPostRequest) (*models.GetAllPost, error) {
	params := make(map[string]interface{})
	filter := ` WHERE deleted_at IS NULL `
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	createdBy := sql.NullString{}
	likeCount := sql.NullInt32{}

	s := `
	SELECT 
		id,
		description,
		photos,
		created_at,
		created_by,
		(SELECT COUNT(*) FROM post_likes l
		WHERE 	l.post_id = p.id) AS like_count
	FROM posts p 
`

	if req.Search != "" {
		filter += ` AND description ILIKE '%' || :search || '%' `
		params["search"] = req.Search
	}

	limit := fmt.Sprintf(" LIMIT %d", req.Limit)
	offsetQ := fmt.Sprintf(" OFFSET %d", offset)
	query := s + " " + filter + " " + limit + " " + offsetQ

	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, q, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.GetAllPost{}
	count := 0
	for rows.Next() {
		var post models.Post
		count++
		err := rows.Scan(
			&post.ID,
			&post.Description,
			&post.Photos,
			&createdAt,
			&createdBy,
			&likeCount,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.Format(time.RFC3339)

		if createdBy.Valid {
			post.CreatedBy = createdBy.String
		}
		if likeCount.Valid {
			post.LikeCount = strconv.Itoa(int(likeCount.Int32))
		}
		resp.Posts = append(resp.Posts, post)
	}

	resp.Count = count
	return resp, nil
}

func (u *postRepo) GetMyPost(ctx context.Context, req *models.GetAllPostRequest) (*models.GetAllPost, error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	params := make(map[string]interface{})
	filter := fmt.Sprintf(`deleted_at IS NULL AND created_by = '%s'`, TokenUser.UserID)
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	likeCount := sql.NullInt32{}

	s := `
	SELECT 
		id,
		description,
		photos,
		created_at,
		created_by,
		(SELECT COUNT(*) FROM post_likes l
		 WHERE deleted_at is null and
		  l.post_id = p.id) AS like_count
		FROM posts p `

	limit := fmt.Sprintf(" LIMIT %d", req.Limit)
	offsetQ := fmt.Sprintf(" OFFSET %d", offset)
	query := s + " WHERE " + filter + limit + " " + offsetQ
	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, q, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.GetAllPost{}
	count := 0
	for rows.Next() {
		var post models.Post
		count++
		err := rows.Scan(
			&post.ID,
			&post.Description,
			&post.Photos,
			&createdAt,
			&TokenUser.UserID,
			&likeCount,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.Format(time.RFC3339)
		if likeCount.Valid {
			post.LikeCount = strconv.Itoa(int(likeCount.Int32))
		}
		resp.Posts = append(resp.Posts, post)
	}

	resp.Count = count
	return resp, nil
}

func (u *postRepo) GetAllDeletedPost(ctx context.Context, req *models.GetAllPostRequest) (*models.GetAllPost, error) {
	params := make(map[string]interface{})
	filter := ` `
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	createdBy := sql.NullString{}

	s := `
	SELECT 
		id,
		description,
		photos,
		created_at,
		created_by
	FROM
		posts
		where deleted_at IS  NOT NULL`

	if req.Search != "" {
		filter += ` AND description ILIKE '%' || :search || '%' `
		params["search"] = req.Search
	}

	limit := fmt.Sprintf(" LIMIT %d", req.Limit)
	offsetQ := fmt.Sprintf(" OFFSET %d", offset)
	query := s + filter + limit + offsetQ

	q, pArr := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, q, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.GetAllPost{}
	count := 0
	for rows.Next() {
		var post models.Post
		count++
		err := rows.Scan(
			&post.ID,
			&post.Description,
			&post.Photos,
			&createdAt,
			&createdBy,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt.Format(time.RFC3339)
		post.CreatedBy = createdBy.String

		resp.Posts = append(resp.Posts, post)
	}

	resp.Count = count
	return resp, nil
}

func (u *postRepo) UpdatePost(ctx context.Context, req *models.UpdatePost) (string, error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	query := `
			UPDATE "posts" 
				SET 
				"description" = $1,
				"photos" = $2,
				"updated_by"=$3,
				"updated_at" = NOW()
				WHERE "id" = $4 `

	result, err := u.db.Exec(
		context.Background(),
		query,
		req.Description,
		req.Photos,
		TokenUser.UserID,
		req.ID,
	)

	if err != nil {
		return "Error Update Post", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("post not found")
	}

	return req.ID, nil
}

func (b *postRepo) DeletePost(ctx context.Context, req *models.DeletePostRequest) (resp string, err error) {
	TokenUser := ctx.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "posts" 
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
