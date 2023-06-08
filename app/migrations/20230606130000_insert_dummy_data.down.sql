-- +migrate Down
-- +migrate StatementBegin
DELETE FROM todos WHERE ID IN (1, 2, 3); 
-- +migrate StatementCommit