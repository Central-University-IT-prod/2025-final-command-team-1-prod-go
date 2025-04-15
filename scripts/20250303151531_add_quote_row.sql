-- +goose Up
ALTER TABLE posts ADD quote TEXT DEFAULT '';
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
ALTER TABLE posts DROP COLUMN quote;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
