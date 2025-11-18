-- name: UpsertUserPresence :one
INSERT INTO
    user_presence (
        user_id,
        status,
        custom_status,
        custom_status_emoji,
        activity,
        last_seen
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        CURRENT_TIMESTAMP
    )
ON CONFLICT (user_id) DO
UPDATE
SET
    status = EXCLUDED.status,
    custom_status = EXCLUDED.custom_status,
    custom_status_emoji = EXCLUDED.custom_status_emoji,
    activity = EXCLUDED.activity,
    last_seen = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
RETURNING
    *;

-- name: GetUserPresence :one
SELECT * FROM user_presence WHERE user_id = $1 LIMIT 1;

-- name: GetMultipleUserPresences :many
SELECT * FROM user_presence WHERE user_id = ANY ($1::int[]);

-- name: UpdatePresenceStatus :one
UPDATE user_presence
SET
    status = $2,
    last_seen = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
RETURNING
    *;

-- name: SetCustomStatus :one
UPDATE user_presence
SET
    custom_status = $2,
    custom_status_emoji = $3,
    custom_status_expires_at = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
RETURNING
    *;

-- name: ClearExpiredCustomStatuses :one
UPDATE user_presence
SET
    custom_status = NULL,
    custom_status_emoji = NULL,
    custom_status_expires_at = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE
    custom_status_expires_at IS NOT NULL
    AND custom_status_expires_at < CURRENT_TIMESTAMP
RETURNING
    *;
