CREATE TABLE IF NOT EXISTS consumed (
   id INTEGER               PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
   user_id INTEGER          NOT NULL,
   recipe_id INTEGER        ,
   quantity DOUBLE PRECISION,
   carbs DOUBLE PRECISION   ,
   fats DOUBLE PRECISION    ,
   proteins DOUBLE PRECISION,
   alcohol DOUBLE PRECISION,
   consumed_at TIMESTAMP    ,
   created_at TIMESTAMP     DEFAULT current_timestamp,
   last_edited_at TIMESTAMP ,
   notes TEXT
);

ALTER TABLE consumed ADD CONSTRAINT fk_consumed_consumerid FOREIGN KEY (user_id) REFERENCES user_profile(id) ON DELETE CASCADE;
ALTER TABLE consumed ADD CONSTRAINT fk_consumed_recipeid FOREIGN KEY (recipe_id) REFERENCES recipe(id) ON DELETE SET NULL;