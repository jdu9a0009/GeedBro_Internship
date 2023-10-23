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
                    updated_at,
                    deleted_at
				FROM
				    users
				WHERE
				    id = $1;`
	// WHERE
	// deleted_at IS NOT NULL
	// AND id = $1;`
	var (
		createdAt  sql.NullTime
		updatedAt  sql.NullTime
		deleted_at sql.NullTime
	)
	user := models.User{}

	err = u.db.QueryRow(context.Background(), query, req.Id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsActive,
		&createdAt,
		&updatedAt,
		&deleted_at,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	user.CreatedAt = createdAt.Time.Format(time.RFC3339)
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	if deleted_at.Valid {
		user.DeletedAt = deleted_at.Time.Format(time.RFC3339)
	}

	return &user, nil
}

func (u *userRepo) GetAllUser(ctx context.Context, req *models.GetAllUserRequest) (*models.GetAllUser, error) {
	params := make(map[string]interface{})
	filter := ""
	offset := (req.Page - 1) * req.Limit
	createdAt := time.Time{}
	updatedAt := sql.NullTime{}

	s := `SELECT 
	               id,
				   username,
				   password,
				   is_active,
				   created_at,
				   updated_at
			FROM users
						WHERE deleted_at IS  NULL`

	if req.UserName != "" {
		filter += ` AND username ILIKE '%' || :search || '%' `
		params["search"] = req.UserName
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

	resp := &models.GetAllUser{}
	count := 0
	for rows.Next() {
		var user models.User
		count++
		err := rows.Scan(
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

		resp.Users = append(resp.Users, user)
	}

	resp.Count = count
	return resp, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, req *models.User) (string, error) {

	query := `UPDATE users SET 
					  username = $1, 
					  password = $2, 
					  updated_at = NOW() 
					  WHERE 
					  "is_active" = true AND 
					  id = $3 RETURNING id`
	result, err := u.db.Exec(ctx, query, req.Username, req.Password, req.ID)
	if err != nil {
		return "Error Update User", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("user not found")
	}

	return req.ID, nil
}

func (b *userRepo) DeleteUser(c context.Context, req *models.IdRequest) (resp string, err error) {
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
		req.Id,
	)
	if err != nil {
		return "", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("user with ID %s not found", req.Id)

	}

	return req.Id, nil
}

// func (u *userRepo) GetByLogin(ctx context.Context, req *models.LoginRequest) (*models.LoginRespond, error) {
// 	query := `SELECT
// 		id,
// 		username,
// 		password
// 	FROM users
// 	WHERE username = $1`

// 	user := models.LoginRespond{}
// 	err := u.db.QueryRow(context.Background(), query, req.Username).Scan(
// 		&user.ID,
// 		&user.Username,
// 		&user.Password,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("user not found")
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }
