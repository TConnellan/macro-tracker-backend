-- +goose Up
BEGIN;
CREATE TABLE IF NOT EXISTS recipes (
   id INTEGER               PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
   recipe_name VARCHAR(50)  ,
   creator_id INTEGER       ,
   created_at TIMESTAMP     DEFAULT current_timestamp,
   last_edited_at TIMESTAMP DEFAULT current_timestamp,
   notes TEXT
);
COMMIT;
ALTER TABLE recipes ADD CONSTRAINT fk_recipe_creator FOREIGN KEY (creator_id) REFERENCES users(id);

-- +goose Down
DROP TABLE IF EXISTS recipes CASCADE;