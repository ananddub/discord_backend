-- name: CreateFriend :one
INSERT INTO friends (
    user_id, friend_id, status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetFriendship :one
SELECT * FROM friends
WHERE user_id = $1 AND friend_id = $2 LIMIT 1;

-- name: GetUserFriends :many
SELECT * FROM friends
WHERE user_id = $1 AND status = 'accepted'
ORDER BY created_at DESC;

-- name: GetPendingFriendRequests :many
SELECT * FROM friends
WHERE friend_id = $1 AND status = 'pending'
ORDER BY created_at DESC;

-- name: GetSentFriendRequests :many
SELECT * FROM friends
WHERE user_id = $1 AND status = 'pending'
ORDER BY created_at DESC;

-- name: UpdateFriendStatus :exec
UPDATE friends
SET status = $3, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2;

-- name: UpdateFriendAlias :exec
UPDATE friends
SET alias_name = $3, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2;

-- name: ToggleFavorite :exec
UPDATE friends
SET is_favorite = $3, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2;

-- name: DeleteFriendship :exec
DELETE FROM friends
WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1);

-- name: BlockUser :exec
UPDATE friends
SET status = 'blocked', updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND friend_id = $2;

-- name: GetBlockedUsers :many
SELECT * FROM friends
WHERE user_id = $1 AND status = 'blocked';
