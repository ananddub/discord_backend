package controller

import (
	"context"

	"discord/gen/proto/service/dm"
	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	dmService "discord/internal/dm/service"

	"github.com/jackc/pgx/v5/pgtype"
)

type DMController struct {
	dm.UnimplementedDirectMessageServiceServer
	dmService *dmService.DMService
	queries   *repo.Queries
}

func NewDMController(dmService *dmService.DMService, queries *repo.Queries) *dm.DirectMessageServiceServer {
	controller := &DMController{
		dmService: dmService,
		queries:   queries,
	}
	var grpcController dm.DirectMessageServiceServer = controller
	return &grpcController
}

// CreateDMChannel creates a 1-on-1 DM channel
func (c *DMController) CreateDMChannel(ctx context.Context, req *dm.CreateDMChannelRequest) (*dm.CreateDMChannelResponse, error) {
	if req.GetUserId() == 0 || req.GetRecipientId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	dmChannelID, err := c.dmService.CreateDMChannel(ctx, req.GetUserId(), req.GetRecipientId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.CreateDMChannelResponse{
		DmChannelId: dmChannelID,
		Success:     true,
	}, nil
}

// CreateGroupDM creates a group DM channel
func (c *DMController) CreateGroupDM(ctx context.Context, req *dm.CreateGroupDMRequest) (*dm.CreateGroupDMResponse, error) {
	if req.GetOwnerId() == 0 || len(req.GetUserIds()) < 2 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	dmChannelID, err := c.dmService.CreateGroupDM(ctx, req.GetOwnerId(), req.GetUserIds(), req.GetName())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.CreateGroupDMResponse{
		DmChannelId: dmChannelID,
		Success:     true,
	}, nil
}

// GetDMChannel retrieves a DM channel
func (c *DMController) GetDMChannel(ctx context.Context, req *dm.GetDMChannelRequest) (*dm.GetDMChannelResponse, error) {
	channel, err := c.dmService.GetDMChannel(ctx, req.GetDmChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.GetDMChannelResponse{
		Id:             channel.ID,
		ParticipantIds: channel.ParticipantIDs,
		Name:           channel.Name,
		Icon:           channel.Icon,
	}, nil
}

// GetUserDMChannels retrieves all DM channels for a user
func (c *DMController) GetUserDMChannels(ctx context.Context, req *dm.GetUserDMChannelsRequest) (*dm.GetUserDMChannelsResponse, error) {
	channels, err := c.dmService.GetUserDMChannels(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	channelInfos := make([]*dm.DMChannelInfo, len(channels))
	for i, ch := range channels {
		channelInfos[i] = &dm.DMChannelInfo{
			Id:             ch.ID,
			ParticipantIds: ch.ParticipantIDs,
			Name:           ch.Name,
			Icon:           ch.Icon,
			LastMessageId:  ch.LastMessageID,
			LastMessageAt:  ch.LastMessageAt,
			UnreadCount:    ch.UnreadCount,
		}
	}

	return &dm.GetUserDMChannelsResponse{
		Channels: channelInfos,
	}, nil
}

// CloseDMChannel closes a DM channel for a user
func (c *DMController) CloseDMChannel(ctx context.Context, req *dm.CloseDMChannelRequest) (*dm.CloseDMChannelResponse, error) {
	err := c.dmService.CloseDMChannel(ctx, req.GetDmChannelId(), req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.CloseDMChannelResponse{
		Success: true,
	}, nil
}

// AddUserToGroupDM adds a user to a group DM
func (c *DMController) AddUserToGroupDM(ctx context.Context, req *dm.AddUserToGroupDMRequest) (*dm.AddUserToGroupDMResponse, error) {
	err := c.dmService.AddUserToGroupDM(ctx, req.GetDmChannelId(), req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.AddUserToGroupDMResponse{
		Success: true,
	}, nil
}

// RemoveUserFromGroupDM removes a user from a group DM
func (c *DMController) RemoveUserFromGroupDM(ctx context.Context, req *dm.RemoveUserFromGroupDMRequest) (*dm.RemoveUserFromGroupDMResponse, error) {
	err := c.dmService.RemoveUserFromGroupDM(ctx, req.GetDmChannelId(), req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.RemoveUserFromGroupDMResponse{
		Success: true,
	}, nil
}

// UpdateGroupDM updates group DM information
func (c *DMController) UpdateGroupDM(ctx context.Context, req *dm.UpdateGroupDMRequest) (*dm.UpdateGroupDMResponse, error) {
	var name, icon *string

	if req.GetName() != "" {
		n := req.GetName()
		name = &n
	}

	if req.GetIcon() != "" {
		i := req.GetIcon()
		icon = &i
	}

	err := c.dmService.UpdateGroupDM(ctx, req.GetDmChannelId(), name, icon)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.UpdateGroupDMResponse{
		Success: true,
	}, nil
}

// SendDM sends a direct message
func (c *DMController) SendDM(ctx context.Context, req *dm.SendDMRequest) (*dm.SendDMResponse, error) {
	// Validate input
	if req.GetDmChannelId() == 0 || req.GetSenderId() == 0 || req.GetContent() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Verify DM channel exists and user is participant
	channel, err := c.dmService.GetDMChannel(ctx, req.GetDmChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Check if sender is a participant
	isParticipant := false
	for _, participantID := range channel.ParticipantIDs {
		if participantID == req.GetSenderId() {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrPermissionDenied)
	}

	// Create message in database
	message, err := c.queries.CreateMessage(ctx, repo.CreateMessageParams{
		ChannelID:       req.GetDmChannelId(),
		SenderID:        req.GetSenderId(),
		Content:         req.GetContent(),
		MessageType:     pgtype.Text{String: "TEXT", Valid: true},
		MentionEveryone: pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Update last message info in DM channel
	err = c.dmService.UpdateLastMessage(ctx, req.GetDmChannelId(), message.ID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.SendDMResponse{
		MessageId: message.ID,
		CreatedAt: message.CreatedAt.Time.Unix(),
		Success:   true,
	}, nil
}

// GetDMMessages retrieves DM messages (streaming response)
func (c *DMController) GetDMMessages(req *dm.GetDMMessagesRequest, stream dm.DirectMessageService_GetDMMessagesServer) error {
	ctx := stream.Context()

	// Verify DM channel exists
	_, err := c.dmService.GetDMChannel(ctx, req.GetDmChannelId())
	if err != nil {
		return commonErrors.ToGRPCError(err)
	}

	// Get messages from database
	var messages []repo.Message
	if req.GetBeforeMessageId() != 0 {
		// Get messages before a specific message ID
		messages, err = c.queries.GetMessagesBefore(ctx, repo.GetMessagesBeforeParams{
			ChannelID: req.GetDmChannelId(),
			ID:        req.GetBeforeMessageId(),
			Limit:     req.GetLimit(),
		})
	} else {
		// Get latest messages
		limit := req.GetLimit()
		if limit == 0 {
			limit = 50 // Default limit
		}
		messages, err = c.queries.GetChannelMessages(ctx, repo.GetChannelMessagesParams{
			ChannelID: req.GetDmChannelId(),
			Limit:     limit,
			Offset:    0,
		})
	}

	if err != nil {
		return commonErrors.ToGRPCError(err)
	}

	// Stream messages to client
	for _, msg := range messages {
		if err := stream.Send(&dm.GetDMMessagesResponse{
			Id:        msg.ID,
			SenderId:  msg.SenderID,
			Content:   msg.Content,
			IsRead:    false, // TODO: Check against last_read_message_id
			IsEdited:  msg.IsEdited.Bool,
			CreatedAt: msg.CreatedAt.Time.Unix(),
		}); err != nil {
			return err
		}
	}

	return nil
}

// EditDM edits a direct message
func (c *DMController) EditDM(ctx context.Context, req *dm.EditDMRequest) (*dm.EditDMResponse, error) {
	// Validate input
	if req.GetMessageId() == 0 || req.GetContent() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Get message to verify it exists
	message, err := c.queries.GetMessageByID(ctx, req.GetMessageId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrNotFound)
	}

	// TODO: Verify user owns the message
	_ = message

	// Update message in database
	_, err = c.queries.UpdateMessage(ctx, repo.UpdateMessageParams{
		ID:      req.GetMessageId(),
		Content: req.GetContent(),
	})
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.EditDMResponse{
		Success: true,
	}, nil
}

// DeleteDM deletes a direct message
func (c *DMController) DeleteDM(ctx context.Context, req *dm.DeleteDMRequest) (*dm.DeleteDMResponse, error) {
	// Validate input
	if req.GetMessageId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Get message to verify it exists
	message, err := c.queries.GetMessageByID(ctx, req.GetMessageId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrNotFound)
	}

	// TODO: Verify user has permission to delete (message owner or admin)
	_ = message

	// Delete message from database
	err = c.queries.SoftDeleteMessage(ctx, req.GetMessageId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.DeleteDMResponse{
		Success: true,
	}, nil
}

// MarkAsRead marks messages as read
func (c *DMController) MarkAsRead(ctx context.Context, req *dm.MarkAsReadRequest) (*dm.MarkAsReadResponse, error) {
	err := c.dmService.MarkAsRead(ctx, req.GetDmChannelId(), req.GetUserId(), req.GetMessageId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &dm.MarkAsReadResponse{
		Success: true,
	}, nil
}
