-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) DEFAULT '#99AAB5',
    hoist BOOLEAN DEFAULT FALSE,
    position INTEGER DEFAULT 0,
    permissions BIGINT DEFAULT 0,
    mentionable BOOLEAN DEFAULT FALSE,
    icon VARCHAR(255),
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(server_id, name)
);

-- Create indexes
CREATE INDEX idx_roles_server_id ON roles(server_id);
CREATE INDEX idx_roles_position ON roles(server_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_roles_position;
DROP INDEX IF EXISTS idx_roles_server_id;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
