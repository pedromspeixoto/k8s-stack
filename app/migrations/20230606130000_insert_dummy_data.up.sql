-- +migrate Up
-- +migrate StatementBegin
INSERT INTO todos (todo_id, description, expiration_date, created_at, updated_at)
VALUES
    (1, 'Task 1', '2023-06-30', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 'Task 2', '2023-07-15', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 'Task 3', '2023-08-10', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- +migrate StatementCommit