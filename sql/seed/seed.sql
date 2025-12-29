-- Seed data for Discord Backend

-- Users
INSERT INTO users (username, email, password, full_name, status) VALUES
('john_doe', 'john@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'John Doe', 'online'),
('jane_smith', 'jane@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'Jane Smith', 'idle'),
('bob_wilson', 'bob@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'Bob Wilson', 'offline'),
('alice_brown', 'alice@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'Alice Brown', 'dnd'),
('charlie_davis', 'charlie@example.com', '$2a$10$abcdefghijklmnopqrstuv', 'Charlie Davis', 'online')
RETURNING id;

-- Servers (using actual user IDs)
INSERT INTO servers (name, owner_id, member_count) 
SELECT 'Gaming Hub', id, 150 FROM users WHERE username = 'john_doe'
UNION ALL
SELECT 'Tech Talk', id, 80 FROM users WHERE username = 'jane_smith'
UNION ALL
SELECT 'Music Lovers', id, 200 FROM users WHERE username = 'bob_wilson'
UNION ALL
SELECT 'Study Group', id, 45 FROM users WHERE username = 'alice_brown';

-- Channels
INSERT INTO channels (server_id, name, type, topic)
SELECT s.id, 'general', 'text', 'General discussion' FROM servers s WHERE s.name = 'Gaming Hub'
UNION ALL
SELECT s.id, 'gaming', 'text', 'Gaming chat' FROM servers s WHERE s.name = 'Gaming Hub'
UNION ALL
SELECT s.id, 'voice-1', 'voice', 'Voice channel' FROM servers s WHERE s.name = 'Gaming Hub'
UNION ALL
SELECT s.id, 'tech-news', 'text', 'Latest tech news' FROM servers s WHERE s.name = 'Tech Talk'
UNION ALL
SELECT s.id, 'help', 'text', 'Get help here' FROM servers s WHERE s.name = 'Tech Talk'
UNION ALL
SELECT s.id, 'music-share', 'text', 'Share your music' FROM servers s WHERE s.name = 'Music Lovers'
UNION ALL
SELECT s.id, 'study-room', 'text', 'Study together' FROM servers s WHERE s.name = 'Study Group';

-- Messages
INSERT INTO messages (channel_id, sender_id, content, message_type)
SELECT c.id, u.id, 'Hello everyone!', 'default' 
FROM channels c, users u 
WHERE c.name = 'general' AND u.username = 'john_doe'
UNION ALL
SELECT c.id, u.id, 'Hey! How are you?', 'default'
FROM channels c, users u
WHERE c.name = 'general' AND u.username = 'jane_smith'
UNION ALL
SELECT c.id, u.id, 'Anyone up for a game?', 'default'
FROM channels c, users u
WHERE c.name = 'gaming' AND u.username = 'bob_wilson';

-- Friends
INSERT INTO friends (user_id, friend_id, is_accepted, is_pending)
SELECT u1.id, u2.id, true, false
FROM users u1, users u2
WHERE u1.username = 'john_doe' AND u2.username = 'jane_smith'
UNION ALL
SELECT u1.id, u2.id, true, false
FROM users u1, users u2
WHERE u1.username = 'john_doe' AND u2.username = 'bob_wilson'
UNION ALL
SELECT u1.id, u2.id, false, true
FROM users u1, users u2
WHERE u1.username = 'bob_wilson' AND u2.username = 'charlie_davis';

-- Server Members
INSERT INTO server_members (server_id, user_id)
SELECT s.id, u.id FROM servers s, users u WHERE s.name = 'Gaming Hub' AND u.username = 'john_doe'
UNION ALL
SELECT s.id, u.id FROM servers s, users u WHERE s.name = 'Gaming Hub' AND u.username = 'jane_smith'
UNION ALL
SELECT s.id, u.id FROM servers s, users u WHERE s.name = 'Tech Talk' AND u.username = 'jane_smith'
UNION ALL
SELECT s.id, u.id FROM servers s, users u WHERE s.name = 'Music Lovers' AND u.username = 'bob_wilson'
UNION ALL
SELECT s.id, u.id FROM servers s, users u WHERE s.name = 'Study Group' AND u.username = 'alice_brown';
