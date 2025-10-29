-- name: CreateVoiceState :one
INSERT INTO voice_states (
    user_id, channel_id, server_id, session_id, is_muted, is_deafened, self_mute, self_deaf
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetVoiceState :one
SELECT * FROM voice_states
WHERE user_id = $1 AND channel_id = $2 LIMIT 1;

-- name: GetChannelVoiceStates :many
SELECT * FROM voice_states
WHERE channel_id = $1;

-- name: GetUserVoiceState :one
SELECT * FROM voice_states
WHERE user_id = $1 LIMIT 1;

-- name: UpdateVoiceState :one
UPDATE voice_states
SET 
    is_muted = COALESCE(sqlc.narg('is_muted'), is_muted),
    is_deafened = COALESCE(sqlc.narg('is_deafened'), is_deafened),
    self_mute = COALESCE(sqlc.narg('self_mute'), self_mute),
    self_deaf = COALESCE(sqlc.narg('self_deaf'), self_deaf),
    self_video = COALESCE(sqlc.narg('self_video'), self_video),
    self_stream = COALESCE(sqlc.narg('self_stream'), self_stream),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.arg('user_id') AND channel_id = sqlc.arg('channel_id')
RETURNING *;

-- name: DeleteVoiceState :exec
DELETE FROM voice_states
WHERE user_id = $1 AND channel_id = $2;

-- name: DeleteUserVoiceStates :exec
DELETE FROM voice_states
WHERE user_id = $1;
