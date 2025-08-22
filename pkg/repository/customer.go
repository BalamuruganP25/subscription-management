package repository

import (
	"context"
	"fmt"
)

func (r *CurdRepository) CreateCustomer(ctx context.Context, id, name, email, phone string) error {
	query := `INSERT INTO customers (id, name, email_id, phone_number) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, id, name, email, phone)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}
