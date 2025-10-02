CREATE TABLE channel_members (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    joined_at TIMESTAMP DEFAULT NOW() NOT NULL,
    UNIQUE (channel_id, user_id)
);

CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_accepted BOOLEAN DEFAULT FALSE NOT NULL,
    is_blocked BOOLEAN DEFAULT FALSE NOT NULL,
    is_favorite BOOLEAN DEFAULT FALSE NOT NULL,
    is_rejected BOOLEAN DEFAULT FALSE NOT NULL,
    is_pending BOOLEAN DEFAULT TRUE NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    UNIQUE (user_id, friend_id)
);
