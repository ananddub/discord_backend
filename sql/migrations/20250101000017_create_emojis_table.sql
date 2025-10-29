-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS emojis (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    creator_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    require_colons BOOLEAN DEFAULT TRUE,
    managed BOOLEAN DEFAULT FALSE,
    animated BOOLEAN DEFAULT FALSE,
    available BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(server_id, name)
);

-- Create indexes
CREATE INDEX idx_emojis_server_id ON emojis(server_id);
CREATE INDEX idx_emojis_name ON emojis(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_emojis_name;
DROP INDEX IF EXISTS idx_emojis_server_id;
DROP TABLE IF EXISTS emojis;
-- +goose StatementEnd
