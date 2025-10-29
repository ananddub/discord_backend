-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channel_permissions (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    allow_permissions BIGINT DEFAULT 0,
    deny_permissions BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CHECK ((role_id IS NOT NULL AND user_id IS NULL) OR (role_id IS NULL AND user_id IS NOT NULL)),
    UNIQUE(channel_id, role_id, user_id)
);

-- Create indexes
CREATE INDEX idx_channel_permissions_channel_id ON channel_permissions(channel_id);
CREATE INDEX idx_channel_permissions_role_id ON channel_permissions(role_id);
CREATE INDEX idx_channel_permissions_user_id ON channel_permissions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_channel_permissions_user_id;
DROP INDEX IF EXISTS idx_channel_permissions_role_id;
DROP INDEX IF EXISTS idx_channel_permissions_channel_id;
DROP TABLE IF EXISTS channel_permissions;
-- +goose StatementEnd
