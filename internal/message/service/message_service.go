package service

import (
	"context"
	"errors"

	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	messageRepo "discord/internal/message/repository"
	"discord/pkg/pubsub"
)

type MessageService struct {
	messageRepo *messageRepo.MessageRepository
	pubsub      *pubsub.PubSub
}

func NewMessageService(messageRepo *messageRepo.MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		pubsub:      pubsub.Get(),
	}
}

// SendMessage sends a new message
func (s *MessageService) SendMessage(ctx context.Context, channelID, senderID int32, content string, replyToMessageID *int32, mentionEveryone bool) (repo.Message, error) {
	if content == "" {
		return repo.Message{}, commonErrors.ErrInvalidInput
	}

	// Validate reply message if provided
	if replyToMessageID != nil {
		_, err := s.messageRepo.GetMessageByID(ctx, *replyToMessageID)
		if err != nil {
			return repo.Message{}, errors.New("reply message not found")
		}
	}

	return s.messageRepo.CreateMessage(ctx, channelID, senderID, content, "TEXT", replyToMessageID, mentionEveryone)
}

// GetMessage retrieves a single message
func (s *MessageService) GetMessage(ctx context.Context, messageID int32) (repo.Message, error) {
	return s.messageRepo.GetMessageByID(ctx, messageID)
}

// GetMessages retrieves messages with pagination
func (s *MessageService) GetMessages(ctx context.Context, channelID int32, limit, offset int32, beforeMessageID, afterMessageID *int32) ([]repo.Message, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	if beforeMessageID != nil {
		return s.messageRepo.GetMessagesBefore(ctx, channelID, *beforeMessageID, limit)
	}

	if afterMessageID != nil {
		return s.messageRepo.GetMessagesAfter(ctx, channelID, *afterMessageID, limit)
	}

	return s.messageRepo.GetChannelMessages(ctx, channelID, limit, offset)
}

// EditMessage edits an existing message
func (s *MessageService) EditMessage(ctx context.Context, messageID, userID int32, content string) (repo.Message, error) {
	if content == "" {
		return repo.Message{}, commonErrors.ErrInvalidInput
	}

	// Verify message exists and user owns it
	message, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return repo.Message{}, commonErrors.ErrNotFound
	}

	if message.SenderID != userID {
		return repo.Message{}, commonErrors.ErrPermissionDenied
	}

	return s.messageRepo.UpdateMessage(ctx, messageID, content)
}

// DeleteMessage deletes a message
func (s *MessageService) DeleteMessage(ctx context.Context, messageID, userID int32) error {
	// Verify message exists and user owns it
	message, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if message.SenderID != userID {
		return commonErrors.ErrPermissionDenied
	}

	// Delete attachments first
	_ = s.messageRepo.DeleteMessageAttachments(ctx, messageID)

	// Delete reactions
	_ = s.messageRepo.DeleteAllReactions(ctx, messageID)

	return s.messageRepo.DeleteMessage(ctx, messageID)
}

// PinMessage pins a message
func (s *MessageService) PinMessage(ctx context.Context, messageID int32) error {
	// Verify message exists
	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	return s.messageRepo.PinMessage(ctx, messageID)
}

// UnpinMessage unpins a message
func (s *MessageService) UnpinMessage(ctx context.Context, messageID int32) error {
	return s.messageRepo.UnpinMessage(ctx, messageID)
}

// GetPinnedMessages retrieves all pinned messages in a channel
func (s *MessageService) GetPinnedMessages(ctx context.Context, channelID int32) ([]repo.Message, error) {
	return s.messageRepo.GetPinnedMessages(ctx, channelID)
}

// AddReaction adds a reaction to a message
func (s *MessageService) AddReaction(ctx context.Context, messageID, userID int32, emoji string, emojiID *string) error {
	// Verify message exists
	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Check if user already reacted with this emoji
	reactions, err := s.messageRepo.GetReactionsByEmoji(ctx, messageID, emoji)
	if err == nil {
		for _, reaction := range reactions {
			if reaction.UserID == userID {
				return errors.New("already reacted with this emoji")
			}
		}
	}

	_, err = s.messageRepo.CreateReaction(ctx, messageID, userID, emoji, emojiID)
	return err
}

// RemoveReaction removes a reaction from a message
func (s *MessageService) RemoveReaction(ctx context.Context, messageID, userID int32, emoji string) error {
	return s.messageRepo.DeleteReaction(ctx, messageID, userID, emoji)
}

// GetReactions retrieves reactions for a message
func (s *MessageService) GetReactions(ctx context.Context, messageID int32, emoji *string) ([]repo.MessageReaction, error) {
	if emoji != nil {
		return s.messageRepo.GetReactionsByEmoji(ctx, messageID, *emoji)
	}
	return s.messageRepo.GetMessageReactions(ctx, messageID)
}

// BulkDeleteMessages deletes multiple messages
func (s *MessageService) BulkDeleteMessages(ctx context.Context, messageIDs []int32) error {
	if len(messageIDs) == 0 {
		return commonErrors.ErrInvalidInput
	}

	if len(messageIDs) > 100 {
		return errors.New("cannot delete more than 100 messages at once")
	}

	return s.messageRepo.BulkDeleteMessages(ctx, messageIDs)
}

// SearchMessages searches for messages in a channel
func (s *MessageService) SearchMessages(ctx context.Context, channelID int32, query string, limit, offset int32) ([]repo.Message, error) {
	if query == "" {
		return nil, commonErrors.ErrInvalidInput
	}

	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	return s.messageRepo.SearchMessages(ctx, channelID, query, limit, offset)
}

// CreateAttachment creates an attachment for a message
func (s *MessageService) CreateAttachment(ctx context.Context, messageID int32, fileURL, fileName, fileType string, fileSize int64, width, height *int32) (repo.MessageAttachment, error) {
	// Verify message exists
	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return repo.MessageAttachment{}, commonErrors.ErrNotFound
	}

	return s.messageRepo.CreateAttachment(ctx, messageID, fileURL, fileName, fileType, fileSize, width, height)
}

// GetMessageAttachments retrieves attachments for a message
func (s *MessageService) GetMessageAttachments(ctx context.Context, messageID int32) ([]repo.MessageAttachment, error) {
	return s.messageRepo.GetMessageAttachments(ctx, messageID)
}
