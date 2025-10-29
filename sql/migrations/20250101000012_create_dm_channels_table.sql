-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dm_channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    icon VARCHAR(255),
    owner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    is_group BOOLEAN DEFAULT FALSE,
    last_message_id INTEGER,
    last_message_at TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_dm_channels_owner_id ON dm_channels(owner_id);
CREATE INDEX idx_dm_channels_last_message_at ON dm_channels(last_message_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_dm_channels_last_message_at;
DROP INDEX IF EXISTS idx_dm_channels_owner_id;
DROP TABLE IF EXISTS dm_channels;
-- +goose StatementEnd
