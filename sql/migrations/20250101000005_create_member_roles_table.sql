-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS member_roles (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL REFERENCES server_members(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(member_id, role_id)
);

-- Create indexes
CREATE INDEX idx_member_roles_member_id ON member_roles(member_id);
CREATE INDEX idx_member_roles_role_id ON member_roles(role_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_member_roles_role_id;
DROP INDEX IF EXISTS idx_member_roles_member_id;
DROP TABLE IF EXISTS member_roles;
-- +goose StatementEnd
