-- +goose Up
CREATE INDEX IF NOT EXISTS idx_users_email ON users USING HASH(email);

-- +goose Down
DROP INDEX IF EXISTS idx_username;
DROP INDEX IF EXISTS idx_email;