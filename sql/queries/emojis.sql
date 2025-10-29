-- name: CreateEmoji :one
INSERT INTO emojis (
    server_id, name, image_url, creator_id, require_colons, animated
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetEmojiByID :one
SELECT * FROM emojis
WHERE id = $1 LIMIT 1;

-- name: GetServerEmojis :many
SELECT * FROM emojis
WHERE server_id = $1 AND available = TRUE
ORDER BY name;

-- name: UpdateEmoji :one
UPDATE emojis
SET 
    name = COALESCE(sqlc.narg('name'), name),
    available = COALESCE(sqlc.narg('available'), available)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEmoji :exec
DELETE FROM emojis
WHERE id = $1;

-- name: SearchEmojis :many
SELECT * FROM emojis
WHERE server_id = $1 AND name ILIKE '%' || $2 || '%' AND available = TRUE
ORDER BY name
LIMIT $3;
