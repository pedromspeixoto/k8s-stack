-- +migrate Down
-- +migrate StatementBegin
DROP TABLE IF EXISTS todos;
-- +migrate StatementCommit