-- name: CreateFriend :one
INSERT INTO friends (
    user_id, friend_id, status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetFriendship :one
SELECT * FROM friends
WHERE user_id = $1 AND friend_id = $2 AND is_deleted = FALSE LIMIT 1;

-- name: GetUserFriends :many
SELECT * FROM friends
WHERE user_id = $1 AND status = 'accepted' AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetPendingFriendRequests :many
SELECT * FROM friends
WHERE friend_id = $1 AND status = 'pending' AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetSentFriendRequests :many
SELECT * FROM friends
WHERE user_id = $1 AND status = 'pending' AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: UpdateFriendStatus :exec
UPDATE friends
SET status = $3, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2 AND is_deleted = FALSE;

-- name: UpdateFriendAlias :exec
UPDATE friends
SET alias_name = $3, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2 AND is_deleted = FALSE;

-- name: ToggleFavorite :exec
UPDATE friends
SET is_favorite = $3, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2 AND is_deleted = FALSE;

-- name: SoftDeleteFriendship :exec
UPDATE friends
SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1);

-- name: HardDeleteFriendship :exec
DELETE FROM friends
WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1);

-- name: RestoreFriendship :exec
UPDATE friends
SET is_deleted = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2;

-- name: BlockUser :exec
UPDATE friends
SET status = 'blocked', updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2 AND is_deleted = FALSE;

-- name: GetBlockedUsers :many
SELECT * FROM friends
WHERE user_id = $1 AND status = 'blocked';
