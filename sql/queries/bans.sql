-- name: CreateBan :one
INSERT INTO bans (
    server_id, user_id, moderator_id, reason, expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBan :one
SELECT * FROM bans
WHERE server_id = $1 AND user_id = $2 LIMIT 1;

-- name: GetServerBans :many
SELECT * FROM bans
WHERE server_id = $1
ORDER BY created_at DESC;

-- name: DeleteBan :exec
DELETE FROM bans
WHERE server_id = $1 AND user_id = $2;

-- name: DeleteExpiredBans :exec
DELETE FROM bans
WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP;

-- name: IsUserBanned :one
SELECT EXISTS(
    SELECT 1 FROM bans
    WHERE server_id = $1 AND user_id = $2
    AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
) AS is_banned;
