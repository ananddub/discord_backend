-- name: CreateUser :one
INSERT INTO users (username, email,password, full_name, profile_pic, bio, color_code, background_color, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUserStatus :exec
UPDATE users SET status = $2, updated_at = NOW() WHERE id = $1;

