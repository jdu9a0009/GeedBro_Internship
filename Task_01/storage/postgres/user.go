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

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) CreateUser(ctx context.Context, req *models.CreateUser) (string, error) {
	id := uuid.NewString()

	query := `INSERT INTO users(
	                        id,
							username,
							password,
							created_at)
							VALUES($1, $2, $3,now())`

	_, err := u.db.Exec(ctx, query,
		id,
		req.Username,
		req.Password,
	)

	if err != nil {
		return "Error CreateUser", err
	}

	return id, nil
}

func (u userRepo) GetUser(ctx context.Context, req *models.IdRequest) (rep *models.User, err error) {
	query := `SELECT 
                    id,
                    username,
                    password,
                    is_active,
                    created_at,

				FROM
				    users
				WHERE
				    is_active = true AND
				    id = $1;`

	var (
		createdAt sql.NullTime
	)
	user := models.User{}

	err = u.db.QueryRow(context.Background(), query, req.Id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	user.CreatedAt = createdAt.Time.Format(time.RFC3339)
	return &user, nil
}

func (u *userRepo) GetAllUser(ctx context.Context, req *models.GetAllUserRequest) (resp *models.GetAllUser, err error) {
	params := make(map[string]interface{})
	filter := `   WHERE is_active `
	resp = &models.GetAllUser{}
	resp.Users = make([]models.User, 0)
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	updatedAt := sql.NullTime{}

	query := `SELECT 
	               COUNT(*) OVER(),
	               id,
				   username,
				   password,
				   is_active,
				   created_at,
				   updated_at
			FROM users`

	if req.UserName != "" {
		filter += ` AND username ILIKE '%' || :search || '%' `
		params["search"] = req.UserName
	}

	params["limit"] = req.Limit
	params["offset"] = offset

	query = query + filter + "ORDER BY created_at desc  LIMIT :limit OFFSET :offset"
	resQuery, pArr := helper.ReplaceQueryParams(query, params)

	rows, err := u.db.Query(context.Background(), resQuery, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&count,
			&user.ID,
			&user.Username,
			&user.Password,
			&user.IsActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = createdAt.Format(time.RFC3339)
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
		}
		resp.Count = count
		resp.Users = append(resp.Users, user)
	}

	return resp, nil
}

func (u *userRepo) GetAllDeletedUser(ctx context.Context, req *models.GetAllUserRequest) (resp *models.GetAllUser, err error) {
	params := make(map[string]interface{})
	filter := ""

	resp = &models.GetAllUser{}

	resp.Users = make([]models.User, 0)

	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	updatedAt := sql.NullTime{}
	deletedAt := sql.NullTime{}

	fmt.Println("storage", req)
	query := `SELECT 
	               COUNT(*) OVER(),
	               id,
				   username,
				   password,
				   is_active,
				   created_at,
				   updated_at,
				   deleted_at
				   
			FROM users
			WHERE deleted_at is not null`

	if req.UserName != "" {
		filter += ` AND username ILIKE '%' || :search || '%' `
		params["search"] = req.UserName
	}

	params["limit"] = req.Limit
	params["offset"] = offset

	query = query + filter + " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"
	resQuery, pArr := helper.ReplaceQueryParams(query, params)

	rows, err := u.db.Query(context.Background(), resQuery, pArr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&count,
			&user.ID,
			&user.Username,
			&user.Password,
			&user.IsActive,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = createdAt.Format(time.RFC3339)
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
		}

		resp.Users = append(resp.Users, user)
	}
	resp.Count = count
	return resp, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, req *models.User) (string, error) {
	TokenInfo := ctx.Value("user_info").(helper.TokenInfo)
	query := `UPDATE users SET 
					  username = $1, 
					  password = $2, 
					  updated_at = NOW() 
					  WHERE 
					  "is_active"  AND 
					  id = $3 RETURNING id`
	result, err := u.db.Exec(context.Background(), query, req.Username, req.Password, TokenInfo.UserID)
	if err != nil {
		return "Error Update User", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("user not found")
	}

	return req.ID, nil
}

func (b *userRepo) DeleteUser(c context.Context, req *models.IdRequest) (resp string, err error) {
	TokenInfo := c.Value("user_info").(helper.TokenInfo)

	query := `
	 	UPDATE "users" 
		SET 
			"is_active" = $1,
			"deleted_at" = NOW() 
		WHERE 
			"is_active" = true AND
			"id" = $2
`
	result, err := b.db.Exec(
		context.Background(),
		query,
		false,
		TokenInfo.UserID,
	)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("user with ID %s not found", req.Id)

	}

	return req.Id, nil
}

func (u *userRepo) GetByLogin(ctx context.Context, req *models.LoginRequest) (*models.LoginDataRespond, error) {
	query := `SELECT
		id,
		username,
		password
	FROM users
	WHERE username = $1`

	user := models.LoginDataRespond{}
	err := u.db.QueryRow(context.Background(), query, req.Username).Scan(
		&user.UserId,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
