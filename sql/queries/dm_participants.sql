-- name: AddDMParticipant :one
INSERT INTO dm_participants (
    dm_channel_id, user_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetDMParticipants :many
SELECT * FROM dm_participants
WHERE dm_channel_id = $1;

-- name: RemoveDMParticipant :exec
DELETE FROM dm_participants
WHERE dm_channel_id = $1 AND user_id = $2;

-- name: UpdateLastReadMessage :exec
UPDATE dm_participants
SET last_read_message_id = $3
WHERE dm_channel_id = $1 AND user_id = $2;

-- name: GetDMChannelForUsers :one
SELECT dc.* FROM dm_channels dc
INNER JOIN dm_participants dp1 ON dc.id = dp1.dm_channel_id
INNER JOIN dm_participants dp2 ON dc.id = dp2.dm_channel_id
WHERE dp1.user_id = $1 AND dp2.user_id = $2 AND dc.is_group = FALSE
LIMIT 1;
