-- name: CreateEmoji :one
INSERT INTO
    emojis (
        server_id,
        name,
        image_url,
        creator_id,
        require_colons,
        animated
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetEmojiByID :one
SELECT * FROM emojis WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetServerEmojis :many
SELECT *
FROM emojis
WHERE
    server_id = $1
    AND available = TRUE
    AND is_deleted = FALSE
ORDER BY name;

-- name: UpdateEmoji :one
UPDATE emojis
SET
    name = COALESCE(sqlc.narg ('name'), name),
    available = COALESCE(
        sqlc.narg ('available'),
        available
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg ('id')
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SoftDeleteEmoji :one
UPDATE emojis
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: HardDeleteEmoji :one
DELETE FROM emojis WHERE id = $1 RETURNING *;

-- name: RestoreEmoji :one
UPDATE emojis
SET
    is_deleted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: SearchEmojis :many
SELECT *
FROM emojis
WHERE
    server_id = $1
    AND name ILIKE '%' || $2 || '%'
    AND available = TRUE
    AND is_deleted = FALSE
ORDER BY name
LIMIT $3;
