-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(255),
    banner VARCHAR(255),
    description TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    region VARCHAR(50) DEFAULT 'auto',
    member_count INTEGER DEFAULT 1,
    is_verified BOOLEAN DEFAULT FALSE,
    vanity_url VARCHAR(50) UNIQUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_servers_owner_id ON servers(owner_id);
CREATE INDEX idx_servers_vanity_url ON servers(vanity_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_servers_vanity_url;
DROP INDEX IF EXISTS idx_servers_owner_id;
DROP TABLE IF EXISTS servers;
-- +goose StatementEnd
