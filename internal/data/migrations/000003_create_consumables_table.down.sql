ALTER TABLE consumables DROP CONSTRAINT IF EXISTS fk_consumable_creator;
DROP TABLE IF EXISTS consumables CASCADE;

