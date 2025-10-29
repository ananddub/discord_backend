-- name: CreateMessageAttachment :one
INSERT INTO message_attachments (
    message_id, file_url, file_name, file_type, file_size, width, height
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMessageAttachments :many
SELECT * FROM message_attachments
WHERE message_id = $1;

-- name: DeleteMessageAttachment :exec
DELETE FROM message_attachments
WHERE id = $1;

-- name: DeleteMessageAttachments :exec
DELETE FROM message_attachments
WHERE message_id = $1;
