-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dm_participants (
    id SERIAL PRIMARY KEY,
    dm_channel_id INTEGER NOT NULL REFERENCES dm_channels(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_message_id INTEGER,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(dm_channel_id, user_id)
);

-- Create indexes
CREATE INDEX idx_dm_participants_dm_channel_id ON dm_participants(dm_channel_id);
CREATE INDEX idx_dm_participants_user_id ON dm_participants(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_dm_participants_user_id;
DROP INDEX IF EXISTS idx_dm_participants_dm_channel_id;
DROP TABLE IF EXISTS dm_participants;
-- +goose StatementEnd
