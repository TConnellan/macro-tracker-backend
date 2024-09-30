CREATE DATABASE test_macrotracker
WITH OWNER 'postgres'
ENCODING 'UTF8'
TABLESPACE 'pg-default';


CREATE USER test_web WITH PASSWORD 'pass';
GRANT CREATE, CONNECT ON DATABASE test_snippetbox TO test_web;
\c test_macrotracker;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO test_web;

GRANT CREATE ON SCHEMA public TO test_web;

CREATE EXTENSION IF NOT EXISTS citext;

-- GRANT ALL PRIVILEGES ON DATABASE test_macrotracker TO test_web;
