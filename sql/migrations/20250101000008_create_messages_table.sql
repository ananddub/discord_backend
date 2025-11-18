-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER REFERENCES channels (id) ON DELETE CASCADE,
    receiver_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    ischannel BOOLEAN DEFAULT FALSE,
    sender_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'default' CHECK (
        message_type IN (
            'default',
            'reply',
            'system',
            'user_join',
            'user_leave',
            'call',
            'pinned'
        )
    ),
    reply_to_message_id INTEGER REFERENCES messages (id) ON DELETE SET NULL,
    is_edited BOOLEAN DEFAULT FALSE,
    is_pinned BOOLEAN DEFAULT FALSE,
    mention_everyone BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    edited_at TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_messages_channel_id ON messages (channel_id);

CREATE INDEX idx_messages_sender_id ON messages (sender_id);

CREATE INDEX idx_messages_created_at ON messages (channel_id, created_at DESC);

CREATE INDEX idx_messages_pinned ON messages (channel_id, is_pinned)
WHERE
    is_pinned = TRUE;

CREATE INDEX idx_messages_reply_to ON messages (reply_to_message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_reply_to;

DROP INDEX IF EXISTS idx_messages_pinned;

DROP INDEX IF EXISTS idx_messages_created_at;

DROP INDEX IF EXISTS idx_messages_sender_id;

DROP INDEX IF EXISTS idx_messages_channel_id;

DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
