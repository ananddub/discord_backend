-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS message_mentions (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CHECK ((user_id IS NOT NULL AND role_id IS NULL) OR (user_id IS NULL AND role_id IS NOT NULL))
);

-- Create indexes
CREATE INDEX idx_message_mentions_message_id ON message_mentions(message_id);
CREATE INDEX idx_message_mentions_user_id ON message_mentions(user_id);
CREATE INDEX idx_message_mentions_role_id ON message_mentions(role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_message_mentions_role_id;
DROP INDEX IF EXISTS idx_message_mentions_user_id;
DROP INDEX IF EXISTS idx_message_mentions_message_id;
DROP TABLE IF EXISTS message_mentions;
-- +goose StatementEnd
