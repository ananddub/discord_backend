-- name: CreateInvite :one
INSERT INTO invites (
    code, server_id, channel_id, inviter_id, max_uses, max_age, temporary, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetInviteByCode :one
SELECT * FROM invites
WHERE code = $1 LIMIT 1;

-- name: GetServerInvites :many
SELECT * FROM invites
WHERE server_id = $1
ORDER BY created_at DESC;

-- name: IncrementInviteUses :exec
UPDATE invites
SET uses = uses + 1
WHERE code = $1;

-- name: DeleteInvite :exec
DELETE FROM invites
WHERE code = $1;

-- name: DeleteExpiredInvites :exec
DELETE FROM invites
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP;

-- name: DeleteInvitesByServer :exec
DELETE FROM invites
WHERE server_id = $1;
