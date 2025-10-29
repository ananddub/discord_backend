-- name: CreateDMChannel :one
INSERT INTO dm_channels (
    name, icon, owner_id, is_group
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetDMChannelByID :one
SELECT * FROM dm_channels
WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: UpdateDMChannel :one
UPDATE dm_channels
SET 
    name = COALESCE(sqlc.narg('name'), name),
    icon = COALESCE(sqlc.narg('icon'), icon),
    last_message_id = COALESCE(sqlc.narg('last_message_id'), last_message_id),
    last_message_at = COALESCE(sqlc.narg('last_message_at'), last_message_at),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id') AND is_deleted = FALSE
RETURNING *;

-- name: SoftDeleteDMChannel :exec
UPDATE dm_channels
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: HardDeleteDMChannel :exec
DELETE FROM dm_channels
WHERE id = $1;

-- name: RestoreDMChannel :exec
UPDATE dm_channels
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetUserDMChannels :many
SELECT dc.* FROM dm_channels dc
INNER JOIN dm_participants dp ON dc.id = dp.dm_channel_id
WHERE dp.user_id = $1 AND dc.is_deleted = FALSE
ORDER BY dc.last_message_at DESC NULLS LAST;
