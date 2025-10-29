-- name: CreateInvite :one
INSERT INTO invites (
    code, server_id, channel_id, inviter_id, max_uses, max_age, temporary, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetInviteByCode :one
SELECT * FROM invites
WHERE code = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetServerInvites :many
SELECT * FROM invites
WHERE server_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: IncrementInviteUses :exec
UPDATE invites
SET uses = uses + 1, updated_at = CURRENT_TIMESTAMP
WHERE code = $1 AND is_deleted = FALSE;

-- name: SoftDeleteInvite :exec
UPDATE invites
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE code = $1;

-- name: HardDeleteInvite :exec
DELETE FROM invites
WHERE code = $1;

-- name: RestoreInvite :exec
UPDATE invites
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE code = $1;

-- name: SoftDeleteExpiredInvites :exec
UPDATE invites
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP AND is_deleted = FALSE;

-- name: HardDeleteExpiredInvites :exec
DELETE FROM invites
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP;

-- name: SoftDeleteInvitesByServer :exec
UPDATE invites
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE server_id = $1 AND is_deleted = FALSE;

-- name: HardDeleteInvitesByServer :exec
DELETE FROM invites
WHERE server_id = $1;
