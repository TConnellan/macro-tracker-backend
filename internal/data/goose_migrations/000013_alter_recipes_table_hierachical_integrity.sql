-- +goose Up
ALTER TABLE recipes ADD CONSTRAINT ancestor_references_are_descending CHECK (parent_recipe_id IS NULL OR id > parent_recipe_id);

-- +goose Down
DROP TABLE IF EXISTS tokens;
ALTER TABLE recipes DROP CONSTRAINT ancestor_references_are_descending;