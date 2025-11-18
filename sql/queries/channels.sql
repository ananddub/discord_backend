-- name: CreateChannel :one
INSERT INTO
    channels (
        server_id,
        category_id,
        name,
        type,
        position,
        topic,
        is_nsfw,
        slowmode_delay
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
RETURNING
    *;

-- name: GetChannelByID :one
SELECT * FROM channels WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetServerChannels :many
SELECT *
FROM channels
WHERE
    server_id = $1
    AND is_deleted = FALSE
ORDER BY position ASC;

-- name: GetChannelsByCategory :many
SELECT *
FROM channels
WHERE
    category_id = $1
    AND is_deleted = FALSE
ORDER BY position ASC;

-- name: UpdateChannel :one
UPDATE channels
SET
    name = COALESCE(sqlc.narg ('name'), name),
    topic = COALESCE(sqlc.narg ('topic'), topic),
    position = COALESCE(
        sqlc.narg ('position'),
        position
    ),
    is_nsfw = COALESCE(
        sqlc.narg ('is_nsfw'),
        is_nsfw
    ),
    slowmode_delay = COALESCE(
        sqlc.narg ('slowmode_delay'),
        slowmode_delay
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg ('id')
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SoftDeleteChannel :one
UPDATE channels
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: HardDeleteChannel :one
DELETE FROM channels WHERE id = $1 RETURNING *;

-- name: RestoreChannel :one
UPDATE channels
SET
    is_deleted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateChannelPosition :one
UPDATE channels
SET
    position = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: GetChannelsByType :many
SELECT *
FROM channels
WHERE
    server_id = $1
    AND type = $2
    AND is_deleted = FALSE
ORDER BY position ASC;
