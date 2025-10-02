-- name: CreateMessage :one
INSERT INTO messages (user_id, sender_id, reply_to_message_id, content) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages WHERE id = $1;

-- name: GetUserMessages :many
SELECT * FROM messages WHERE user_id = $1 ORDER BY created_at DESC;

-- name: MarkMessageAsRead :exec
UPDATE messages SET is_read = TRUE WHERE id = $1;
