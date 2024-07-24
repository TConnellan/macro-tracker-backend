CREATE INDEX IF NOT EXISTS idx_consumables_name ON consumable USING GIN (to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS idx_consumables_brand_name ON consumable USING GIN (to_tsvector('simple', brand_name));