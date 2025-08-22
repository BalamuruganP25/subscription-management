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
	CreateUser(ctx context.Context, name, email_id, phone_number string) (string, error)
	GetUser(ctx context.Context, id string) (UserResponse, error)
	UpdateUser(ctx context.Context, id string, phone_number string) error
	DeleteUser(ctx context.Context, id string) error
}

func (r *CurdRepository) CreateUser(ctx context.Context, name, email, phone string) (string, error) {
	query := `INSERT INTO users (name, email_id, phone_number) VALUES (?, ?, ?) RETURNING id`
	var id int64
	err := r.db.QueryRowContext(ctx, query, name, email, phone).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return fmt.Sprintf("%d", id), nil
}

func (r *CurdRepository) GetUser(ctx context.Context, id string) (UserResponse, error) {
	query := `SELECT id, name, email_id, phone_number FROM users WHERE id = ?`
	var user UserResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email_id, &user.Phone_number)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserResponse{}, fmt.Errorf("user not found: %w", err)
		}
		return UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *CurdRepository) UpdateUser(ctx context.Context, id string, phone_number string) error {
	query := `UPDATE users SET phone_number = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, phone_number, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id: %s", id)
	}
	return nil
}

func (r *CurdRepository) DeleteUser(ctx context.Context, id string) error {
	query := `UPDATE users SET status = false WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id: %s", id)
	}
	return nil
}
