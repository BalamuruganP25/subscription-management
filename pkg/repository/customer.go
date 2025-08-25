package repository

import (
	"context"
	"fmt"

	"github.com/lib/pq"
)

func (r *CurdRepository) CreateCustomer(ctx context.Context, id, name, email, phone string) error {
	query := `INSERT INTO customers (id, name, email_id, phone_number) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, id, name, email, phone)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				// This is a unique violation error
				return fmt.Errorf("email already exists: %w", err)
			}
		}

		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

func (r *CurdRepository) CreateSubscription(
	ctx context.Context,
	customerID, priceID, promoCode, subscriptionID, subscriptionStatus string,
) error {
	query := `
        INSERT INTO customer_subscriptions (
            customer_id, price_id, promo_code, subscription_id, subscription_status
        ) VALUES ($1, $2, $3, $4, $5)
    `

	// Convert empty promoCode to nil
	var promoCodeVal interface{}
	if promoCode == "" {
		promoCodeVal = nil
	} else {
		promoCodeVal = promoCode
	}

	_, err := r.db.ExecContext(ctx, query,
		customerID,
		priceID,
		promoCodeVal,
		subscriptionID,
		subscriptionStatus,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

func (r *CurdRepository) UpdateSubscription(ctx context.Context, subscriptionID, status string) error {
	query := `UPDATE customer_subscriptions SET subscription_status = $1 WHERE subscription_id = $2`
	result, err := r.db.ExecContext(ctx, query, status, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found with ID: %s", subscriptionID)
	}
	return nil
}
