-- name: CreateReaction :one
INSERT INTO message_reactions (
    message_id, user_id, emoji, emoji_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetMessageReactions :many
SELECT * FROM message_reactions
WHERE message_id = $1;

-- name: GetReactionsByEmoji :many
SELECT * FROM message_reactions
WHERE message_id = $1 AND emoji = $2;

-- name: DeleteReaction :exec
DELETE FROM message_reactions
WHERE message_id = $1 AND user_id = $2 AND emoji = $3;

-- name: DeleteAllReactions :exec
DELETE FROM message_reactions
WHERE message_id = $1;

-- name: CountReactions :one
SELECT emoji, COUNT(*) as count
FROM message_reactions
WHERE message_id = $1
GROUP BY emoji;
