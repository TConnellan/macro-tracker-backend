-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    username text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    version integer NOT NULL DEFAULT 1
);

for up_file in *.up.sql; do
    down_file="${up_file%.up.sql}.down.sql"
    if [ -f "$down_file" ]; then
        echo -e "\n" >> "$up_file"
        cat "$down_file" >> "$up_file"
    fi
done

-- +goose Down
DROP TABLE IF EXISTS users;