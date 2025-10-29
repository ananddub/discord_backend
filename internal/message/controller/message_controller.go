package controller

import (
	"context"

	"discord/gen/proto/schema"
	messagePb "discord/gen/proto/service/message"
	commonErrors "discord/internal/common/errors"
	messageService "discord/internal/message/service"
)

type MessageController struct {
	messagePb.UnimplementedMessageServiceServer
	messageService *messageService.MessageService
}

func NewMessageController(messageService *messageService.MessageService) *messagePb.MessageServiceServer {
	controller := &MessageController{
		messageService: messageService,
	}
	var grpcController messagePb.MessageServiceServer = controller
	return &grpcController
}

// SendMessage sends a new message
func (c *MessageController) SendMessage(ctx context.Context, req *messagePb.SendMessageRequest) (*messagePb.SendMessageResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetChannelId() == 0 || req.GetContent() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	var replyToID *int32
	if req.GetReplyToMessageId() != 0 {
		id := req.GetReplyToMessageId()
		replyToID = &id
	}

	message, err := c.messageService.SendMessage(ctx, req.GetChannelId(), userID, req.GetContent(), replyToID, false)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.SendMessageResponse{
		Message: &schema.Message{
			Id:        message.ID,
			ChannelId: message.ChannelID,
			SenderId:  message.SenderID,
			Content:   message.Content,
			IsEdited:  message.IsEdited.Bool,
			IsPinned:  message.IsPinned.Bool,
			CreatedAt: message.CreatedAt.Time.Unix(),
		},
		Success: true,
	}, nil
}

// GetMessages retrieves messages (streaming response)
func (c *MessageController) GetMessages(req *messagePb.GetMessagesRequest, stream messagePb.MessageService_GetMessagesServer) error {
	ctx := stream.Context()

	if req.GetChannelId() == 0 {
		return commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	limit := req.GetLimit()
	if limit == 0 {
		limit = 50
	}

	messages, err := c.messageService.GetMessages(ctx, req.GetChannelId(), limit, req.GetOffset(), nil, nil)
	if err != nil {
		return commonErrors.ToGRPCError(err)
	}

	// Stream messages to client
	for _, msg := range messages {
		pbMsg := &schema.Message{
			Id:        msg.ID,
			ChannelId: msg.ChannelID,
			SenderId:  msg.SenderID,
			Content:   msg.Content,
			IsEdited:  msg.IsEdited.Bool,
			IsPinned:  msg.IsPinned.Bool,
			CreatedAt: msg.CreatedAt.Time.Unix(),
		}

		if err := stream.Send(&messagePb.GetMessagesResponse{
			Messages: []*schema.Message{pbMsg},
		}); err != nil {
			return err
		}
	}

	return nil
}

// GetMessage retrieves a single message
func (c *MessageController) GetMessage(ctx context.Context, req *messagePb.GetMessageRequest) (*messagePb.GetMessageResponse, error) {
	if req.GetMessageId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	message, err := c.messageService.GetMessage(ctx, req.GetMessageId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.GetMessageResponse{
		Id:        message.ID,
		ChannelId: message.ChannelID,
		SenderId:  message.SenderID,
		Content:   message.Content,
		IsEdited:  message.IsEdited.Bool,
		CreatedAt: message.CreatedAt.Time.Unix(),
	}, nil
}

// EditMessage edits an existing message
func (c *MessageController) EditMessage(ctx context.Context, req *messagePb.EditMessageRequest) (*messagePb.EditMessageResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetMessageId() == 0 || req.GetContent() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	message, err := c.messageService.EditMessage(ctx, req.GetMessageId(), userID, req.GetContent())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.EditMessageResponse{
		Message: &schema.Message{
			Id:        message.ID,
			ChannelId: message.ChannelID,
			SenderId:  message.SenderID,
			Content:   message.Content,
			IsEdited:  message.IsEdited.Bool,
			CreatedAt: message.CreatedAt.Time.Unix(),
		},
		Success: true,
	}, nil
}

// DeleteMessage deletes a message
func (c *MessageController) DeleteMessage(ctx context.Context, req *messagePb.DeleteMessageRequest) (*messagePb.DeleteMessageResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetMessageId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.messageService.DeleteMessage(ctx, req.GetMessageId(), userID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.DeleteMessageResponse{
		Success: true,
	}, nil
}

// PinMessage pins a message
func (c *MessageController) PinMessage(ctx context.Context, req *messagePb.PinMessageRequest) (*messagePb.PinMessageResponse, error) {
	if req.GetMessageId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// TODO: Check user permissions for pinning messages

	err := c.messageService.PinMessage(ctx, req.GetMessageId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.PinMessageResponse{
		Success: true,
	}, nil
}

// UnpinMessage unpins a message
func (c *MessageController) UnpinMessage(ctx context.Context, req *messagePb.UnpinMessageRequest) (*messagePb.UnpinMessageResponse, error) {
	if req.GetMessageId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.messageService.UnpinMessage(ctx, req.GetMessageId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.UnpinMessageResponse{
		Success: true,
	}, nil
}

// GetPinnedMessages retrieves all pinned messages
func (c *MessageController) GetPinnedMessages(ctx context.Context, req *messagePb.GetPinnedMessagesRequest) (*messagePb.GetPinnedMessagesResponse, error) {
	if req.GetChannelId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	messages, err := c.messageService.GetPinnedMessages(ctx, req.GetChannelId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	messageIDs := make([]int32, len(messages))
	for i, msg := range messages {
		messageIDs[i] = msg.ID
	}

	return &messagePb.GetPinnedMessagesResponse{
		MessageIds: messageIDs,
	}, nil
}

// AddReaction adds a reaction to a message
func (c *MessageController) AddReaction(ctx context.Context, req *messagePb.AddReactionRequest) (*messagePb.AddReactionResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetMessageId() == 0 || req.GetEmoji() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.messageService.AddReaction(ctx, req.GetMessageId(), userID, req.GetEmoji(), nil)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.AddReactionResponse{
		Success: true,
	}, nil
}

// RemoveReaction removes a reaction from a message
func (c *MessageController) RemoveReaction(ctx context.Context, req *messagePb.RemoveReactionRequest) (*messagePb.RemoveReactionResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetMessageId() == 0 || req.GetEmoji() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.messageService.RemoveReaction(ctx, req.GetMessageId(), userID, req.GetEmoji())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.RemoveReactionResponse{
		Success: true,
	}, nil
}

// GetReactions retrieves reactions for a message
func (c *MessageController) GetReactions(ctx context.Context, req *messagePb.GetReactionsRequest) (*messagePb.GetReactionsResponse, error) {
	if req.GetMessageId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	var emoji *string
	if req.GetEmoji() != "" {
		e := req.GetEmoji()
		emoji = &e
	}

	reactions, err := c.messageService.GetReactions(ctx, req.GetMessageId(), emoji)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Group reactions by emoji
	reactionMap := make(map[string]*messagePb.ReactionInfo)
	for _, reaction := range reactions {
		if _, exists := reactionMap[reaction.Emoji]; !exists {
			reactionMap[reaction.Emoji] = &messagePb.ReactionInfo{
				Emoji:   reaction.Emoji,
				Count:   0,
				UserIds: []int32{},
			}
		}
		reactionMap[reaction.Emoji].Count++
		reactionMap[reaction.Emoji].UserIds = append(reactionMap[reaction.Emoji].UserIds, reaction.UserID)
	}

	reactionInfos := make([]*messagePb.ReactionInfo, 0, len(reactionMap))
	for _, info := range reactionMap {
		reactionInfos = append(reactionInfos, info)
	}

	return &messagePb.GetReactionsResponse{
		Reactions: reactionInfos,
	}, nil
}

// SendTyping sends a typing indicator
func (c *MessageController) SendTyping(ctx context.Context, req *messagePb.SendTypingRequest) (*messagePb.SendTypingResponse, error) {
	if req.GetChannelId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// TODO: Implement typing indicator logic (websocket broadcast)
	// For now, just return success

	return &messagePb.SendTypingResponse{
		Success: true,
	}, nil
}

// BulkDeleteMessages deletes multiple messages
func (c *MessageController) BulkDeleteMessages(ctx context.Context, req *messagePb.BulkDeleteMessagesRequest) (*messagePb.BulkDeleteMessagesResponse, error) {
	if len(req.GetMessageIds()) == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// TODO: Check user permissions for bulk delete

	err := c.messageService.BulkDeleteMessages(ctx, req.GetMessageIds())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &messagePb.BulkDeleteMessagesResponse{
		DeletedCount: int32(len(req.GetMessageIds())),
		Success:      true,
	}, nil
}

// SearchMessages searches for messages
func (c *MessageController) SearchMessages(ctx context.Context, req *messagePb.SearchMessagesRequest) (*messagePb.SearchMessagesResponse, error) {
	if req.GetChannelId() == 0 || req.GetQuery() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	limit := req.GetLimit()
	if limit == 0 {
		limit = 25
	}

	messages, err := c.messageService.SearchMessages(ctx, req.GetChannelId(), req.GetQuery(), limit, req.GetOffset())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	results := make([]*messagePb.MessageSearchResult, len(messages))
	for i, msg := range messages {
		results[i] = &messagePb.MessageSearchResult{
			MessageId: msg.ID,
			ChannelId: msg.ChannelID,
			SenderId:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Time.Unix(),
		}
	}

	return &messagePb.SearchMessagesResponse{
		Results:      results,
		TotalResults: int32(len(results)),
	}, nil
}
