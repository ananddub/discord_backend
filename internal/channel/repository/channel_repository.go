package repository

import (
	"context"

	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChannelRepository struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

func NewChannelRepository(db *pgxpool.Pool) *ChannelRepository {
	return &ChannelRepository{
		queries: repo.New(db),
		db:      db,
	}
}

// CreateChannel creates a new channel
func (r *ChannelRepository) CreateChannel(ctx context.Context, params repo.CreateChannelParams) (*repo.Channel, error) {
	channel, err := r.queries.CreateChannel(ctx, params)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// GetChannelByID retrieves channel by ID
func (r *ChannelRepository) GetChannelByID(ctx context.Context, id int32) (*repo.Channel, error) {
	channel, err := r.queries.GetChannelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// GetServerChannels retrieves all channels in a server
func (r *ChannelRepository) GetServerChannels(ctx context.Context, serverID int32) ([]repo.Channel, error) {
	return r.queries.GetServerChannels(ctx, serverID)
}

// GetChannelsByCategory retrieves all channels in a category
func (r *ChannelRepository) GetChannelsByCategory(ctx context.Context, categoryID int32) ([]repo.Channel, error) {
	return r.queries.GetChannelsByCategory(ctx, pgtype.Int4{Int32: categoryID, Valid: true})
}

// GetChannelsByType retrieves all channels of a specific type in a server
func (r *ChannelRepository) GetChannelsByType(ctx context.Context, serverID int32, channelType string) ([]repo.Channel, error) {
	return r.queries.GetChannelsByType(ctx, repo.GetChannelsByTypeParams{
		ServerID: serverID,
		Type:     channelType,
	})
}

// UpdateChannel updates channel information
func (r *ChannelRepository) UpdateChannel(ctx context.Context, params repo.UpdateChannelParams) (*repo.Channel, error) {
	channel, err := r.queries.UpdateChannel(ctx, params)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// DeleteChannel deletes a channel
func (r *ChannelRepository) DeleteChannel(ctx context.Context, id int32) error {
	return r.queries.DeleteChannel(ctx, id)
}

// UpdateChannelPosition updates channel position
func (r *ChannelRepository) UpdateChannelPosition(ctx context.Context, id int32, position int32) error {
	return r.queries.UpdateChannelPosition(ctx, repo.UpdateChannelPositionParams{
		ID:       id,
		Position: pgtype.Int4{Int32: position, Valid: true},
	})
}

// SetChannelPermission sets channel permissions for role or user
func (r *ChannelRepository) SetChannelPermission(ctx context.Context, params repo.SetChannelPermissionParams) (*repo.ChannelPermission, error) {
	permission, err := r.queries.SetChannelPermission(ctx, params)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetChannelPermissions retrieves all channel permissions
func (r *ChannelRepository) GetChannelPermissions(ctx context.Context, channelID int32) ([]repo.ChannelPermission, error) {
	return r.queries.GetChannelPermissions(ctx, channelID)
}

// GetRoleChannelPermissions retrieves channel permissions for a specific role
func (r *ChannelRepository) GetRoleChannelPermissions(ctx context.Context, channelID, roleID int32) (*repo.ChannelPermission, error) {
	permission, err := r.queries.GetRoleChannelPermissions(ctx, repo.GetRoleChannelPermissionsParams{
		ChannelID: channelID,
		RoleID:    pgtype.Int4{Int32: roleID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetUserChannelPermissions retrieves channel permissions for a specific user
func (r *ChannelRepository) GetUserChannelPermissions(ctx context.Context, channelID, userID int32) (*repo.ChannelPermission, error) {
	permission, err := r.queries.GetUserChannelPermissions(ctx, repo.GetUserChannelPermissionsParams{
		ChannelID: channelID,
		UserID:    pgtype.Int4{Int32: userID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// DeleteChannelPermission deletes a channel permission
func (r *ChannelRepository) DeleteChannelPermission(ctx context.Context, id int32) error {
	return r.queries.DeleteChannelPermission(ctx, id)
}

// DeleteRoleChannelPermissions deletes all channel permissions for a role
func (r *ChannelRepository) DeleteRoleChannelPermissions(ctx context.Context, channelID, roleID int32) error {
	return r.queries.DeleteRoleChannelPermissions(ctx, repo.DeleteRoleChannelPermissionsParams{
		ChannelID: channelID,
		RoleID:    pgtype.Int4{Int32: roleID, Valid: true},
	})
}

// DeleteUserChannelPermissions deletes all channel permissions for a user
func (r *ChannelRepository) DeleteUserChannelPermissions(ctx context.Context, channelID, userID int32) error {
	return r.queries.DeleteUserChannelPermissions(ctx, repo.DeleteUserChannelPermissionsParams{
		ChannelID: channelID,
		UserID:    pgtype.Int4{Int32: userID, Valid: true},
	})
}

// GetServerMembers gets members of a server
func (r *ChannelRepository) GetServerMembers(ctx context.Context, serverID int32, limit, offset int32) ([]repo.ServerMember, error) {
	return r.queries.GetServerMembers(ctx, repo.GetServerMembersParams{
		ServerID: serverID,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetChannelVoiceStates gets voice states for a channel
func (r *ChannelRepository) GetChannelVoiceStates(ctx context.Context, channelID int32) ([]repo.VoiceState, error) {
	return r.queries.GetChannelVoiceStates(ctx, channelID)
}
