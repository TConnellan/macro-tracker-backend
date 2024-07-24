CREATE TABLE IF NOT EXISTS recipes (
   id INTEGER               PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
   recipe_name VARCHAR(50)  ,
   creator_id INTEGER       ,
   created_at TIMESTAMP     DEFAULT current_timestamp,
   last_edited_at TIMESTAMP ,
   notes TEXT
);

ALTER TABLE recipe ADD CONSTRAINT fk_recipe_creator FOREIGN KEY (creator_id) REFERENCES users(id);