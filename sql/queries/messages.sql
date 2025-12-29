-- name: CreateMessage :one
INSERT INTO
    messages (
        channel_id,
        sender_id,
        content,
        message_type,
        reply_to_message_id,
        mention_everyone
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetMessageByID :one
SELECT * FROM messages WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetChannelMessages :many
SELECT *
FROM messages
WHERE
    channel_id = $1
    AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetMessagesBefore :many
SELECT *
FROM messages
WHERE
    channel_id = $1
    AND id < $2
    AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $3;

-- name: GetMessagesAfter :many
SELECT *
FROM messages
WHERE
    channel_id = $1
    AND id > $2
    AND is_deleted = FALSE
ORDER BY created_at ASC
LIMIT $3;

-- name: UpdateMessage :one
UPDATE messages
SET
    content = $2,
    is_edited = TRUE,
    edited_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SoftDeleteMessage :one
UPDATE messages
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: HardDeleteMessage :one
DELETE FROM messages WHERE id = $1 RETURNING *;

-- name: RestoreMessage :one
UPDATE messages
SET
    is_deleted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: PinMessage :one
UPDATE messages
SET
    is_pinned = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: UnpinMessage :one
UPDATE messages
SET
    is_pinned = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: GetPinnedMessages :many
SELECT *
FROM messages
WHERE
    channel_id = $1
    AND is_pinned = TRUE
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: BulkSoftDeleteMessages :one
UPDATE messages
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = ANY ($1::int[])
RETURNING
    *;

-- name: BulkHardDeleteMessages :one
DELETE FROM messages WHERE id = ANY ($1::int[]) RETURNING *;

-- name: SearchMessages :many
SELECT *
FROM messages
WHERE
    channel_id = $1
    AND content ILIKE '%' || $2 || '%'
    AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: GetUserMessages :many
SELECT *
FROM messages
WHERE
    sender_id = $1
    AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: CreateChatMessage :one
INSERT INTO
    messages (
        receiver_id,
        sender_id,
        content,
        message_type,
        reply_to_message_id,
        mention_everyone,
        ischannel
    )
VALUES ($1, $2, $3, $4, $5, $6, FALSE)
RETURNING
    *;

-- name: GetChatMessageByID :one
SELECT * FROM messages WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: UpdateChatMessage :one
UPDATE messages
SET
    content = $2,
    is_edited = TRUE,
    edited_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SoftDeleteChatMessage :one
UPDATE messages
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: HardDeleteChatMessage :one
DELETE FROM messages WHERE id = $1 RETURNING *;

-- name: GetChatMessages :many
SELECT *
FROM messages
WHERE (
        (
            sender_id = $1
            AND receiver_id = $2
        )
        OR (
            sender_id = $2
            AND receiver_id = $1
        )
    )
    AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: GetChatMessagesAfter :many
SELECT *
FROM messages
WHERE (
        (
            sender_id = $1
            AND receiver_id = $2
        )
        OR (
            sender_id = $2
            AND receiver_id = $1
        )
    )
    AND id > $3
    AND is_deleted = FALSE
ORDER BY created_at ASC
LIMIT $4;

-- name: SearchChatMessages :many
SELECT *
FROM messages
WHERE (
        (
            sender_id = $1
            AND receiver_id = $2
        )
        OR (
            sender_id = $2
            AND receiver_id = $1
        )
    )
    AND content ILIKE '%' || $3 || '%'
    AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $4
OFFSET
    $5;

-- name: GetChatMessagesBefore :many
SELECT *
FROM messages
WHERE (
        (
            sender_id = $1
            AND receiver_id = $2
        )
        OR (
            sender_id = $2
            AND receiver_id = $1
        )
    )
    AND id < $3
    AND is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $4;

-- name: HardDeleteChatMessages :one
DELETE FROM messages
WHERE (
        sender_id = $1
        AND receiver_id = $2
    )
    OR (
        sender_id = $2
        AND receiver_id = $1
    )
    AND is_deleted = TRUE
RETURNING
    *;