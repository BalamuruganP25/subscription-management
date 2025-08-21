package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type CurdRepository struct {
	db *sql.DB
}

func NewCurdRepo(db *sql.DB) *CurdRepository {
	return &CurdRepository{
		db: db,
	}
}

type CrudRepo interface {
	CreateUser(ctx context.Context, name, email, phone string) (string, error)
}

func (r *CurdRepository) CreateUser(ctx context.Context, name, email, phone string) (string, error) {
	query := `INSERT INTO users (name, email, phone) VALUES (?, ?, ?) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, name, email, phone).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return fmt.Sprintf("%d", id), nil
}
