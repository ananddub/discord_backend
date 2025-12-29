package service

import (
	"context"
	"errors"

	"discord/config"
	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	"discord/internal/dm/repository"
	"discord/pkg/pubsub"
)

type MessageService struct {
	messageRepo *repository.DMRepository
	pubsub      *pubsub.PubSub
}

func NewMessageService() *MessageService {
	db := config.GetWriteDb()
	return &MessageService{
		messageRepo: repository.NewDMRepository(db),
		pubsub:      pubsub.Get(),
	}
}

func (s *MessageService) SendMessage(ctx context.Context, reciver_id, senderID int32, content string, replyToMessageID *int32, mentionEveryone bool) (repo.Message, error) {
	if content == "" {
		return repo.Message{}, commonErrors.ErrInvalidInput
	}

	if replyToMessageID != nil {
		_, err := s.messageRepo.GetChatMessageByID(ctx, *replyToMessageID)
		if err != nil {
			return repo.Message{}, errors.New("reply message not found")
		}
	}

	return s.messageRepo.CreateDMMessage(ctx, reciver_id, senderID, content, "default", replyToMessageID, mentionEveryone)
}

func (s *MessageService) GetMessage(ctx context.Context, messageID int32) (repo.Message, error) {
	return s.messageRepo.GetChatMessageByID(ctx, messageID)
}

func (s *MessageService) GetMessages(ctx context.Context, receiver_id int32, limit, offset int32, beforeMessageID, afterMessageID *int32) ([]repo.Message, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	if beforeMessageID != nil {
		return s.messageRepo.GetMessagesBefore(ctx, receiver_id, *beforeMessageID, limit)
	}

	if afterMessageID != nil {
		return s.messageRepo.GetMessagesAfter(ctx, receiver_id, *afterMessageID, limit)
	}

	return s.messageRepo.GetChatMessages(ctx, receiver_id, limit, offset)
}

func (s *MessageService) EditMessage(ctx context.Context, messageID, userID int32, content string) (repo.Message, error) {
	if content == "" {
		return repo.Message{}, commonErrors.ErrInvalidInput
	}

	message, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return repo.Message{}, commonErrors.ErrNotFound
	}

	if message.SenderID != userID {
		return repo.Message{}, commonErrors.ErrPermissionDenied
	}

	return s.messageRepo.UpdateMessage(ctx, messageID, content)
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID, userID int32) error {

	message, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if message.SenderID != userID {
		return commonErrors.ErrPermissionDenied
	}

	_ = s.messageRepo.DeleteMessageAttachments(ctx, messageID)

	_ = s.messageRepo.DeleteAllReactions(ctx, messageID)

	return s.messageRepo.DeleteMessage(ctx, messageID)
}

func (s *MessageService) PinMessage(ctx context.Context, messageID int32) error {

	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	return s.messageRepo.PinMessage(ctx, messageID)
}

func (s *MessageService) UnpinMessage(ctx context.Context, messageID int32) error {
	return s.messageRepo.UnpinMessage(ctx, messageID)
}

func (s *MessageService) GetPinnedMessages(ctx context.Context, channelID int32) ([]repo.Message, error) {
	return s.messageRepo.GetPinnedMessages(ctx, channelID)
}

func (s *MessageService) AddReaction(ctx context.Context, messageID, userID int32, emoji string, emojiID *string) error {

	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

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

func (s *MessageService) RemoveReaction(ctx context.Context, messageID, userID int32, emoji string) error {
	return s.messageRepo.DeleteReaction(ctx, messageID, userID, emoji)
}

func (s *MessageService) GetReactions(ctx context.Context, messageID int32, emoji *string) ([]repo.MessageReaction, error) {
	if emoji != nil {
		return s.messageRepo.GetReactionsByEmoji(ctx, messageID, *emoji)
	}
	return s.messageRepo.GetMessageReactions(ctx, messageID)
}

func (s *MessageService) BulkDeleteMessages(ctx context.Context, messageIDs []int32) error {
	if len(messageIDs) == 0 {
		return commonErrors.ErrInvalidInput
	}

	if len(messageIDs) > 100 {
		return errors.New("cannot delete more than 100 messages at once")
	}

	return s.messageRepo.BulkDeleteMessages(ctx, messageIDs)
}

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

func (s *MessageService) CreateAttachment(ctx context.Context, messageID int32, fileURL, fileName, fileType string, fileSize int64, width, height *int32) (repo.MessageAttachment, error) {

	_, err := s.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return repo.MessageAttachment{}, commonErrors.ErrNotFound
	}

	return s.messageRepo.CreateAttachment(ctx, messageID, fileURL, fileName, fileType, fileSize, width, height)
}

func (s *MessageService) GetMessageAttachments(ctx context.Context, messageID int32) ([]repo.MessageAttachment, error) {
	return s.messageRepo.GetMessageAttachments(ctx, messageID)
}
