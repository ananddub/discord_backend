-- name: CreateMessage :one
INSERT INTO messages (
    channel_id, sender_id, content, message_type, reply_to_message_id, mention_everyone
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM messages
WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetChannelMessages :many
SELECT * FROM messages
WHERE channel_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetMessagesBefore :many
SELECT * FROM messages
WHERE channel_id = $1 AND id < $2 AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $3;

-- name: GetMessagesAfter :many
SELECT * FROM messages
WHERE channel_id = $1 AND id > $2 AND is_deleted = FALSE
ORDER BY created_at ASC
LIMIT $3;

-- name: UpdateMessage :one
UPDATE messages
SET 
    content = $2,
    is_edited = TRUE,
    edited_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_deleted = FALSE
RETURNING *;

-- name: SoftDeleteMessage :exec
UPDATE messages
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: HardDeleteMessage :exec
DELETE FROM messages
WHERE id = $1;

-- name: RestoreMessage :exec
UPDATE messages
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: PinMessage :exec
UPDATE messages
SET is_pinned = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_deleted = FALSE;

-- name: UnpinMessage :exec
UPDATE messages
SET is_pinned = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_deleted = FALSE;

-- name: GetPinnedMessages :many
SELECT * FROM messages
WHERE channel_id = $1 AND is_pinned = TRUE AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: BulkSoftDeleteMessages :exec
UPDATE messages
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = ANY($1::int[]);

-- name: BulkHardDeleteMessages :exec
DELETE FROM messages
WHERE id = ANY($1::int[]);

-- name: SearchMessages :many
SELECT * FROM messages
WHERE channel_id = $1 
  AND content ILIKE '%' || $2 || '%'
  AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetUserMessages :many
SELECT * FROM messages
WHERE sender_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
