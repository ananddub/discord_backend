-- name: CreateRole :one
INSERT INTO roles (
    server_id, name, color, hoist, position, permissions, mentionable, description
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;

-- name: GetServerRoles :many
SELECT * FROM roles
WHERE server_id = $1
ORDER BY position DESC;

-- name: UpdateRole :one
UPDATE roles
SET 
    name = COALESCE(sqlc.narg('name'), name),
    color = COALESCE(sqlc.narg('color'), color),
    hoist = COALESCE(sqlc.narg('hoist'), hoist),
    position = COALESCE(sqlc.narg('position'), position),
    permissions = COALESCE(sqlc.narg('permissions'), permissions),
    mentionable = COALESCE(sqlc.narg('mentionable'), mentionable),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;

-- name: UpdateRolePosition :exec
UPDATE roles
SET position = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
