BEGIN;

ALTER TABLE recipe_components DROP CONSTRAINT fk_recipecomponent_pantry_item;
ALTER TABLE recipe_components RENAME COLUMN pantry_item_id TO consumable_id;
ALTER TABLE recipe_components ADD CONSTRAINT fk_recipecomponent_consumable FOREIGN KEY (consumable_id) REFERENCES consumable(id) ON DELETE RESTRICT;

DROP TABLE IF EXISTS pantry_items CASCADE;

COMMIT;