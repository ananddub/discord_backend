-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS invites (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    channel_id INTEGER REFERENCES channels(id) ON DELETE CASCADE,
    inviter_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    max_uses INTEGER DEFAULT 0,
    uses INTEGER DEFAULT 0,
    max_age INTEGER DEFAULT 0,
    temporary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_invites_code ON invites(code);
CREATE INDEX idx_invites_server_id ON invites(server_id);
CREATE INDEX idx_invites_inviter_id ON invites(inviter_id);
CREATE INDEX idx_invites_expires_at ON invites(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_invites_expires_at;
DROP INDEX IF EXISTS idx_invites_inviter_id;
DROP INDEX IF EXISTS idx_invites_server_id;
DROP INDEX IF EXISTS idx_invites_code;
DROP TABLE IF EXISTS invites;
-- +goose StatementEnd
