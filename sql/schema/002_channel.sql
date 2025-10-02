-- +goose Up
CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    unique_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    pic TEXT,
    position INTEGER DEFAULT -1 NOT NULL,
    description TEXT,
    group_id INTEGER DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- +goose Down
DROP TABLE channels;
