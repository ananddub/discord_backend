-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    alias_name VARCHAR(100),
    is_pending BOOLEAN DEFAULT TRUE,
    is_accepted BOOLEAN DEFAULT FALSE,
    is_blocked BOOLEAN DEFAULT FALSE,
    is_favorite BOOLEAN DEFAULT FALSE,
    is_muted BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (user_id, friend_id),
    CHECK (user_id != friend_id)
);

-- Create indexes
CREATE INDEX idx_friends_user_id ON friends (user_id);

CREATE INDEX idx_friends_friend_id ON friends (friend_id);

CREATE INDEX idx_friends_pending ON friends (user_id, is_pending);

CREATE INDEX idx_friends_accepted ON friends (user_id, is_accepted);

CREATE INDEX idx_friends_blocked ON friends (user_id, is_blocked);

CREATE INDEX idx_friends_favorite ON friends (user_id, is_favorite);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_friends_favorite;

DROP INDEX IF EXISTS idx_friends_blocked;

DROP INDEX IF EXISTS idx_friends_accepted;

DROP INDEX IF EXISTS idx_friends_pending;

DROP INDEX IF EXISTS idx_friends_friend_id;

DROP INDEX IF EXISTS idx_friends_user_id;

DROP TABLE IF EXISTS friends;
