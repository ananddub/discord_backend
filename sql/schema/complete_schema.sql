-- Discord Backend - Complete Database Schema
-- This file contains the complete schema for reference
-- For migrations, use the files in the migrations/ directory

-- ==============================================
-- USERS & AUTHENTICATION
-- ==============================================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    profile_pic VARCHAR(255),
    bio TEXT,
    color_code VARCHAR(7) DEFAULT '#5865F2',
    background_color VARCHAR(7) DEFAULT '#313338',
    background_pic VARCHAR(255),
    status VARCHAR(50) DEFAULT 'offline' NOT NULL,
    custom_status VARCHAR(255),
    is_bot BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    is_2fa_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_status ON users(status);

-- ==============================================
-- SERVERS (GUILDS)
-- ==============================================

CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(255),
    banner VARCHAR(255),
    description TEXT,
    owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    region VARCHAR(50) DEFAULT 'auto',
    member_count INTEGER DEFAULT 1,
    is_verified BOOLEAN DEFAULT FALSE,
    vanity_url VARCHAR(50) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_servers_owner_id ON servers(owner_id);
CREATE INDEX idx_servers_vanity_url ON servers(vanity_url);

-- ==============================================
-- ROLES & PERMISSIONS
-- ==============================================

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) DEFAULT '#99AAB5',
    hoist BOOLEAN DEFAULT FALSE,
    position INTEGER DEFAULT 0,
    permissions BIGINT DEFAULT 0,
    mentionable BOOLEAN DEFAULT FALSE,
    icon VARCHAR(255),
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(server_id, name)
);

CREATE INDEX idx_roles_server_id ON roles(server_id);
CREATE INDEX idx_roles_position ON roles(server_id, position);

-- ==============================================
-- SERVER MEMBERS
-- ==============================================

CREATE TABLE server_members (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nickname VARCHAR(100),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_muted BOOLEAN DEFAULT FALSE,
    is_deafened BOOLEAN DEFAULT FALSE,
    UNIQUE(server_id, user_id)
);

CREATE INDEX idx_server_members_server_id ON server_members(server_id);
CREATE INDEX idx_server_members_user_id ON server_members(user_id);
CREATE INDEX idx_server_members_joined_at ON server_members(joined_at);

CREATE TABLE member_roles (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL REFERENCES server_members(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(member_id, role_id)
);

CREATE INDEX idx_member_roles_member_id ON member_roles(member_id);
CREATE INDEX idx_member_roles_role_id ON member_roles(role_id);

-- ==============================================
-- CHANNELS
-- ==============================================

CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES channels(id) ON DELETE SET NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'voice', 'category', 'announcement', 'stage', 'forum', 'dm', 'group_dm')),
    position INTEGER DEFAULT 0,
    topic VARCHAR(1024),
    is_nsfw BOOLEAN DEFAULT FALSE,
    slowmode_delay INTEGER DEFAULT 0,
    user_limit INTEGER DEFAULT 0,
    bitrate INTEGER DEFAULT 64000,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_channels_server_id ON channels(server_id);
CREATE INDEX idx_channels_category_id ON channels(category_id);
CREATE INDEX idx_channels_type ON channels(type);
CREATE INDEX idx_channels_position ON channels(server_id, position);

CREATE TABLE channel_permissions (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    allow_permissions BIGINT DEFAULT 0,
    deny_permissions BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CHECK ((role_id IS NOT NULL AND user_id IS NULL) OR (role_id IS NULL AND user_id IS NOT NULL)),
    UNIQUE(channel_id, role_id, user_id)
);

CREATE INDEX idx_channel_permissions_channel_id ON channel_permissions(channel_id);
CREATE INDEX idx_channel_permissions_role_id ON channel_permissions(role_id);
CREATE INDEX idx_channel_permissions_user_id ON channel_permissions(user_id);

-- ==============================================
-- MESSAGES
-- ==============================================

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    sender_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'default' CHECK (message_type IN ('default', 'reply', 'system', 'user_join', 'user_leave', 'call', 'pinned')),
    reply_to_message_id INTEGER REFERENCES messages(id) ON DELETE SET NULL,
    is_edited BOOLEAN DEFAULT FALSE,
    is_pinned BOOLEAN DEFAULT FALSE,
    mention_everyone BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    edited_at TIMESTAMP
);

CREATE INDEX idx_messages_channel_id ON messages(channel_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(channel_id, created_at DESC);
CREATE INDEX idx_messages_pinned ON messages(channel_id, is_pinned) WHERE is_pinned = TRUE;
CREATE INDEX idx_messages_reply_to ON messages(reply_to_message_id);

CREATE TABLE message_attachments (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size BIGINT NOT NULL,
    width INTEGER,
    height INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_message_attachments_message_id ON message_attachments(message_id);

CREATE TABLE message_reactions (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji VARCHAR(100) NOT NULL,
    emoji_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(message_id, user_id, emoji)
);

CREATE INDEX idx_message_reactions_message_id ON message_reactions(message_id);
CREATE INDEX idx_message_reactions_user_id ON message_reactions(user_id);

CREATE TABLE message_mentions (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CHECK ((user_id IS NOT NULL AND role_id IS NULL) OR (user_id IS NULL AND role_id IS NOT NULL))
);

CREATE INDEX idx_message_mentions_message_id ON message_mentions(message_id);
CREATE INDEX idx_message_mentions_user_id ON message_mentions(user_id);
CREATE INDEX idx_message_mentions_role_id ON message_mentions(role_id);

-- ==============================================
-- DIRECT MESSAGES
-- ==============================================

CREATE TABLE dm_channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    icon VARCHAR(255),
    owner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    is_group BOOLEAN DEFAULT FALSE,
    last_message_id INTEGER,
    last_message_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_dm_channels_owner_id ON dm_channels(owner_id);
CREATE INDEX idx_dm_channels_last_message_at ON dm_channels(last_message_at DESC);

CREATE TABLE dm_participants (
    id SERIAL PRIMARY KEY,
    dm_channel_id INTEGER NOT NULL REFERENCES dm_channels(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_read_message_id INTEGER,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(dm_channel_id, user_id)
);

CREATE INDEX idx_dm_participants_dm_channel_id ON dm_participants(dm_channel_id);
CREATE INDEX idx_dm_participants_user_id ON dm_participants(user_id);

-- ==============================================
-- FRIENDS
-- ==============================================

CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected', 'blocked')),
    alias_name VARCHAR(100),
    is_favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, friend_id),
    CHECK (user_id != friend_id)
);

CREATE INDEX idx_friends_user_id ON friends(user_id);
CREATE INDEX idx_friends_friend_id ON friends(friend_id);
CREATE INDEX idx_friends_status ON friends(user_id, status);

-- ==============================================
-- INVITES & BANS
-- ==============================================

CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    channel_id INTEGER REFERENCES channels(id) ON DELETE CASCADE,
    inviter_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    max_uses INTEGER DEFAULT 0,
    uses INTEGER DEFAULT 0,
    max_age INTEGER DEFAULT 0,
    temporary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP
);

CREATE INDEX idx_invites_code ON invites(code);
CREATE INDEX idx_invites_server_id ON invites(server_id);
CREATE INDEX idx_invites_inviter_id ON invites(inviter_id);
CREATE INDEX idx_invites_expires_at ON invites(expires_at);

CREATE TABLE bans (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    moderator_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    reason TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(server_id, user_id)
);

CREATE INDEX idx_bans_server_id ON bans(server_id);
CREATE INDEX idx_bans_user_id ON bans(user_id);
CREATE INDEX idx_bans_expires_at ON bans(expires_at);

-- ==============================================
-- EMOJIS
-- ==============================================

CREATE TABLE emojis (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    creator_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    require_colons BOOLEAN DEFAULT TRUE,
    managed BOOLEAN DEFAULT FALSE,
    animated BOOLEAN DEFAULT FALSE,
    available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(server_id, name)
);

CREATE INDEX idx_emojis_server_id ON emojis(server_id);
CREATE INDEX idx_emojis_name ON emojis(name);

-- ==============================================
-- VOICE
-- ==============================================

CREATE TABLE voice_states (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    server_id INTEGER REFERENCES servers(id) ON DELETE CASCADE,
    session_id VARCHAR(100) NOT NULL,
    is_muted BOOLEAN DEFAULT FALSE,
    is_deafened BOOLEAN DEFAULT FALSE,
    self_mute BOOLEAN DEFAULT FALSE,
    self_deaf BOOLEAN DEFAULT FALSE,
    self_video BOOLEAN DEFAULT FALSE,
    self_stream BOOLEAN DEFAULT FALSE,
    suppress BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE(user_id, channel_id)
);

CREATE INDEX idx_voice_states_user_id ON voice_states(user_id);
CREATE INDEX idx_voice_states_channel_id ON voice_states(channel_id);
CREATE INDEX idx_voice_states_session_id ON voice_states(session_id);

-- ==============================================
-- USER PRESENCE
-- ==============================================

CREATE TABLE user_presence (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    status VARCHAR(20) DEFAULT 'offline' CHECK (status IN ('online', 'idle', 'dnd', 'offline', 'invisible')),
    custom_status VARCHAR(255),
    custom_status_emoji VARCHAR(100),
    custom_status_expires_at TIMESTAMP,
    activity VARCHAR(255),
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_user_presence_user_id ON user_presence(user_id);
CREATE INDEX idx_user_presence_status ON user_presence(status);
CREATE INDEX idx_user_presence_last_seen ON user_presence(last_seen);

-- ==============================================
-- AUDIT LOGS
-- ==============================================

CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    server_id INTEGER NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    target_id INTEGER,
    target_type VARCHAR(50),
    changes JSONB,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_audit_logs_server_id ON audit_logs(server_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(server_id, created_at DESC);
