-- +goose Up
CREATE INDEX IF NOT EXISTS idx_recipecomponents_recipeid ON recipe_components USING BTREE(recipe_id);

-- +goose Down
DROP INDEX IF EXISTS idx_recipecomponents_recipeid;