package repository

import (
	"context"
	"time"

	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Helper function to convert unix milliseconds to pgtype.Timestamp
func toPgTimestamp(unixMilli int64) pgtype.Timestamp {
	t := time.UnixMilli(unixMilli)
	return pgtype.Timestamp{
		Time:  t,
		Valid: true,
	}
}

type SyncRepository struct {
	db      *pgxpool.Pool
	queries *repo.Queries
}

func NewSyncRepository(db *pgxpool.Pool) *SyncRepository {
	return &SyncRepository{
		db:      db,
		queries: repo.New(db),
	}
}

// GetServerTimestamp returns current server timestamp
func (r *SyncRepository) GetServerTimestamp(ctx context.Context) (int64, error) {
	result, err := r.queries.GetServerTimestamp(ctx)
	if err != nil {
		return 0, err
	}
	return int64(result), nil
}

// SyncFriends returns friends updated after timestamp
func (r *SyncRepository) SyncFriends(ctx context.Context, userID int32, lastUpdatedAt int64, limit, offset int32) ([]repo.SyncFriendsRow, error) {
	return r.queries.SyncFriends(ctx, repo.SyncFriendsParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
		Limit:     limit,
		Offset:    offset,
	})
}

// SyncPendingFriendRequests returns pending friend requests
func (r *SyncRepository) SyncPendingFriendRequests(ctx context.Context, userID int32, lastUpdatedAt int64, limit int32) ([]repo.SyncPendingFriendRequestsRow, error) {
	return r.queries.SyncPendingFriendRequests(ctx, repo.SyncPendingFriendRequestsParams{
		FriendID:  userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
		Limit:     limit,
	})
}

// CountUpdatedFriends returns count of updated friends
func (r *SyncRepository) CountUpdatedFriends(ctx context.Context, userID int32, lastUpdatedAt int64) (int64, error) {
	return r.queries.CountUpdatedFriends(ctx, repo.CountUpdatedFriendsParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncMessages returns messages updated after timestamp
func (r *SyncRepository) SyncMessages(ctx context.Context, channelID int32, lastUpdatedAt int64, limit, offset int32) ([]repo.SyncMessagesRow, error) {
	return r.queries.SyncMessages(ctx, repo.SyncMessagesParams{
		ChannelID: channelID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
		Limit:     limit,
		Offset:    offset,
	})
}

// SyncUserMessages returns all messages for user's channels
func (r *SyncRepository) SyncUserMessages(ctx context.Context, userID int32, lastUpdatedAt int64, limit, offset int32) ([]repo.SyncUserMessagesRow, error) {
	return r.queries.SyncUserMessages(ctx, repo.SyncUserMessagesParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
		Limit:     limit,
		Offset:    offset,
	})
}

// CountUpdatedMessages returns count of updated messages
func (r *SyncRepository) CountUpdatedMessages(ctx context.Context, userID int32, lastUpdatedAt int64) (int64, error) {
	return r.queries.CountUpdatedMessages(ctx, repo.CountUpdatedMessagesParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncServers returns servers updated after timestamp
func (r *SyncRepository) SyncServers(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncServersRow, error) {
	return r.queries.SyncServers(ctx, repo.SyncServersParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// CountUpdatedServers returns count of updated servers
func (r *SyncRepository) CountUpdatedServers(ctx context.Context, userID int32, lastUpdatedAt int64) (int64, error) {
	return r.queries.CountUpdatedServers(ctx, repo.CountUpdatedServersParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncChannels returns channels updated after timestamp
func (r *SyncRepository) SyncChannels(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncChannelsRow, error) {
	return r.queries.SyncChannels(ctx, repo.SyncChannelsParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// CountUpdatedChannels returns count of updated channels
func (r *SyncRepository) CountUpdatedChannels(ctx context.Context, userID int32, lastUpdatedAt int64) (int64, error) {
	return r.queries.CountUpdatedChannels(ctx, repo.CountUpdatedChannelsParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncUserProfile returns user profile if updated
func (r *SyncRepository) SyncUserProfile(ctx context.Context, userID int32, lastUpdatedAt int64) (repo.SyncUserProfileRow, error) {
	return r.queries.SyncUserProfile(ctx, repo.SyncUserProfileParams{
		ID:        userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// GetUserLastUpdate returns user's last update timestamp
func (r *SyncRepository) GetUserLastUpdate(ctx context.Context, userID int32) (pgtype.Timestamp, error) {
	return r.queries.GetUserLastUpdate(ctx, userID)
}

// SyncVoiceChannels returns voice channels updated after timestamp
func (r *SyncRepository) SyncVoiceChannels(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncVoiceChannelsRow, error) {
	return r.queries.SyncVoiceChannels(ctx, repo.SyncVoiceChannelsParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncTextChannels returns text channels updated after timestamp
func (r *SyncRepository) SyncTextChannels(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncTextChannelsRow, error) {
	return r.queries.SyncTextChannels(ctx, repo.SyncTextChannelsParams{
		UserID:    userID,
		UpdatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncDirectMessages returns DM channels updated after timestamp
func (r *SyncRepository) SyncDirectMessages(ctx context.Context, userID int32, lastUpdatedAt int64, limit, offset int32) ([]repo.SyncDirectMessagesRow, error) {
	return r.queries.SyncDirectMessages(ctx, repo.SyncDirectMessagesParams{
		UserID:        userID,
		LastMessageAt: toPgTimestamp(lastUpdatedAt),
		Limit:         limit,
		Offset:        offset,
	})
}

// SyncPermissions returns permissions updated after timestamp
func (r *SyncRepository) SyncPermissions(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncPermissionsRow, error) {
	return r.queries.SyncPermissions(ctx, repo.SyncPermissionsParams{
		UserID:    userID,
		CreatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncVoiceStates returns voice states updated after timestamp
func (r *SyncRepository) SyncVoiceStates(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncVoiceStatesRow, error) {
	return r.queries.SyncVoiceStates(ctx, repo.SyncVoiceStatesParams{
		UserID:   userID,
		JoinedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncMessageAttachments returns attachments for message IDs
func (r *SyncRepository) SyncMessageAttachments(ctx context.Context, messageIDs []int32) ([]repo.SyncMessageAttachmentsRow, error) {
	return r.queries.SyncMessageAttachments(ctx, messageIDs)
}

// SyncMessageReactions returns reactions for message IDs
func (r *SyncRepository) SyncMessageReactions(ctx context.Context, messageIDs []int32) ([]repo.SyncMessageReactionsRow, error) {
	return r.queries.SyncMessageReactions(ctx, messageIDs)
}

// SyncServerMembers returns server members updated after timestamp
func (r *SyncRepository) SyncServerMembers(ctx context.Context, serverID int32, lastUpdatedAt int64) ([]repo.SyncServerMembersRow, error) {
	return r.queries.SyncServerMembers(ctx, repo.SyncServerMembersParams{
		ServerID: serverID,
		JoinedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncBans returns bans updated after timestamp
func (r *SyncRepository) SyncBans(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncBansRow, error) {
	return r.queries.SyncBans(ctx, repo.SyncBansParams{
		UserID:    userID,
		CreatedAt: toPgTimestamp(lastUpdatedAt),
	})
}

// SyncInvites returns invites updated after timestamp
func (r *SyncRepository) SyncInvites(ctx context.Context, userID int32, lastUpdatedAt int64) ([]repo.SyncInvitesRow, error) {
	return r.queries.SyncInvites(ctx, repo.SyncInvitesParams{
		UserID:    userID,
		CreatedAt: toPgTimestamp(lastUpdatedAt),
	})
}
