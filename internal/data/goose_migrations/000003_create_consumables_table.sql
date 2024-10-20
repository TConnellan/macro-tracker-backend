-- +goose Up
Begin;
CREATE TABLE IF NOT EXISTS consumables (
   id INTEGER               PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
   creator_id INTEGER,
   created_at timestamp default current_timestamp,
   name VARCHAR(50)    ,
   brand_name VARCHAR(50)   ,
   size DOUBLE PRECISION    ,
   units VARCHAR(10)        ,
   carbs DOUBLE PRECISION,
   fats DOUBLE PRECISION    ,
   proteins  DOUBLE PRECISION,
   alcohol DOUBLE PRECISION
);
COMMIT;
ALTER TABLE consumables ADD CONSTRAINT fk_consumable_creator FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE consumables DROP CONSTRAINT IF EXISTS fk_consumable_creator;
DROP TABLE IF EXISTS consumables CASCADE;

