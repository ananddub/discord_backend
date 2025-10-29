-- name: CreateUser :one
INSERT INTO users (
    username, email, password, full_name, profile_pic, bio
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    profile_pic = COALESCE(sqlc.narg('profile_pic'), profile_pic),
    bio = COALESCE(sqlc.narg('bio'), bio),
    color_code = COALESCE(sqlc.narg('color_code'), color_code),
    background_color = COALESCE(sqlc.narg('background_color'), background_color),
    background_pic = COALESCE(sqlc.narg('background_pic'), background_pic),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateUserStatus :exec
UPDATE users
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchUsers :many
SELECT * FROM users
WHERE username ILIKE '%' || $1 || '%'
   OR email ILIKE '%' || $1 || '%'
   OR full_name ILIKE '%' || $1 || '%'
ORDER BY username
LIMIT $2 OFFSET $3;

-- name: Enable2FA :exec
UPDATE users
SET is_2fa_enabled = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
