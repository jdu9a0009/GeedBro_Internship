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
							name,
							phone,
							username,
							password,
							created_at)
							VALUES($1, $2, $3, $4, $5, now())`

	_, err := u.db.Exec(ctx, query,
		id,
		req.Name,
		req.Phone,
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
				   name,
				   phone,
				   username,
				   password,
				   created_at,
				   updated_at from users
				   WHERE id=$1`
	var (
		createdAt time.Time
		updatedAt sql.NullTime
	)
	user := models.User{}

	err = u.db.QueryRow(ctx, query, req.Id).Scan(
		&user.ID,
		&user.Name,
		&user.Phone,
		&user.Username,
		&user.Password,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
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
				   name,
				   phone,
				   username,
				   password,
				   created_at,
				   updated_at from users
				   `
	if req.Name != "" {
		filter += ` WHERE name ILIKE '%' || :search || '%' `
		params["search"] = req.Name
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
			&user.Name,
			&user.Phone,
			&user.Username,
			&user.Password,
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
	                  name = $1, phone = $2, username = $3, password = $4, updated_at = NOW() WHERE id = $5 RETURNING id`
	result, err := u.db.Exec(ctx, query, req.Name, req.Phone, req.Password, req.ID)
	if err != nil {
		return "Error Update User", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("user not found")
	}

	return req.ID, nil
}
func (u *userRepo) DeleteUser(ctx context.Context, req *models.IdRequest) (resp string, err error) {
	query := `DELETE FROM users WHERE id = $1 RETURNING id`

	result, err := u.db.Exec(ctx, query, req.Id)
	if err != nil {
		return "Error from Delete User", err
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("User not found")
	}

	return req.Id, nil
}

func (u userRepo) GetByLogin(ctx context.Context, req *models.LoginRequest) (rep *models.LoginRespond, err error) {
	query := `SELECT 
				   phone,
				   username,
				   password
				   from users
				   WHERE "username"=$1`

	user := models.LoginRespond{}

	err = u.db.QueryRow(ctx, query, req.Username).Scan(
		&user.Phone,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
