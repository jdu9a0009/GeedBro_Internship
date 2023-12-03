package postgres

import (
	"context"

	"fmt"
	"user/models"

	"github.com/jackc/pgx"
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

func (u *userRepo) CreateUser(ctx context.Context, req *models.CreateUserReq) (*models.CreateUserResp, error) {
	var User models.CreateUserResp
	var lastInsertedId int

	query := `INSERT INTO users(
							username,
							password,
							email)
							VALUES($1, $2, $3)`

	_, err := u.db.Exec(ctx, query,
		req.Username,
		req.Password,
		req.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	User.ID = lastInsertedId
	User.Username = req.Username
	User.Email = req.Email
	return &User, nil
}

func (u userRepo) GetUserByEmail(ctx context.Context, req *models.LoginUserReq) (rep *models.User, err error) {
	query := `SELECT  
	                id,
                    username,
                    password,
                    email

				FROM
				    users
				WHERE
				    email = $1;`

	user := models.User{}

	err = u.db.QueryRow(context.Background(), query, req.Email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
