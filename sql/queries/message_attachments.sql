-- name: CreateMessageAttachment :one
INSERT INTO
    message_attachments (
        message_id,
        file_url,
        file_name,
        file_type,
        file_size,
        width,
        height
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: GetMessageAttachments :many
SELECT *
FROM message_attachments
WHERE
    message_id = $1
    AND is_deleted = FALSE;

-- name: SoftDeleteMessageAttachment :one
UPDATE message_attachments
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: HardDeleteMessageAttachment :one
DELETE FROM message_attachments WHERE id = $1 RETURNING *;

-- name: RestoreMessageAttachment :one
UPDATE message_attachments
SET
    is_deleted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: SoftDeleteMessageAttachments :one
UPDATE message_attachments
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    message_id = $1
RETURNING
    *;

-- name: HardDeleteMessageAttachments :one
DELETE FROM message_attachments WHERE message_id = $1 RETURNING *;
