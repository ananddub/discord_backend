package controller

import (
	"context"

	"discord/gen/proto/schema"
	voicePb "discord/gen/proto/service/voice_channel"
	commonErrors "discord/internal/common/errors"
	voiceService "discord/internal/voice/service"
)

type VoiceController struct {
	voicePb.UnimplementedVoiceChannelServiceServer
	voiceService *voiceService.VoiceService
}

func NewVoiceController(voiceService *voiceService.VoiceService) *voicePb.VoiceChannelServiceServer {
	controller := &VoiceController{
		voiceService: voiceService,
	}
	var grpcController voicePb.VoiceChannelServiceServer = controller
	return &grpcController
}

// CreateVoiceChannel creates a new voice channel (note: basic implementation)
func (c *VoiceController) CreateVoiceChannel(ctx context.Context, req *voicePb.CreateVoiceChannelRequest) (*voicePb.CreateVoiceChannelResponse, error) {
	if req.GetChannelId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// TODO: Implement actual voice channel creation in channels table
	// For now, return success with basic response
	return &voicePb.CreateVoiceChannelResponse{
		VoiceChannel: &schema.VoiceChannel{
			Id: req.GetChannelId(),
		},
		Success: true,
	}, nil
}

// JoinVoiceChannel allows a user to join a voice channel
func (c *VoiceController) JoinVoiceChannel(ctx context.Context, req *voicePb.JoinVoiceChannelRequest) (*voicePb.JoinVoiceChannelResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetVoiceChannelId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Validate permissions
	err := c.voiceService.ValidateVoicePermissions(ctx, userID, req.GetVoiceChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Join voice channel
	_, err = c.voiceService.JoinVoiceChannel(ctx, userID, req.GetVoiceChannelId(), nil)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &voicePb.JoinVoiceChannelResponse{
		Success: true,
	}, nil
}

// LeaveVoiceChannel allows a user to leave a voice channel
func (c *VoiceController) LeaveVoiceChannel(ctx context.Context, req *voicePb.LeaveVoiceChannelRequest) (*voicePb.LeaveVoiceChannelResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetVoiceChannelId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Leave voice channel
	err := c.voiceService.LeaveVoiceChannel(ctx, userID, req.GetVoiceChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &voicePb.LeaveVoiceChannelResponse{
		Success: true,
	}, nil
}

// SendVoiceChat sends a voice chat message (stub implementation)
func (c *VoiceController) SendVoiceChat(ctx context.Context, req *voicePb.SendVoiceChatRequest) (*voicePb.SendVoiceChatResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetVoiceChannelId() == 0 || req.GetChat() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Verify user is in the voice channel
	isInChannel := c.voiceService.IsUserInVoiceChannel(ctx, userID, req.GetVoiceChannelId())
	if !isInChannel {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrPermissionDenied)
	}

	// TODO: Implement actual voice chat message handling
	return &voicePb.SendVoiceChatResponse{
		VoiceChat: &schema.VoiceChat{
			Id:             1, // Placeholder
			VoiceChannelId: req.GetVoiceChannelId(),
			UserId:         userID,
			Chat:           req.GetChat(),
		},
		Success: true,
	}, nil
}
