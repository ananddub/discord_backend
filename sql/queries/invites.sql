-- name: CreateInvite :one
INSERT INTO
    invites (
        code,
        server_id,
        channel_id,
        inviter_id,
        max_uses,
        max_age,
        temporary,
        expires_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
RETURNING
    *;

-- name: GetInviteByCode :one
SELECT * FROM invites WHERE code = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetServerInvites :many
SELECT *
FROM invites
WHERE
    server_id = $1
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: IncrementInviteUses :one
UPDATE invites
SET
    uses = uses + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    code = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SoftDeleteInvite :one
UPDATE invites
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    code = $1
RETURNING
    *;

-- name: HardDeleteInvite :one
DELETE FROM invites WHERE code = $1 RETURNING *;

-- name: RestoreInvite :one
UPDATE invites
SET
    is_deleted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    code = $1
RETURNING
    *;

-- name: SoftDeleteExpiredInvites :one
UPDATE invites
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    expires_at IS NOT NULL
    AND expires_at < CURRENT_TIMESTAMP
    AND is_deleted = FALSE
RETURNING
    *;

-- name: HardDeleteExpiredInvites :one
DELETE FROM invites
WHERE
    expires_at IS NOT NULL
    AND expires_at < CURRENT_TIMESTAMP
RETURNING
    *;

-- name: SoftDeleteInvitesByServer :one
UPDATE invites
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    server_id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: HardDeleteInvitesByServer :one
DELETE FROM invites WHERE server_id = $1 RETURNING *;
