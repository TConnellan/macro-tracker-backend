CREATE TABLE IF NOT EXISTS recipe_components (
   id INTEGER               PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
   recipe_id INTEGER        NOT NULL,
   component_id INTEGER     NOT NULL,
   created_at TIMESTAMP     DEFAULT current_timestamp,
   quantity DOUBLE PRECISION,   
   step_no INTEGER          ,
   step_description TEXT    
);

ALTER TABLE recipe_component ADD CONSTRAINT fk_recipecomponent_recipe FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE RESTRICT;
ALTER TABLE recipe_component ADD CONSTRAINT fk_recipecomponent_component FOREIGN KEY (component_id) REFERENCES consumables(id) ON DELETE RESTRICT;