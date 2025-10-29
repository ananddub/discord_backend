-- name: CreateMessageAttachment :one
INSERT INTO message_attachments (
    message_id, file_url, file_name, file_type, file_size, width, height
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMessageAttachments :many
SELECT * FROM message_attachments
WHERE message_id = $1 AND is_deleted = FALSE;

-- name: SoftDeleteMessageAttachment :exec
UPDATE message_attachments
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: HardDeleteMessageAttachment :exec
DELETE FROM message_attachments
WHERE id = $1;

-- name: RestoreMessageAttachment :exec
UPDATE message_attachments
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: SoftDeleteMessageAttachments :exec
UPDATE message_attachments
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE message_id = $1;

-- name: HardDeleteMessageAttachments :exec
DELETE FROM message_attachments
WHERE message_id = $1;
