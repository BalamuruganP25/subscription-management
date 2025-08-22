
-- +goose Up

-- Now create the customers table
CREATE TABLE customers (
    id VARCHAR(200) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email_id VARCHAR(100) NOT NULL UNIQUE,
    phone_number VARCHAR(15) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;