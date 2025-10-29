-- name: CreateMessage :one
INSERT INTO messages (
    channel_id, sender_id, content, message_type, reply_to_message_id, mention_everyone
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: GetChannelMessages :many
SELECT * FROM messages
WHERE channel_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetMessagesBefore :many
SELECT * FROM messages
WHERE channel_id = $1 AND id < $2
ORDER BY created_at DESC
LIMIT $3;

-- name: GetMessagesAfter :many
SELECT * FROM messages
WHERE channel_id = $1 AND id > $2
ORDER BY created_at ASC
LIMIT $3;

-- name: UpdateMessage :one
UPDATE messages
SET 
    content = $2,
    is_edited = TRUE,
    edited_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;

-- name: PinMessage :exec
UPDATE messages
SET is_pinned = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UnpinMessage :exec
UPDATE messages
SET is_pinned = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetPinnedMessages :many
SELECT * FROM messages
WHERE channel_id = $1 AND is_pinned = TRUE
ORDER BY created_at DESC;

-- name: BulkDeleteMessages :exec
DELETE FROM messages
WHERE id = ANY($1::int[]);

-- name: SearchMessages :many
SELECT * FROM messages
WHERE channel_id = $1 
  AND content ILIKE '%' || $2 || '%'
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetUserMessages :many
SELECT * FROM messages
WHERE sender_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
