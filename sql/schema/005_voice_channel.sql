CREATE TABLE voice_channels (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT FALSE NOT NULL,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE voice_chats (
    id SERIAL PRIMARY KEY,
    voice_channel_id INTEGER NOT NULL REFERENCES voice_channels(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat TEXT NOT NULL,
    reply_to_message_id INTEGER REFERENCES voice_chats(id) ON DELETE SET NULL,
    joined_at TIMESTAMP DEFAULT NOW() NOT NULL,
    left_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);
