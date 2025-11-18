-- name: CreateUser :one
INSERT INTO
    users (
        username,
        email,
        password,
        full_name,
        profile_pic,
        bio
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND is_deleted = FALSE LIMIT 1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE
    username = $1
    AND is_deleted = FALSE
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE(
        sqlc.narg ('full_name'),
        full_name
    ),
    profile_pic = COALESCE(
        sqlc.narg ('profile_pic'),
        profile_pic
    ),
    bio = COALESCE(sqlc.narg ('bio'), bio),
    color_code = COALESCE(
        sqlc.narg ('color_code'),
        color_code
    ),
    background_color = COALESCE(
        sqlc.narg ('background_color'),
        background_color
    ),
    background_pic = COALESCE(
        sqlc.narg ('background_pic'),
        background_pic
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg ('id')
    AND is_deleted = FALSE
RETURNING
    *;

-- name: UpdateUserStatus :one
UPDATE users
SET
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: UpdateUserPassword :one
UPDATE users
SET
    password = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SoftDeleteUser :one
UPDATE users
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: HardDeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING *;

-- name: RestoreUser :one
UPDATE users
SET
    is_deleted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
RETURNING
    *;

-- name: ListUsers :many
SELECT *
FROM users
WHERE
    is_deleted = FALSE
ORDER BY created_at DESC
LIMIT $1
OFFSET
    $2;

-- name: SearchUsers :many
SELECT *
FROM users
WHERE
    is_deleted = FALSE
    AND (
        username ILIKE '%' || $1 || '%'
        OR email ILIKE '%' || $1 || '%'
        OR full_name ILIKE '%' || $1 || '%'
    )
ORDER BY username
LIMIT $2
OFFSET
    $3;

-- name: Enable2FA :one
UPDATE users
SET
    is_2fa_enabled = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND is_deleted = FALSE
RETURNING
    *;

-- name: ConnectedUser :many
WITH
    f AS (
        SELECT
            CASE
                WHEN fr.user_id = $1 THEN fr.friend_id
                ELSE fr.user_id
            END AS id
        FROM friends fr
        WHERE (
                fr.user_id = $1
                OR fr.friend_id = $1
            )
            AND fr.status != 'blocked'
    ),
    blocked AS (
        SELECT
            CASE
                WHEN fr.user_id = $1 THEN fr.friend_id
                ELSE fr.user_id
            END AS id
        FROM friends fr
        WHERE (
                fr.user_id = $1
                OR fr.friend_id = $1
            )
            AND fr.status = 'blocked'
    ),
    c AS (
        SELECT
            CASE
                WHEN m.sender_id = $1 THEN m.receiver_id
                ELSE m.sender_id
            END AS id
        FROM messages m
        WHERE (
                m.sender_id = $1
                OR m.receiver_id = $1
            )
            AND m.sender_id != m.receiver_id
            AND (
                CASE
                    WHEN m.sender_id = $1 THEN m.receiver_id
                    ELSE m.sender_id
                END
            ) NOT IN (
                SELECT id
                FROM blocked
            )
    ),
    q AS (
        SELECT DISTINCT
            id
        FROM f
        UNION
        SELECT DISTINCT
            id
        FROM c
    )
SELECT u.*
FROM users u
WHERE
    u.id IN (
        SELECT id
        FROM q
    )
    AND u.is_deleted = FALSE;
