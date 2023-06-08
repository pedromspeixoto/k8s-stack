-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    todo_id INT NOT NULL,
    description TEXT,
    expiration_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
-- +migrate StatementCommit