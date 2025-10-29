-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS voice_states (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    server_id INTEGER REFERENCES servers(id) ON DELETE CASCADE,
    session_id VARCHAR(100) NOT NULL,
    is_muted BOOLEAN DEFAULT FALSE,
    is_deafened BOOLEAN DEFAULT FALSE,
    self_mute BOOLEAN DEFAULT FALSE,
    self_deaf BOOLEAN DEFAULT FALSE,
    self_video BOOLEAN DEFAULT FALSE,
    self_stream BOOLEAN DEFAULT FALSE,
    suppress BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, channel_id)
);

-- Create indexes
CREATE INDEX idx_voice_states_user_id ON voice_states(user_id);
CREATE INDEX idx_voice_states_channel_id ON voice_states(channel_id);
CREATE INDEX idx_voice_states_session_id ON voice_states(session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_voice_states_session_id;
DROP INDEX IF EXISTS idx_voice_states_channel_id;
DROP INDEX IF EXISTS idx_voice_states_user_id;
DROP TABLE IF EXISTS voice_states;
-- +goose StatementEnd
