CREATE INDEX IF NOT EXISTS idx_users_username ON users (name);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
-- these could be made on a hash column of these values