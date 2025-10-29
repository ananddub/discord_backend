-- name: SyncFriends :many
-- Get all friends updated after last_updated_at timestamp
SELECT 
    f.id,
    f.user_id,
    f.friend_id,
    f.status,
    f.alias_name,
    f.is_favorite,
    f.created_at,
    f.updated_at,
    u.id as friend_user_id,
    u.username as friend_username,
    u.email as friend_email,
    u.full_name as friend_full_name,
    u.profile_pic as friend_profile_pic,
    u.bio as friend_bio,
    u.color_code as friend_color_code,
    u.background_color as friend_background_color,
    u.background_pic as friend_background_pic,
    u.status as friend_status,
    u.custom_status as friend_custom_status,
    u.is_bot as friend_is_bot,
    u.is_verified as friend_is_verified,
    u.created_at as friend_created_at,
    u.updated_at as friend_updated_at
FROM friends f
LEFT JOIN users u ON f.friend_id = u.id
WHERE f.user_id = $1 
    AND f.updated_at > $2
ORDER BY f.updated_at DESC
LIMIT $3 OFFSET $4;

-- name: SyncPendingFriendRequests :many
-- Get pending friend requests updated after last_updated_at
SELECT 
    f.id,
    f.user_id,
    f.friend_id,
    f.status,
    f.alias_name,
    f.is_favorite,
    f.created_at,
    f.updated_at,
    u.id as requester_id,
    u.username as requester_username,
    u.profile_pic as requester_profile_pic,
    u.status as requester_status
FROM friends f
LEFT JOIN users u ON f.user_id = u.id
WHERE f.friend_id = $1 
    AND f.status = 'pending'
    AND f.updated_at > $2
ORDER BY f.created_at DESC
LIMIT $3;

-- name: SyncDeletedFriends :many
-- Get deleted friendships (track separately in audit or use soft delete)
-- For now, returning empty - implement soft delete if needed
SELECT id FROM friends WHERE 1=0;

-- name: SyncMessages :many
-- Get messages updated after last_updated_at
SELECT 
    m.id,
    m.channel_id,
    m.sender_id,
    m.content,
    m.message_type,
    m.reply_to_message_id,
    m.is_edited,
    m.is_pinned,
    m.mention_everyone,
    m.created_at,
    m.updated_at,
    m.edited_at,
    u.username as sender_username,
    u.profile_pic as sender_profile_pic
FROM messages m
LEFT JOIN users u ON m.sender_id = u.id
WHERE m.channel_id = $1
    AND m.updated_at > $2
ORDER BY m.created_at DESC
LIMIT $3 OFFSET $4;

-- name: SyncUserMessages :many
-- Get all messages for channels user has access to
SELECT 
    m.id,
    m.channel_id,
    m.sender_id,
    m.content,
    m.message_type,
    m.reply_to_message_id,
    m.is_edited,
    m.is_pinned,
    m.mention_everyone,
    m.created_at,
    m.updated_at,
    m.edited_at
FROM messages m
INNER JOIN channels c ON m.channel_id = c.id
INNER JOIN server_members sm ON c.server_id = sm.server_id
WHERE sm.user_id = $1
    AND m.updated_at > $2
ORDER BY m.created_at DESC
LIMIT $3 OFFSET $4;

-- name: SyncServers :many
-- Get servers user is member of, updated after last_updated_at
SELECT 
    s.id,
    s.name,
    s.icon,
    s.banner,
    s.description,
    s.owner_id,
    s.region,
    s.member_count,
    s.is_verified,
    s.vanity_url,
    s.created_at,
    s.updated_at
FROM servers s
INNER JOIN server_members sm ON s.id = sm.server_id
WHERE sm.user_id = $1
    AND s.updated_at > $2
ORDER BY s.updated_at DESC;

-- name: SyncChannels :many
-- Get channels from user's servers, updated after last_updated_at
SELECT 
    c.id,
    c.server_id,
    c.category_id,
    c.name,
    c.type,
    c.position,
    c.topic,
    c.is_nsfw,
    c.slowmode_delay,
    c.user_limit,
    c.bitrate,
    c.is_private,
    c.created_at,
    c.updated_at
FROM channels c
INNER JOIN server_members sm ON c.server_id = sm.server_id
WHERE sm.user_id = $1
    AND c.updated_at > $2
ORDER BY c.server_id, c.position;

-- name: SyncUserProfile :one
-- Get user profile if updated after last_updated_at
SELECT 
    id,
    username,
    email,
    full_name,
    profile_pic,
    bio,
    color_code,
    background_color,
    background_pic,
    status,
    custom_status,
    is_bot,
    is_verified,
    is_2fa_enabled,
    created_at,
    updated_at
FROM users
WHERE id = $1
    AND updated_at > $2;

-- name: SyncVoiceChannels :many
-- Get voice channels from user's servers, updated after last_updated_at
SELECT 
    c.id,
    c.server_id,
    c.category_id,
    c.name,
    c.type,
    c.position,
    c.topic,
    c.user_limit,
    c.bitrate,
    c.is_private,
    c.created_at,
    c.updated_at
FROM channels c
INNER JOIN server_members sm ON c.server_id = sm.server_id
WHERE sm.user_id = $1
    AND c.type IN ('voice', 'stage')
    AND c.updated_at > $2
ORDER BY c.server_id, c.position;

-- name: SyncTextChannels :many
-- Get text channels from user's servers, updated after last_updated_at
SELECT 
    c.id,
    c.server_id,
    c.category_id,
    c.name,
    c.type,
    c.position,
    c.topic,
    c.is_nsfw,
    c.slowmode_delay,
    c.is_private,
    c.created_at,
    c.updated_at
FROM channels c
INNER JOIN server_members sm ON c.server_id = sm.server_id
WHERE sm.user_id = $1
    AND c.type IN ('text', 'announcement', 'forum')
    AND c.updated_at > $2
ORDER BY c.server_id, c.position;

-- name: SyncDirectMessages :many
-- Get DM channels user is part of
SELECT 
    dc.id,
    dc.name,
    dc.icon,
    dc.owner_id,
    dc.is_group,
    dc.last_message_id,
    dc.last_message_at,
    dc.created_at
FROM dm_channels dc
INNER JOIN dm_participants dp ON dc.id = dp.dm_channel_id
WHERE dp.user_id = $1
    AND (dc.last_message_at > $2 OR dc.created_at > $2)
ORDER BY dc.last_message_at DESC NULLS LAST
LIMIT $3 OFFSET $4;

-- name: SyncPermissions :many
-- Get channel permissions updated after last_updated_at
SELECT 
    cp.id,
    cp.channel_id,
    cp.role_id,
    cp.user_id,
    cp.allow_permissions,
    cp.deny_permissions,
    cp.created_at
FROM channel_permissions cp
INNER JOIN channels c ON cp.channel_id = c.id
INNER JOIN server_members sm ON c.server_id = sm.server_id
WHERE sm.user_id = $1
    AND cp.created_at > $2
ORDER BY cp.created_at DESC;

-- name: SyncVoiceStates :many
-- Get active voice states for user's servers
SELECT 
    vs.id,
    vs.user_id,
    vs.channel_id,
    vs.server_id,
    vs.session_id,
    vs.is_muted,
    vs.is_deafened,
    vs.self_mute,
    vs.self_deaf,
    vs.self_video,
    vs.self_stream,
    vs.suppress,
    vs.joined_at
FROM voice_states vs
INNER JOIN server_members sm ON vs.server_id = sm.server_id
WHERE sm.user_id = $1
    AND vs.joined_at > $2
ORDER BY vs.joined_at DESC;

-- name: SyncMessageAttachments :many
-- Get message attachments for synced messages
SELECT 
    ma.id,
    ma.message_id,
    ma.file_name,
    ma.file_url,
    ma.file_size,
    ma.file_type,
    ma.width,
    ma.height,
    ma.created_at
FROM message_attachments ma
WHERE ma.message_id = ANY($1::int[])
ORDER BY ma.message_id, ma.id;

-- name: SyncMessageReactions :many
-- Get message reactions for synced messages
SELECT 
    mr.id,
    mr.message_id,
    mr.user_id,
    mr.emoji,
    mr.created_at
FROM message_reactions mr
WHERE mr.message_id = ANY($1::int[])
ORDER BY mr.message_id, mr.created_at;

-- name: GetServerTimestamp :one
-- Get current server timestamp for sync
SELECT EXTRACT(EPOCH FROM NOW()) * 1000 as server_timestamp;

-- name: CountUpdatedFriends :one
-- Count friends updated after timestamp
SELECT COUNT(*) as total
FROM friends
WHERE user_id = $1 
    AND updated_at > $2;

-- name: CountUpdatedMessages :one
-- Count messages updated after timestamp
SELECT COUNT(*) as total
FROM messages m
INNER JOIN channels c ON m.channel_id = c.id
INNER JOIN server_members sm ON c.server_id = sm.server_id
WHERE sm.user_id = $1
    AND m.updated_at > $2;

-- name: CountUpdatedServers :one
-- Count servers updated after timestamp
SELECT COUNT(*) as total
FROM servers s
INNER JOIN server_members sm ON s.id = sm.server_id
WHERE sm.user_id = $1
    AND s.updated_at > $2;

-- name: CountUpdatedChannels :one
-- Count channels updated after timestamp
SELECT COUNT(*) as total
FROM channels c
INNER JOIN server_members sm ON c.server_id = sm.server_id
WHERE sm.user_id = $1
    AND c.updated_at > $2;

-- name: GetUserLastUpdate :one
-- Get user's last profile update timestamp
SELECT updated_at
FROM users
WHERE id = $1;

-- name: SyncServerMembers :many
-- Get server members updated after last_updated_at
SELECT 
    sm.id,
    sm.server_id,
    sm.user_id,
    sm.nickname,
    sm.joined_at,
    u.username,
    u.profile_pic,
    u.status
FROM server_members sm
LEFT JOIN users u ON sm.user_id = u.id
WHERE sm.server_id = $1
    AND sm.joined_at > $2
ORDER BY sm.joined_at DESC;

-- name: SyncBans :many
-- Get server bans updated after last_updated_at
SELECT 
    b.id,
    b.server_id,
    b.user_id,
    b.moderator_id,
    b.reason,
    b.expires_at,
    b.created_at
FROM bans b
INNER JOIN server_members sm ON b.server_id = sm.server_id
WHERE sm.user_id = $1
    AND b.created_at > $2
ORDER BY b.created_at DESC;

-- name: SyncInvites :many
-- Get server invites updated after last_updated_at
SELECT 
    i.id,
    i.code,
    i.server_id,
    i.channel_id,
    i.inviter_id,
    i.max_uses,
    i.uses,
    i.max_age,
    i.temporary,
    i.created_at,
    i.expires_at
FROM invites i
INNER JOIN server_members sm ON i.server_id = sm.server_id
WHERE sm.user_id = $1
    AND i.created_at > $2
ORDER BY i.created_at DESC;
