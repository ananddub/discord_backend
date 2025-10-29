-- name: AddServerMember :one
INSERT INTO server_members (
    server_id, user_id, nickname
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetServerMember :one
SELECT * FROM server_members
WHERE server_id = $1 AND user_id = $2 LIMIT 1;

-- name: GetServerMembers :many
SELECT * FROM server_members
WHERE server_id = $1
ORDER BY joined_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateMemberNickname :exec
UPDATE server_members
SET nickname = $3, updated_at = CURRENT_TIMESTAMP
WHERE server_id = $1 AND user_id = $2;

-- name: RemoveServerMember :exec
DELETE FROM server_members
WHERE server_id = $1 AND user_id = $2;

-- name: GetUserServerMemberships :many
SELECT * FROM server_members
WHERE user_id = $1
ORDER BY joined_at DESC;

-- name: CountServerMembers :one
SELECT COUNT(*) FROM server_members
WHERE server_id = $1;

-- name: UpdateMemberMuteStatus :exec
UPDATE server_members
SET is_muted = $3, is_deafened = $4, updated_at = CURRENT_TIMESTAMP
WHERE server_id = $1 AND user_id = $2;
