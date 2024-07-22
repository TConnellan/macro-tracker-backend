CREATE TABLE IF NOT EXISTS consumables (
   id INTEGER               PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
   name VARCHAR(50)    ,
   brand_name VARCHAR(50)   ,
   size DOUBLE PRECISION    ,
   units VARCHAR(10)        ,
   carbohydrates DOUBLE PRECISION,
   fat DOUBLE PRECISION    ,
   protein  DOUBLE PRECISION,
   alcohol DOUBLE PRECISION
);