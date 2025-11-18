-- name: SetChannelPermission :one
INSERT INTO
    channel_permissions (
        channel_id,
        role_id,
        user_id,
        allow_permissions,
        deny_permissions
    )
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (channel_id, role_id, user_id) DO
UPDATE
SET
    allow_permissions = EXCLUDED.allow_permissions,
    deny_permissions = EXCLUDED.deny_permissions,
    updated_at = CURRENT_TIMESTAMP
RETURNING
    *;

-- name: GetChannelPermissions :many
SELECT * FROM channel_permissions WHERE channel_id = $1;

-- name: GetRoleChannelPermissions :one
SELECT *
FROM channel_permissions
WHERE
    channel_id = $1
    AND role_id = $2
LIMIT 1;

-- name: GetUserChannelPermissions :one
SELECT *
FROM channel_permissions
WHERE
    channel_id = $1
    AND user_id = $2
LIMIT 1;

-- name: DeleteChannelPermission :one
DELETE FROM channel_permissions WHERE id = $1 RETURNING *;

-- name: DeleteRoleChannelPermissions :one
DELETE FROM channel_permissions
WHERE
    channel_id = $1
    AND role_id = $2
RETURNING
    *;

-- name: DeleteUserChannelPermissions :one
DELETE FROM channel_permissions
WHERE
    channel_id = $1
    AND user_id = $2
RETURNING
    *;
