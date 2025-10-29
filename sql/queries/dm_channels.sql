-- name: CreateDMChannel :one
INSERT INTO dm_channels (
    name, icon, owner_id, is_group
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetDMChannelByID :one
SELECT * FROM dm_channels
WHERE id = $1 LIMIT 1;

-- name: UpdateDMChannel :one
UPDATE dm_channels
SET 
    name = COALESCE(sqlc.narg('name'), name),
    icon = COALESCE(sqlc.narg('icon'), icon),
    last_message_id = COALESCE(sqlc.narg('last_message_id'), last_message_id),
    last_message_at = COALESCE(sqlc.narg('last_message_at'), last_message_at)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteDMChannel :exec
DELETE FROM dm_channels
WHERE id = $1;

-- name: GetUserDMChannels :many
SELECT dc.* FROM dm_channels dc
INNER JOIN dm_participants dp ON dc.id = dp.dm_channel_id
WHERE dp.user_id = $1
ORDER BY dc.last_message_at DESC NULLS LAST;
