-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS channels (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES channels(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'voice', 'category', 'announcement', 'stage', 'forum', 'dm', 'group_dm')),
    position INTEGER DEFAULT 0,
    topic VARCHAR(1024),
    is_nsfw BOOLEAN DEFAULT FALSE,
    slowmode_delay INTEGER DEFAULT 0,
    user_limit INTEGER DEFAULT 0,
    bitrate INTEGER DEFAULT 64000,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_channels_server_id ON channels(server_id);
CREATE INDEX idx_channels_category_id ON channels(category_id);
CREATE INDEX idx_channels_type ON channels(type);
CREATE INDEX idx_channels_position ON channels(server_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_channels_position;
DROP INDEX IF EXISTS idx_channels_type;
DROP INDEX IF EXISTS idx_channels_category_id;
DROP INDEX IF EXISTS idx_channels_server_id;
DROP TABLE IF EXISTS channels;
-- +goose StatementEnd
