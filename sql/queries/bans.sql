-- name: CreateBan :one
INSERT INTO bans (
    server_id, user_id, moderator_id, reason, expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBan :one
SELECT * FROM bans
WHERE server_id = $1 AND user_id = $2 AND is_deleted = FALSE LIMIT 1;

-- name: GetServerBans :many
SELECT * FROM bans
WHERE server_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: SoftDeleteBan :exec
UPDATE bans
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE server_id = $1 AND user_id = $2;

-- name: HardDeleteBan :exec
DELETE FROM bans
WHERE server_id = $1 AND user_id = $2;

-- name: RestoreBan :exec
UPDATE bans
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE server_id = $1 AND user_id = $2;

-- name: SoftDeleteExpiredBans :exec
UPDATE bans
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP AND is_deleted = FALSE;

-- name: HardDeleteExpiredBans :exec
DELETE FROM bans
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP;

-- name: IsUserBanned :one
SELECT EXISTS(
    SELECT 1 FROM bans
    WHERE server_id = $1 AND user_id = $2 AND is_deleted = FALSE
    AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
) AS is_banned;
