-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    last_transaction TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users,
    transaction_date DATE NOT NULL,
    description TEXT NOT NULL,
    amount TEXT NOT NULL,
    transaction_type VARCHAR(10) CHECK (transaction_type IN ('debit', 'credit')),
    balance TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
DROP TABLE users;
-- +goose StatementEnd
