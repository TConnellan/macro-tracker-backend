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
(2, '2024-01-01 10:00:00', 'Red Wine', 'Penfolds', 150, 'ml', 3.8, 0, 0, 13.5);
COMMIT;