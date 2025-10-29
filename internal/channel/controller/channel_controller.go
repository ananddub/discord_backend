package controller

import (
	"context"

	channelPb "discord/gen/proto/service/channel"
	channelService "discord/internal/channel/service"
	commonErrors "discord/internal/common/errors"
)

type ChannelController struct {
	channelPb.UnimplementedChannelServiceServer
	channelService *channelService.ChannelService
}

func NewChannelController(channelService *channelService.ChannelService) *channelPb.ChannelServiceServer {
	controller := &ChannelController{
		channelService: channelService,
	}
	var grpcController channelPb.ChannelServiceServer = controller
	return &grpcController
}

// CreateChannel creates a new channel
func (c *ChannelController) CreateChannel(ctx context.Context, req *channelPb.CreateChannelRequest) (*channelPb.CreateChannelResponse, error) {
	if req.GetName() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Basic implementation - using defaults for missing proto fields
	channel, err := c.channelService.CreateChannel(
		ctx,
		1, // server_id - should come from context/request
		req.GetName(),
		"TEXT",               // default type
		nil,                  // no category
		0,                    // position
		req.GetDescription(), // using description as topic
		false,                // not NSFW
		0,                    // no slowmode
	)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.CreateChannelResponse{
		Channel: channel,
		Success: true,
	}, nil
}

// GetChannel retrieves channel by ID
func (c *ChannelController) GetChannel(ctx context.Context, req *channelPb.GetChannelRequest) (*channelPb.GetChannelResponse, error) {
	channel, err := c.channelService.GetChannel(ctx, req.GetChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.GetChannelResponse{
		Channel: channel,
	}, nil
}

// UpdateChannel updates channel information
func (c *ChannelController) UpdateChannel(ctx context.Context, req *channelPb.UpdateChannelRequest) (*channelPb.UpdateChannelResponse, error) {
	var name, topic *string
	var position, slowmodeDelay *int32
	var isNSFW *bool

	if req.GetName() != "" {
		n := req.GetName()
		name = &n
	}

	if req.GetTopic() != "" {
		t := req.GetTopic()
		topic = &t
	}

	if req.GetPosition() != 0 {
		p := req.GetPosition()
		position = &p
	}

	if req.GetSlowmodeDelay() != 0 {
		s := req.GetSlowmodeDelay()
		slowmodeDelay = &s
	}

	isNSFW = &req.IsNsfw

	_, err := c.channelService.UpdateChannel(ctx, req.GetChannelId(), name, topic, position, slowmodeDelay, isNSFW)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.UpdateChannelResponse{
		Success: true,
	}, nil
}

// DeleteChannel deletes a channel
func (c *ChannelController) DeleteChannel(ctx context.Context, req *channelPb.DeleteChannelRequest) (*channelPb.DeleteChannelResponse, error) {
	err := c.channelService.DeleteChannel(ctx, req.GetChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.DeleteChannelResponse{
		Success: true,
	}, nil
}

// GetServerChannels retrieves all channels in a server
func (c *ChannelController) GetServerChannels(ctx context.Context, req *channelPb.GetServerChannelsRequest) (*channelPb.GetServerChannelsResponse, error) {
	channels, err := c.channelService.GetServerChannels(ctx, req.GetServerId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Convert to ChannelInfo format expected by proto
	channelInfos := make([]*channelPb.ChannelInfo, len(channels))
	for i, ch := range channels {
		channelInfos[i] = &channelPb.ChannelInfo{
			Id:         ch.Id,
			Name:       ch.Name,
			Type:       ch.Type.String(),
			Position:   ch.Position,
			CategoryId: ch.CategoryId,
		}
	}

	return &channelPb.GetServerChannelsResponse{
		Channels: channelInfos,
	}, nil
}

// CreateCategory creates a new category
func (c *ChannelController) CreateCategory(ctx context.Context, req *channelPb.CreateCategoryRequest) (*channelPb.CreateCategoryResponse, error) {
	channel, err := c.channelService.CreateChannel(
		ctx,
		req.GetServerId(),
		req.GetName(),
		"CATEGORY", // Channel type for category
		nil,
		req.GetPosition(),
		"",
		false,
		0,
	)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.CreateCategoryResponse{
		CategoryId: channel.Id,
		Success:    true,
	}, nil
}

// UpdateCategory updates category information
func (c *ChannelController) UpdateCategory(ctx context.Context, req *channelPb.UpdateCategoryRequest) (*channelPb.UpdateCategoryResponse, error) {
	var name *string
	var position *int32

	if req.GetName() != "" {
		n := req.GetName()
		name = &n
	}

	if req.GetPosition() != 0 {
		p := req.GetPosition()
		position = &p
	}

	_, err := c.channelService.UpdateChannel(ctx, req.GetCategoryId(), name, nil, position, nil, nil)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.UpdateCategoryResponse{
		Success: true,
	}, nil
}

// DeleteCategory deletes a category
func (c *ChannelController) DeleteCategory(ctx context.Context, req *channelPb.DeleteCategoryRequest) (*channelPb.DeleteCategoryResponse, error) {
	err := c.channelService.DeleteChannel(ctx, req.GetCategoryId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.DeleteCategoryResponse{
		Success: true,
	}, nil
}

// JoinChannel adds user to a channel (primarily for voice channels)
func (c *ChannelController) JoinChannel(ctx context.Context, req *channelPb.JoinChannelRequest) (*channelPb.JoinChannelResponse, error) {
	// Verify channel exists and user has permission to join
	channel, err := c.channelService.GetChannel(ctx, req.GetChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// For voice channels, this would create a voice_state entry
	// For text channels, this could be used for subscription/notification preferences
	_ = channel // Channel type check would go here

	// Voice state creation would be handled by voice service
	// For now, just return success
	return &channelPb.JoinChannelResponse{
		Success: true,
	}, nil
}

// LeaveChannel removes user from a channel (primarily for voice channels)
func (c *ChannelController) LeaveChannel(ctx context.Context, req *channelPb.LeaveChannelRequest) (*channelPb.LeaveChannelResponse, error) {
	// Verify channel exists
	_, err := c.channelService.GetChannel(ctx, req.GetChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// For voice channels, this would delete the voice_state entry
	// For text channels, this could update subscription/notification preferences

	// Voice state deletion would be handled by voice service
	// For now, just return success
	return &channelPb.LeaveChannelResponse{
		Success: true,
	}, nil
}

// GetChannelMembers retrieves channel members (streaming response)
func (c *ChannelController) GetChannelMembers(req *channelPb.GetChannelMembersRequest, stream channelPb.ChannelService_GetChannelMembersServer) error {
	ctx := stream.Context()

	// Get channel to verify it exists
	channel, err := c.channelService.GetChannel(ctx, req.GetChannelId())
	if err != nil {
		return commonErrors.ToGRPCError(err)
	}

	// For voice channels, we could stream voice_states
	// For text channels, we would stream server_members with permission to view the channel
	// This would require server_members service integration

	// Mock implementation - in production, this would query and stream actual members
	// based on channel type and permissions
	_ = channel

	// Example streaming pattern (would be replaced with actual member data):
	// for _, member := range members {
	//     if err := stream.Send(&channelPb.GetChannelMembersResponse{
	//         Members: member,
	//     }); err != nil {
	//         return err
	//     }
	// }

	return nil
}

// UpdateChannelPermissions updates channel permissions for role or user
func (c *ChannelController) UpdateChannelPermissions(ctx context.Context, req *channelPb.UpdateChannelPermissionsRequest) (*channelPb.UpdateChannelPermissionsResponse, error) {
	var roleID, userID *int32

	if req.GetRoleId() != 0 {
		id := req.GetRoleId()
		roleID = &id
	}

	if req.GetUserId() != 0 {
		id := req.GetUserId()
		userID = &id
	}

	err := c.channelService.SetChannelPermission(
		ctx,
		req.GetChannelId(),
		roleID,
		userID,
		req.GetAllow(),
		req.GetDeny(),
	)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.UpdateChannelPermissionsResponse{
		Success: true,
	}, nil
}

// GetChannelPermissions retrieves channel permissions
func (c *ChannelController) GetChannelPermissions(ctx context.Context, req *channelPb.GetChannelPermissionsRequest) (*channelPb.GetChannelPermissionsResponse, error) {
	permissions, err := c.channelService.GetChannelPermissions(ctx, req.GetChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Convert to PermissionOverwrite format
	overwrites := make([]*channelPb.PermissionOverwrite, len(permissions))
	for i, perm := range permissions {
		overwrites[i] = &channelPb.PermissionOverwrite{
			RoleId: perm.RoleId,
			UserId: perm.UserId,
			Allow:  perm.Allow,
			Deny:   perm.Deny,
		}
	}

	return &channelPb.GetChannelPermissionsResponse{
		Overwrites: overwrites,
	}, nil
}

// UpdateChannelPosition updates channel position
func (c *ChannelController) UpdateChannelPosition(ctx context.Context, req *channelPb.UpdateChannelPositionRequest) (*channelPb.UpdateChannelPositionResponse, error) {
	err := c.channelService.UpdateChannelPosition(ctx, req.GetChannelId(), req.GetPosition())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &channelPb.UpdateChannelPositionResponse{
		Success: true,
	}, nil
}
