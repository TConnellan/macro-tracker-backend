CREATE TABLE IF NOT EXISTS consumables (
   id INTEGER               PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
   name VARCHAR(50)    ,
   brand_name VARCHAR(50)   ,
   size DOUBLE PRECISION    ,
   units VARCHAR(10)        ,
   carbs DOUBLE PRECISION,
   fats DOUBLE PRECISION    ,
   proteins  DOUBLE PRECISION,
   alcohol DOUBLE PRECISION
);