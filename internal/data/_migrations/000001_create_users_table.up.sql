CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    username text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    version integer NOT NULL DEFAULT 1
);