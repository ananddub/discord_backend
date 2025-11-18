package controller

import (
	"context"
	"fmt"
	"log"

	schema "discord/gen/proto/schema"
	"discord/gen/proto/service/dm"
	"discord/gen/repo"
	"discord/internal/dm/service"
	"discord/pkg/pubsub"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type DMController struct {
	dm.UnimplementedDirectMessageServiceServer
	service *service.MessageService
	db      *pgxpool.Pool
}

func NewDMController(db *pgxpool.Pool) *DMController {
	return &DMController{
		service: service.NewMessageService(),
		db:      db,
	}
}

// Helper function to convert repo.Message to schema.Message
func messageToProto(msg repo.Message) *schema.Message {
	return &schema.Message{
		Id:        msg.ID,
		ChannelId: int32(msg.ReceiverID.Int32), // For DMs, receiver_id is stored as channel_id
		SenderId:  msg.SenderID,
		Content:   msg.Content,
		IsEdited:  msg.IsEdited.Bool,
		CreatedAt: msg.CreatedAt.Time.Unix() * 1000,
	}
}

// ==================== MESSAGE OPERATIONS ====================

// SendMessage sends a direct message
func (c *DMController) SendMessage(ctx context.Context, req *dm.SendMessageRequest) (*dm.SendMessageResponse, error) {
	if req == nil {
		return &dm.SendMessageResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if req.ReceiverId == 0 {
		return &dm.SendMessageResponse{Success: false}, fmt.Errorf("receiver_id is required")
	}

	if req.Content == "" {
		return &dm.SendMessageResponse{Success: false}, fmt.Errorf("content cannot be empty")
	}

	message, err := c.service.SendMessage(ctx, req.ReceiverId, req.SenderId, req.Content, &req.ReplyToMessageId, false)
	if err != nil {
		return &dm.SendMessageResponse{Success: false}, fmt.Errorf("failed to send message: %w", err)
	}

	return &dm.SendMessageResponse{
		Success: true,
		Message: messageToProto(message),
	}, nil
}

// GetMessages retrieves messages between two users (streaming)
func (c *DMController) GetMessages(req *dm.GetMessagesRequest, stream dm.DirectMessageService_GetMessagesServer) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	ctx := stream.Context()

	if req.UserId == 0 {
		return fmt.Errorf("user_id is required")
	}

	// Validate limits
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Get messages from service
	messages, err := c.service.GetMessages(ctx, req.UserId, limit, offset, nil, nil)
	if err != nil {
		log.Printf("Failed to get messages: %v", err)
		return fmt.Errorf("failed to get messages: %w", err)
	}

	// Stream messages back
	for _, msg := range messages {
		resp := &dm.GetMessagesResponse{
			Messages: []*schema.Message{
				messageToProto(msg),
			},
		}
		if err := stream.Send(resp); err != nil {
			log.Printf("Failed to send message: %v", err)
			return err
		}
	}

	return nil
}

// GetMessage retrieves a single message by ID
func (c *DMController) GetMessage(ctx context.Context, req *dm.GetMessageRequest) (*dm.GetMessageResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return nil, fmt.Errorf("message_id is required")
	}

	message, err := c.service.GetMessage(ctx, req.MessageId)
	if err != nil {
		log.Printf("Failed to get message: %v", err)
		return nil, fmt.Errorf("message not found")
	}

	return &dm.GetMessageResponse{
		Id:         message.ID,
		ReceiverId: int32(message.ReceiverID.Int32),
		SenderId:   message.SenderID,
		Content:    message.Content,
		IsEdited:   message.IsEdited.Bool,
		CreatedAt:  message.CreatedAt.Time.Unix() * 1000,
	}, nil
}

// EditMessage edits a message
func (c *DMController) EditMessage(ctx context.Context, req *dm.EditMessageRequest) (*dm.EditMessageResponse, error) {
	if req == nil {
		return &dm.EditMessageResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return &dm.EditMessageResponse{Success: false}, fmt.Errorf("message_id is required")
	}

	if req.Content == "" {
		return &dm.EditMessageResponse{Success: false}, fmt.Errorf("content cannot be empty")
	}

	message, err := c.service.EditMessage(ctx, req.MessageId, 1, req.Content)
	if err != nil {
		log.Printf("Failed to update message: %v", err)
		return &dm.EditMessageResponse{Success: false}, fmt.Errorf("failed to update message: %w", err)
	}

	return &dm.EditMessageResponse{
		Success: true,
		Message: messageToProto(message),
	}, nil
}

// DeleteMessage deletes a message (soft delete)
func (c *DMController) DeleteMessage(ctx context.Context, req *dm.DeleteMessageRequest) (*dm.DeleteMessageResponse, error) {
	if req == nil {
		return &dm.DeleteMessageResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return &dm.DeleteMessageResponse{Success: false}, fmt.Errorf("message_id is required")
	}

	err := c.service.DeleteMessage(ctx, req.MessageId, 1) // 1 = sender ID from context
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
		return &dm.DeleteMessageResponse{Success: false}, fmt.Errorf("failed to delete message: %w", err)
	}

	return &dm.DeleteMessageResponse{Success: true}, nil
}

// ==================== MESSAGE INTERACTIONS ====================

// PinMessage pins a message
func (c *DMController) PinMessage(ctx context.Context, req *dm.PinMessageRequest) (*dm.PinMessageResponse, error) {
	if req == nil {
		return &dm.PinMessageResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return &dm.PinMessageResponse{Success: false}, fmt.Errorf("message_id is required")
	}

	err := c.service.PinMessage(ctx, req.MessageId)
	if err != nil {
		log.Printf("Failed to pin message: %v", err)
		return &dm.PinMessageResponse{Success: false}, fmt.Errorf("failed to pin message: %w", err)
	}

	return &dm.PinMessageResponse{Success: true}, nil
}

// UnpinMessage unpins a message
func (c *DMController) UnpinMessage(ctx context.Context, req *dm.UnpinMessageRequest) (*dm.UnpinMessageResponse, error) {
	if req == nil {
		return &dm.UnpinMessageResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return &dm.UnpinMessageResponse{Success: false}, fmt.Errorf("message_id is required")
	}

	err := c.service.UnpinMessage(ctx, req.MessageId)
	if err != nil {
		log.Printf("Failed to unpin message: %v", err)
		return &dm.UnpinMessageResponse{Success: false}, fmt.Errorf("failed to unpin message: %w", err)
	}

	return &dm.UnpinMessageResponse{Success: true}, nil
}

// GetPinnedMessages gets pinned messages
func (c *DMController) GetPinnedMessages(ctx context.Context, req *dm.GetPinnedMessagesRequest) (*dm.GetPinnedMessagesResponse, error) {
	if req == nil {
		return &dm.GetPinnedMessagesResponse{}, fmt.Errorf("request cannot be nil")
	}

	// TODO: Implement GetPinnedMessages in repo queries
	// For now, return empty response
	return &dm.GetPinnedMessagesResponse{
		MessageIds: []int32{},
	}, nil
}

// AddReaction adds a reaction to a message
func (c *DMController) AddReaction(ctx context.Context, req *dm.AddReactionRequest) (*dm.AddReactionResponse, error) {
	if req == nil {
		return &dm.AddReactionResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return &dm.AddReactionResponse{Success: false}, fmt.Errorf("message_id is required")
	}

	if req.Emoji == "" {
		return &dm.AddReactionResponse{Success: false}, fmt.Errorf("emoji is required")
	}

	err := c.service.AddReaction(ctx, req.MessageId, req.SenderId, req.Emoji, nil)
	if err != nil {
		log.Printf("Failed to add reaction: %v", err)
		return &dm.AddReactionResponse{Success: false}, fmt.Errorf("failed to add reaction: %w", err)
	}

	return &dm.AddReactionResponse{Success: true}, nil
}

// RemoveReaction removes a reaction from a message
func (c *DMController) RemoveReaction(ctx context.Context, req *dm.RemoveReactionRequest) (*dm.RemoveReactionResponse, error) {
	if req == nil {
		return &dm.RemoveReactionResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return &dm.RemoveReactionResponse{Success: false}, fmt.Errorf("message_id is required")
	}

	if req.Emoji == "" {
		return &dm.RemoveReactionResponse{Success: false}, fmt.Errorf("emoji is required")
	}

	err := c.service.RemoveReaction(ctx, req.MessageId, req.SenderId, req.Emoji)
	if err != nil {
		log.Printf("Failed to remove reaction: %v", err)
		return &dm.RemoveReactionResponse{Success: false}, fmt.Errorf("failed to remove reaction: %w", err)
	}

	return &dm.RemoveReactionResponse{Success: true}, nil
}

// GetReactions gets reactions on a message
func (c *DMController) GetReactions(ctx context.Context, req *dm.GetReactionsRequest) (*dm.GetReactionsResponse, error) {
	if req == nil {
		return &dm.GetReactionsResponse{}, fmt.Errorf("request cannot be nil")
	}

	if req.MessageId == 0 {
		return &dm.GetReactionsResponse{}, fmt.Errorf("message_id is required")
	}

	reactions, err := c.service.GetReactions(ctx, req.MessageId, nil)
	if err != nil {
		log.Printf("Failed to get reactions: %v", err)
		return &dm.GetReactionsResponse{}, fmt.Errorf("failed to get reactions: %w", err)
	}

	// Group reactions by emoji
	reactionMap := make(map[string]*dm.ReactionInfo)
	for _, reaction := range reactions {
		if _, exists := reactionMap[reaction.Emoji]; !exists {
			reactionMap[reaction.Emoji] = &dm.ReactionInfo{
				Emoji:   reaction.Emoji,
				Count:   0,
				UserIds: []int32{},
			}
		}
		reactionMap[reaction.Emoji].Count++
		reactionMap[reaction.Emoji].UserIds = append(reactionMap[reaction.Emoji].UserIds, reaction.UserID)
	}

	var reactionList []*dm.ReactionInfo
	for _, reaction := range reactionMap {
		reactionList = append(reactionList, reaction)
	}

	return &dm.GetReactionsResponse{
		Reactions: reactionList,
	}, nil
}

// ==================== BULK OPERATIONS ====================

// BulkDeleteMessages deletes multiple messages
func (c *DMController) BulkDeleteMessages(ctx context.Context, req *dm.BulkDeleteMessagesRequest) (*dm.BulkDeleteMessagesResponse, error) {
	if req == nil {
		return &dm.BulkDeleteMessagesResponse{Success: false}, fmt.Errorf("request cannot be nil")
	}

	if len(req.MessageIds) == 0 {
		return &dm.BulkDeleteMessagesResponse{Success: false}, fmt.Errorf("message_ids cannot be empty")
	}

	err := c.service.BulkDeleteMessages(ctx, req.MessageIds)
	if err != nil {
		log.Printf("Failed to bulk delete messages: %v", err)
		return &dm.BulkDeleteMessagesResponse{Success: false}, fmt.Errorf("failed to bulk delete messages: %w", err)
	}

	return &dm.BulkDeleteMessagesResponse{
		DeletedCount: int32(len(req.MessageIds)),
		Success:      true,
	}, nil
}

// SearchMessages searches for messages
func (c *DMController) SearchMessages(ctx context.Context, req *dm.SearchMessagesRequest) (*dm.SearchMessagesResponse, error) {
	if req == nil {
		return &dm.SearchMessagesResponse{}, fmt.Errorf("request cannot be nil")
	}

	if req.Query == "" {
		return &dm.SearchMessagesResponse{}, fmt.Errorf("query cannot be empty")
	}

	// Validate limits
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	messages, err := c.service.SearchMessages(ctx, req.ReceiverId, req.Query, limit, offset)
	if err != nil {
		log.Printf("Failed to search messages: %v", err)
		return &dm.SearchMessagesResponse{}, fmt.Errorf("failed to search messages: %w", err)
	}

	var results []*dm.MessageSearchResult
	for _, msg := range messages {
		result := &dm.MessageSearchResult{
			MessageId:  msg.ID,
			ReceiverId: int32(msg.ReceiverID.Int32),
			SenderId:   msg.SenderID,
			Content:    msg.Content,
			CreatedAt:  msg.CreatedAt.Time.Unix() * 1000,
		}
		results = append(results, result)
	}

	return &dm.SearchMessagesResponse{
		Results:      results,
		TotalResults: int32(len(results)),
	}, nil
}
func (c *DMController) SendTyping(stream grpc.BidiStreamingServer[dm.SendTypingRequest, dm.SendTypingResponse]) error {
	pub := pubsub.Get()
	ctx := stream.Context()
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		return fmt.Errorf("user_id not found in context")
	}
	go func() {
		for {
			value, err := stream.Recv()
			if err != nil {
				log.Println("Receive error:", err)
				return
			}

			pub.Publish(fmt.Sprintf("msgtyping:%d", value.ReceiverId), value)
			log.Println("Received typing event:", value)
		}
	}()
	ch := pub.Subscribe(fmt.Sprintf("msgtyping:%d", userID))
	defer ch.Close()
	for msg := range ch.Receive() {
		typingEvent, ok := msg.(*dm.SendTypingRequest)
		if !ok {
			log.Println("Invalid message type")
			continue
		}
		stream.Send(&dm.SendTypingResponse{
			SenderId:   typingEvent.SenderId,
			IsTyping:   typingEvent.IsTyping,
			ReceiverId: typingEvent.ReceiverId,
		})
	}
	return nil
}
