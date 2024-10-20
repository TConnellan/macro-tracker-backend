-- +goose Up
CREATE INDEX IF NOT EXISTS idx_consumables_name ON consumables USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS idx_consumables_brand_name ON consumables USING GIN (to_tsvector('simple', brand_name));

-- +goose Down
DROP INDEX IF EXISTS idx_consumables_name;
DROP INDEX IF EXISTS idx_consumables_brand_name;