-- +goose Up

-- Now create the customers table
CREATE TABLE customer_subscriptions (
    customer_id VARCHAR(200) NOT NULL,
    price_id VARCHAR(100) NOT NULL,
    promo_code VARCHAR(100) NULL,
    subscription_id VARCHAR(100) NOT NULL,
    subscription_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS customer_subscriptions;