-- +goose Up
ALTER TABLE recipes ADD column parent_recipe_id INTEGER DEFAULT NULL;
ALTER TABLE recipes ADD column is_latest BOOLEAN DEFAULT TRUE;

ALTER TABLE recipes ADD CONSTRAINT recipe_child_parent_id FOREIGN KEY (parent_recipe_id) REFERENCES recipes(id) ON DELETE RESTRICT;

-- +goose Down
ALTER TABLE recipes DROP COLUMN parent_recipe_id;
ALTER TABLE recipes DROP COLUMN is_latest;