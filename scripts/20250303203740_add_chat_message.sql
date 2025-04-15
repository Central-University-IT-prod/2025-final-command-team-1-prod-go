-- +goose Up
CREATE TABLE IF NOT EXISTS chat_messages (
    email varchar(64),
    text TEXT,
    role varchar(10),
    created_at timestamp,

    CONSTRAINT fk_user_email FOREIGN KEY (email) REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE 
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE chat_messages
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd