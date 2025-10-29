-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_presence (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    status VARCHAR(20) DEFAULT 'offline' CHECK (status IN ('online', 'idle', 'dnd', 'offline', 'invisible')),
    custom_status VARCHAR(255),
    custom_status_emoji VARCHAR(100),
    custom_status_expires_at TIMESTAMP,
    activity VARCHAR(255),
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_user_presence_user_id ON user_presence(user_id);
CREATE INDEX idx_user_presence_status ON user_presence(status);
CREATE INDEX idx_user_presence_last_seen ON user_presence(last_seen);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_presence_last_seen;
DROP INDEX IF EXISTS idx_user_presence_status;
DROP INDEX IF EXISTS idx_user_presence_user_id;
DROP TABLE IF EXISTS user_presence;
-- +goose StatementEnd
