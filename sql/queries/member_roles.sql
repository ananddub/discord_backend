-- name: AssignRoleToMember :one
INSERT INTO member_roles (
    member_id, role_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: RemoveRoleFromMember :exec
DELETE FROM member_roles
WHERE member_id = $1 AND role_id = $2;

-- name: GetMemberRoles :many
SELECT r.* FROM roles r
INNER JOIN member_roles mr ON r.id = mr.role_id
WHERE mr.member_id = $1
ORDER BY r.position DESC;

-- name: GetRoleMembers :many
SELECT sm.* FROM server_members sm
INNER JOIN member_roles mr ON sm.id = mr.member_id
WHERE mr.role_id = $1;

-- name: RemoveAllMemberRoles :exec
DELETE FROM member_roles
WHERE member_id = $1;
