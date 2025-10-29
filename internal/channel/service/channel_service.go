package service

import (
	"context"
	"errors"

	"discord/gen/proto/schema"
	"discord/gen/repo"
	channelRepo "discord/internal/channel/repository"
	commonErrors "discord/internal/common/errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type ChannelService struct {
	channelRepo *channelRepo.ChannelRepository
}

func NewChannelService(channelRepo *channelRepo.ChannelRepository) *ChannelService {
	return &ChannelService{
		channelRepo: channelRepo,
	}
}

// CreateChannel creates a new channel
func (s *ChannelService) CreateChannel(ctx context.Context, serverID int32, name, channelType string, categoryID *int32, position int32, topic string, isNSFW bool, slowmodeDelay int32) (*schema.Channel, error) {
	// Validate input
	if name == "" {
		return nil, commonErrors.ErrInvalidInput
	}

	// Create channel params
	params := repo.CreateChannelParams{
		ServerID: serverID,
		Name:     name,
		Type:     channelType,
		Position: pgtype.Int4{Int32: position, Valid: true},
	}

	if categoryID != nil {
		params.CategoryID = pgtype.Int4{Int32: *categoryID, Valid: true}
	}

	if topic != "" {
		params.Topic = pgtype.Text{String: topic, Valid: true}
	}

	params.IsNsfw = pgtype.Bool{Bool: isNSFW, Valid: true}
	params.SlowmodeDelay = pgtype.Int4{Int32: slowmodeDelay, Valid: true}

	// Create channel
	channel, err := s.channelRepo.CreateChannel(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.toProtoChannel(channel), nil
}

// GetChannel retrieves channel by ID
func (s *ChannelService) GetChannel(ctx context.Context, channelID int32) (*schema.Channel, error) {
	channel, err := s.channelRepo.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, commonErrors.ErrNotFound
	}

	return s.toProtoChannel(channel), nil
}

// GetServerChannels retrieves all channels in a server
func (s *ChannelService) GetServerChannels(ctx context.Context, serverID int32) ([]*schema.Channel, error) {
	channels, err := s.channelRepo.GetServerChannels(ctx, serverID)
	if err != nil {
		return nil, err
	}

	protoChannels := make([]*schema.Channel, len(channels))
	for i, channel := range channels {
		protoChannels[i] = s.toProtoChannel(&channel)
	}

	return protoChannels, nil
}

// GetChannelsByCategory retrieves all channels in a category
func (s *ChannelService) GetChannelsByCategory(ctx context.Context, categoryID int32) ([]*schema.Channel, error) {
	channels, err := s.channelRepo.GetChannelsByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	protoChannels := make([]*schema.Channel, len(channels))
	for i, channel := range channels {
		protoChannels[i] = s.toProtoChannel(&channel)
	}

	return protoChannels, nil
}

// GetChannelsByType retrieves channels by type
func (s *ChannelService) GetChannelsByType(ctx context.Context, serverID int32, channelType string) ([]*schema.Channel, error) {
	channels, err := s.channelRepo.GetChannelsByType(ctx, serverID, channelType)
	if err != nil {
		return nil, err
	}

	protoChannels := make([]*schema.Channel, len(channels))
	for i, channel := range channels {
		protoChannels[i] = s.toProtoChannel(&channel)
	}

	return protoChannels, nil
}

// UpdateChannel updates channel information
func (s *ChannelService) UpdateChannel(ctx context.Context, channelID int32, name, topic *string, position, slowmodeDelay *int32, isNSFW *bool) (*schema.Channel, error) {
	params := repo.UpdateChannelParams{
		ID: channelID,
	}

	if name != nil {
		params.Name = pgtype.Text{String: *name, Valid: true}
	}

	if topic != nil {
		params.Topic = pgtype.Text{String: *topic, Valid: true}
	}

	if position != nil {
		params.Position = pgtype.Int4{Int32: *position, Valid: true}
	}

	if isNSFW != nil {
		params.IsNsfw = pgtype.Bool{Bool: *isNSFW, Valid: true}
	}

	if slowmodeDelay != nil {
		params.SlowmodeDelay = pgtype.Int4{Int32: *slowmodeDelay, Valid: true}
	}

	channel, err := s.channelRepo.UpdateChannel(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.toProtoChannel(channel), nil
}

// DeleteChannel deletes a channel
func (s *ChannelService) DeleteChannel(ctx context.Context, channelID int32) error {
	return s.channelRepo.DeleteChannel(ctx, channelID)
}

// UpdateChannelPosition updates channel position
func (s *ChannelService) UpdateChannelPosition(ctx context.Context, channelID int32, position int32) error {
	return s.channelRepo.UpdateChannelPosition(ctx, channelID, position)
}

// SetChannelPermission sets channel permissions for role or user
func (s *ChannelService) SetChannelPermission(ctx context.Context, channelID int32, roleID, userID *int32, allowPermissions, denyPermissions int64) error {
	if roleID == nil && userID == nil {
		return errors.New("either role_id or user_id must be provided")
	}

	params := repo.SetChannelPermissionParams{
		ChannelID:        channelID,
		AllowPermissions: pgtype.Int8{Int64: allowPermissions, Valid: true},
		DenyPermissions:  pgtype.Int8{Int64: denyPermissions, Valid: true},
	}

	if roleID != nil {
		params.RoleID = pgtype.Int4{Int32: *roleID, Valid: true}
	}

	if userID != nil {
		params.UserID = pgtype.Int4{Int32: *userID, Valid: true}
	}

	_, err := s.channelRepo.SetChannelPermission(ctx, params)
	return err
}

// GetChannelPermissions retrieves all channel permissions
func (s *ChannelService) GetChannelPermissions(ctx context.Context, channelID int32) ([]*schema.ChannelPermissionOverwrite, error) {
	permissions, err := s.channelRepo.GetChannelPermissions(ctx, channelID)
	if err != nil {
		return nil, err
	}

	protoPermissions := make([]*schema.ChannelPermissionOverwrite, len(permissions))
	for i, perm := range permissions {
		protoPermissions[i] = s.toProtoChannelPermission(&perm)
	}

	return protoPermissions, nil
}

// DeleteChannelPermission deletes a channel permission
func (s *ChannelService) DeleteChannelPermission(ctx context.Context, permissionID int32) error {
	return s.channelRepo.DeleteChannelPermission(ctx, permissionID)
}

// GetChannelMembers gets members who have access to a channel
func (s *ChannelService) GetChannelMembers(ctx context.Context, serverID, channelID int32) ([]ChannelMemberInfo, error) {
	// Get server members
	members, err := s.channelRepo.GetServerMembers(ctx, serverID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// Get voice states for this channel (if it's a voice channel)
	voiceStates, _ := s.channelRepo.GetChannelVoiceStates(ctx, channelID)
	voiceStateMap := make(map[int32]*repo.VoiceState)
	for _, vs := range voiceStates {
		voiceStateMap[vs.UserID] = &vs
	}

	// Build member info list
	memberInfos := make([]ChannelMemberInfo, 0, len(members))
	for _, member := range members {
		info := ChannelMemberInfo{
			UserId:   member.UserID,
			JoinedAt: member.JoinedAt.Time.UnixMilli(),
		}

		if member.Nickname.Valid {
			info.Nickname = member.Nickname.String
		}

		// Add voice state info if user is in voice channel
		if vs, exists := voiceStateMap[member.UserID]; exists {
			if vs.IsMuted.Valid {
				info.IsMuted = vs.IsMuted.Bool
			}
			if vs.IsDeafened.Valid {
				info.IsDeafened = vs.IsDeafened.Bool
			}
		}

		memberInfos = append(memberInfos, info)
	}

	return memberInfos, nil
}

// ChannelMemberInfo represents member information for a channel
type ChannelMemberInfo struct {
	UserId     int32
	Username   string
	Nickname   string
	ProfilePic string
	Status     string
	JoinedAt   int64
	IsMuted    bool
	IsDeafened bool
}

// Helper function to convert repo.Channel to proto Channel
func (s *ChannelService) toProtoChannel(channel *repo.Channel) *schema.Channel {
	protoChannel := &schema.Channel{
		Id:       channel.ID,
		ServerId: channel.ServerID,
		Name:     channel.Name,
		Type:     schema.ChannelType(schema.ChannelType_value[channel.Type]),
	}

	if channel.Position.Valid {
		protoChannel.Position = channel.Position.Int32
	}

	if channel.CategoryID.Valid {
		protoChannel.CategoryId = channel.CategoryID.Int32
	}

	if channel.Topic.Valid {
		protoChannel.Topic = channel.Topic.String
	}

	if channel.IsNsfw.Valid {
		protoChannel.IsNsfw = channel.IsNsfw.Bool
	}

	if channel.SlowmodeDelay.Valid {
		protoChannel.SlowmodeDelay = channel.SlowmodeDelay.Int32
	}

	protoChannel.CreatedAt = channel.CreatedAt.Time.Unix()
	protoChannel.UpdatedAt = channel.UpdatedAt.Time.Unix()

	return protoChannel
}

// Helper function to convert repo.ChannelPermission to proto ChannelPermissionOverwrite
func (s *ChannelService) toProtoChannelPermission(perm *repo.ChannelPermission) *schema.ChannelPermissionOverwrite {
	protoPerm := &schema.ChannelPermissionOverwrite{
		Id:        perm.ID,
		ChannelId: perm.ChannelID,
	}

	if perm.AllowPermissions.Valid {
		protoPerm.Allow = perm.AllowPermissions.Int64
	}

	if perm.DenyPermissions.Valid {
		protoPerm.Deny = perm.DenyPermissions.Int64
	}

	if perm.RoleID.Valid {
		protoPerm.RoleId = perm.RoleID.Int32
	}

	if perm.UserID.Valid {
		protoPerm.UserId = perm.UserID.Int32
	}

	return protoPerm
}
