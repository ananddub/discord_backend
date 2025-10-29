-- name: CreateServer :one
INSERT INTO servers (
    name, icon, banner, description, owner_id, region
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetServerByID :one
SELECT * FROM servers
WHERE id = $1 AND is_deleted = FALSE LIMIT 1;

-- name: UpdateServer :one
UPDATE servers
SET 
    name = COALESCE(sqlc.narg('name'), name),
    icon = COALESCE(sqlc.narg('icon'), icon),
    banner = COALESCE(sqlc.narg('banner'), banner),
    description = COALESCE(sqlc.narg('description'), description),
    region = COALESCE(sqlc.narg('region'), region),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id') AND is_deleted = FALSE
RETURNING *;

-- name: SoftDeleteServer :exec
UPDATE servers
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: HardDeleteServer :exec
DELETE FROM servers
WHERE id = $1;

-- name: RestoreServer :exec
UPDATE servers
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetUserServers :many
SELECT s.* FROM servers s
INNER JOIN server_members sm ON s.id = sm.server_id
WHERE sm.user_id = $1 AND s.is_deleted = FALSE
ORDER BY sm.joined_at DESC;

-- name: IncrementMemberCount :exec
UPDATE servers
SET member_count = member_count + 1, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_deleted = FALSE;

-- name: DecrementMemberCount :exec
UPDATE servers
SET member_count = member_count - 1, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_deleted = FALSE;

-- name: GetServersByOwner :many
SELECT * FROM servers
WHERE owner_id = $1 AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: UpdateServerOwner :exec
UPDATE servers
SET owner_id = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND is_deleted = FALSE;
