BEGIN;
INSERT INTO users (username, email, password_hash, created_at) VALUES (
 'John Doe',
 'John@email.com',
 '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
 '2024-01-01 10:00:00'
),(
 'Jack Brabham',
 'Jack@email.com',
 '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
 '2024-01-01 10:00:00'    
),(
 'user numero 3',
 'user3@email.com',
 '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
 '2024-01-01 10:00:00'    
),(
 'user numero 4',
 'user4@email.com',
 '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
 '2024-01-01 10:00:00'    
);

INSERT INTO consumables (creator_id, created_at, name, brand_name, size, units, carbs, fats, proteins, alcohol) VALUES 
(1, '2024-01-01 10:00:00', 'Oats', 'Uncle Tobys', 100, 'g', 40, 0.5, 3, 0),
(1, '2024-01-01 10:00:00', 'Cavendish Banana', 'Coles', 100, 'g', 38, 0.1, 2, 0),
(1, '2024-01-01 10:00:00', 'Greek Yogurt', 'Jalna', 90, 'g', 3.8, 5.0, 9.0, 0),
(1, '2024-01-01 10:00:00', 'Wholemeal Bread', 'Tip Top', 110, 'g', 41.8, 2.2, 8.8, 0),
(1, '2024-01-01 10:00:00', 'Red Apple', 'Aldi', 95, 'g', 14.0, 0.2, 0.3, 0),
(1, '2024-01-01 10:00:00', 'Chicken Breast', 'IGA', 105, 'g', 0, 2.6, 22.5, 0),
(2, '2024-01-01 10:00:00', 'Almond Milk', 'Vitasoy', 250, 'ml', 0.8, 1.2, 0.5, 0),
(2, '2024-01-01 10:00:00', 'Sweet Potato', 'Woolworths', 150, 'g', 27.5, 0.1, 2.0, 0),
(2, '2024-01-01 10:00:00', 'Salmon Fillet', 'Tassal', 125, 'g', 0, 12.5, 25.0, 0),
(2, '2024-01-01 10:00:00', 'Quinoa', 'Coles', 85, 'g', 15.6, 2.4, 4.8, 0),
(2, '2024-01-01 10:00:00', 'Red Wine', 'Penfolds', 150, 'ml', 3.8, 0, 0, 13.5),
(3, '2024-01-01 10:00:00', 'Full Cream Milk', 'Dairy Farmers', 250, 'ml', 12.5, 8.8, 8.5, 0),
(3, '2024-01-01 10:00:00', 'Skim Milk', 'Pauls', 250, 'ml', 12.8, 0.3, 8.8, 0),
(3, '2024-01-01 10:00:00', 'Soy Milk', 'Sanitarium', 250, 'ml', 5.5, 3.2, 7.5, 0),
(3, '2024-01-01 10:00:00', 'Oat Milk', 'Oatly', 250, 'ml', 16.0, 3.0, 1.0, 0),
(3, '2024-01-01 10:00:00', 'Coconut Milk', 'Pure Harvest', 250, 'ml', 3.0, 5.0, 0.5, 0),
(4, '2024-01-01 10:00:00', 'Lasagne Pasta Large', 'San Remo', 62.5, 'g', 46.6, 0.9, 7.9, 0),
(4, '2024-01-01 10:00:00', 'No Added Hormone Beef 5 Star Extra Trim Mince', 'Coles', 100, 'g', 0.5, 2, 21.3, 0);


INSERT INTO recipes (recipe_name, creator_id, created_at, last_edited_at, notes, parent_recipe_id, is_latest) VALUES 
('Lasagne', 1, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe', NULL, 'true'),
('recipe2', 2, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 2', NULL, 'true'),
('recipe3', 2, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 3', NULL, 'true'),
('recipe4', 3, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 4', 1, 'true'),
('Recipe5', 2, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 5', NULL, 'true'),
('doesntmatchsearch', 2, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'not matching', NULL, 'true'),
('recipe7', 3, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 7', NULL, 'true'),
('recipe8', 2, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 8', 1, 'true'),
('recipe9', 4, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 9', 8, 'true'),
('recipe10', 4, '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'a recipe 10', 9, 'true');

INSERT INTO pantry_items (user_id, consumable_id, name, created_at, last_modified) VALUES
(1, 17, 'lasagne sheet', '2024-01-01 10:00:00', '2024-01-01 10:00:00'),
(1, 18, 'Minced Beef', '2024-01-01 10:00:00', '2024-01-01 10:00:00'),
(1, 18, 'Minced Beef', '2024-01-01 10:00:00', '2024-01-01 10:00:00');

INSERT INTO recipe_components (recipe_id, pantry_item_id, created_at, quantity, step_no, step_description) VALUES
(1, 1, '2024-01-01 10:00:00', 4, 1, 'step 1'),
(1, 2, '2024-01-01 10:00:00', 5, 2, 'step 2'),
(7, 3, '2024-01-01 10:00:00', 5, 1, 'step 3');

INSERT INTO consumed (user_id, recipe_id, quantity, carbs, fats, proteins, alcohol, consumed_at, created_at, last_edited_at, notes) VALUES
(1, 1, 1, 1, 1, 1, 1, '2024-01-01 10:00:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes'),
(1, 2, 5, 7, 7, 7, 1, '2024-01-01 10:00:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes 2'),
(1, 2, 8, 3, 2, 1, 0, '2024-01-01 10:00:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes 3'),
(2, 2, 8, 3, 2, 1, 0, '2024-01-01 10:00:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes 4'),
(3, 1, 1, 1, 1, 1, 1, '2024-01-01 10:00:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes'),
(3, 2, 5, 7, 7, 7, 1, '2024-01-01 10:01:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes 2'),
(3, 2, 8, 3, 2, 1, 0, '2024-01-01 10:10:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes 3'),
(3, 2, 8, 3, 2, 1, 0, '2024-01-02 10:00:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', 'notes 3');

COMMIT;