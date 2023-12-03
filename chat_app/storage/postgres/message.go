package postgres

import (
	"context"
	"fmt"
	"user/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type messageRepo struct {
	db *pgxpool.Pool
}

func NewMessageRepo(db *pgxpool.Pool) *messageRepo {
	return &messageRepo{
		db: db,
	}
}

func (b *messageRepo) CreateMessage(c context.Context, req *models.Message) (string, error) {

	query := `INSERT INTO messages(
					room_id, 
					user_id, 
					content,
					created_at
					) VALUES ($1, $2, $3, NOW())`
	_, err := b.db.Exec(context.Background(), query,
		req.RoomID,
		req.Username,
		req.Content,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create message: %w", err)
	}
	fmt.Println("message inserted ))")
	return "created", nil
}
