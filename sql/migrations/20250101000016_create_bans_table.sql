-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bans (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    moderator_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    reason TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(server_id, user_id)
);

-- Create indexes
CREATE INDEX idx_bans_server_id ON bans(server_id);
CREATE INDEX idx_bans_user_id ON bans(user_id);
CREATE INDEX idx_bans_expires_at ON bans(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_bans_expires_at;
DROP INDEX IF EXISTS idx_bans_user_id;
DROP INDEX IF EXISTS idx_bans_server_id;
DROP TABLE IF EXISTS bans;
-- +goose StatementEnd
