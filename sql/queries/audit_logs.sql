-- name: CreateAuditLog :one
INSERT INTO
    audit_logs (
        server_id,
        user_id,
        action,
        target_id,
        target_type,
        changes,
        reason
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: GetAuditLogs :many
SELECT *
FROM audit_logs
WHERE
    server_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetAuditLogsByAction :many
SELECT *
FROM audit_logs
WHERE
    server_id = $1
    AND action = $2
ORDER BY created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: GetAuditLogsByUser :many
SELECT *
FROM audit_logs
WHERE
    server_id = $1
    AND user_id = $2
ORDER BY created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: DeleteOldAuditLogs :one
DELETE FROM audit_logs WHERE created_at < $1 RETURNING *;
