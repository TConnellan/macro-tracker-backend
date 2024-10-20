-- +goose Up
CREATE INDEX IF NOT EXISTS idx_consumed_consumerid ON consumed USING BTREE(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_consumed_consumerid;