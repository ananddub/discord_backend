-- +goose Up
CREATE TABLE text_groups (
    id SERIAL PRIMARY KEY,
    channel_id TEXT NOT NULL,
    topic TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE text_channels (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES text_groups(id) ON DELETE CASCADE,
    topic TEXT,
    is_archived BOOLEAN DEFAULT FALSE NOT NULL,
    archived_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE text_messages (
    id SERIAL PRIMARY KEY,
    text_channel_id INTEGER NOT NULL REFERENCES text_channels(id) ON DELETE CASCADE,
    sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reply_to_message_id INTEGER REFERENCES text_messages(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW() NOT NULL,
    is_edited BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- +goose Down
DROP TABLE text_messages;
DROP TABLE text_channels;
DROP TABLE text_groups;
