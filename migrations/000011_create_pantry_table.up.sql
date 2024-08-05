BEGIN;

CREATE TABLE IF NOT EXISTS pantry_items (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL,
    consumable_id INTEGER NOT NULL,
    name VARCHAR(50),
    created_at TIMESTAMP DEFAULT current_timestamp,
    last_modified TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE pantry_items ADD CONSTRAINT fk_pantryitem_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE pantry_items ADD CONSTRAINT fk_pantryitem_consumable FOREIGN KEY (consumable_id) REFERENCES consumables(id) ON DELETE RESTRICT;

INSERT INTO pantry_items(user_id, consumable_id, name)
VALUES (
    SELECT DISTINCT RC.user_id, C.consumable_id, C.name
    FROM recipe_components RC INNER JOIN consumables C ON RC.consumable_id = C.id
);

ALTER TABLE recipe_components DROP CONSTRAINT fk_recipecomponent_consumable;

ALTER TABLE recipe_components RENAME COLUMN consumable_id TO pantry_item_id;

ALTER TABLE recipe_components ADD CONSTRAINT 
fk_recipecomponent_pantry_item FOREIGN KEY (pantry_item_id) REFERENCES pantry_items(id) ON DELETE RESTRICT;

UPDATE recipe_components RC
SET pantry_item_id = P.id
FROM pantry_items P
WHERE RC.user_id = P.user_id AND RC.consumable_id = P.consumable_id;

COMMIT;
