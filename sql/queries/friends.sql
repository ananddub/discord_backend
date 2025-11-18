-- name: CreateFriend :one
INSERT INTO
    friends (
        user_id,
        friend_id,
        is_pending,
        is_accepted
    )
VALUES ($1, $2, TRUE, FALSE)
RETURNING
    *;

-- name: GetFriendship :one
SELECT *
FROM friends
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_deleted = FALSE
LIMIT 1;

-- name: GetAcceptedFriends :many
SELECT *
FROM friends
WHERE
    user_id = $1
    AND is_accepted = TRUE
    AND is_blocked = FALSE
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetPendingFriendRequests :many
SELECT *
FROM friends
WHERE
    friend_id = $1
    AND is_pending = TRUE
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetSentFriendRequests :many
SELECT *
FROM friends
WHERE
    user_id = $1
    AND is_pending = TRUE
    AND is_accepted = FALSE
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetBlockedUsers :many
SELECT *
FROM friends
WHERE
    user_id = $1
    AND is_blocked = TRUE
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetMutedFriends :many
SELECT *
FROM friends
WHERE
    user_id = $1
    AND is_muted = TRUE
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: GetFavoriteFriends :many
SELECT *
FROM friends
WHERE
    user_id = $1
    AND is_favorite = TRUE
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: UpdateFriendAlias :one
UPDATE friends
SET
    alias_name = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SetFavoriteFlag :one
UPDATE friends
SET
    is_favorite = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SetMutedFlag :one
UPDATE friends
SET
    is_muted = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_deleted = FALSE
RETURNING
    *;

-- name: ToggleFavoriteFlag :one
UPDATE friends
SET
    is_favorite = NOT is_favorite,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_deleted = FALSE
RETURNING
    *;

-- name: ToggleMutedFlag :one
UPDATE friends
SET
    is_muted = NOT is_muted,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_deleted = FALSE
RETURNING
    *;

-- name: AcceptFriendRequest :one
UPDATE friends
SET
    is_pending = FALSE,
    is_accepted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_pending = TRUE
    AND is_deleted = FALSE
RETURNING
    *;

-- name: RejectFriendRequest :one
UPDATE friends
SET
    is_pending = FALSE,
    is_accepted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_pending = TRUE
    AND is_deleted = FALSE
RETURNING
    *;

-- name: BlockUser :one
UPDATE friends
SET
    is_blocked = TRUE,
    is_accepted = FALSE,
    is_muted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_deleted = FALSE
RETURNING
    *;

-- name: UnblockUser :one
UPDATE friends
SET
    is_blocked = FALSE,
    is_muted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
    AND is_blocked = TRUE
    AND is_deleted = FALSE
RETURNING
    *;

-- name: UpdateFriendFlags :one
UPDATE friends
SET
    is_favorite = COALESCE(
        sqlc.narg ('is_favorite'),
        is_favorite
    ),
    is_muted = COALESCE(
        sqlc.narg ('is_muted'),
        is_muted
    ),
    alias_name = COALESCE(
        sqlc.narg ('alias_name'),
        alias_name
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = sqlc.arg ('user_id')
    AND friend_id = sqlc.arg ('friend_id')
    AND is_deleted = FALSE
RETURNING
    *;

-- name: SoftDeleteFriendship :one
UPDATE friends
SET
    is_deleted = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE (
        user_id = $1
        AND friend_id = $2
    )
    OR (
        user_id = $2
        AND friend_id = $1
    )
RETURNING
    *;

-- name: HardDeleteFriendship :one
DELETE FROM friends
WHERE (
        user_id = $1
        AND friend_id = $2
    )
    OR (
        user_id = $2
        AND friend_id = $1
    )
RETURNING
    *;

-- name: RestoreFriendship :one
UPDATE friends
SET
    is_deleted = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND friend_id = $2
RETURNING
    *;

-- name: GetFriendsWithFlags :many
SELECT *
FROM friends
WHERE
    user_id = $1
    AND is_accepted = $2
    AND is_blocked = $3
    AND is_muted = $4
    AND is_favorite = $5
    AND is_deleted = FALSE
ORDER BY created_at DESC;

-- name: CountAcceptedFriends :one
SELECT COUNT(*) as total
FROM friends
WHERE
    user_id = $1
    AND is_accepted = TRUE
    AND is_blocked = FALSE
    AND is_deleted = FALSE;

-- name: CountBlockedUsers :one
SELECT COUNT(*) as total
FROM friends
WHERE
    user_id = $1
    AND is_blocked = TRUE
    AND is_deleted = FALSE;

-- name: CountPendingRequests :one
SELECT COUNT(*) as total
FROM friends
WHERE
    friend_id = $1
    AND is_pending = TRUE
    AND is_deleted = FALSE;
