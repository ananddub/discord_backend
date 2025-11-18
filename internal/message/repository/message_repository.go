package repository

import (
	"context"
	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepository struct {
	db      *pgxpool.Pool
	queries *repo.Queries
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		db:      db,
		queries: repo.New(db),
	}
}

func (r *MessageRepository) CreateMessage(ctx context.Context, channelID, senderID int32, content, messageType string, replyToMessageID *int32, mentionEveryone bool) (repo.Message, error) {
	var replyTo pgtype.Int4
	if replyToMessageID != nil {
		replyTo = pgtype.Int4{Int32: *replyToMessageID, Valid: true}
	}

	return r.queries.CreateMessage(ctx, repo.CreateMessageParams{
		ChannelID:        pgtype.Int4{Int32: channelID, Valid: true},
		SenderID:         senderID,
		Content:          content,
		MessageType:      pgtype.Text{String: messageType, Valid: true},
		ReplyToMessageID: replyTo,
		MentionEveryone:  pgtype.Bool{Bool: mentionEveryone, Valid: true},
	})
}

func (r *MessageRepository) GetMessageByID(ctx context.Context, messageID int32) (repo.Message, error) {
	return r.queries.GetMessageByID(ctx, messageID)
}

func (r *MessageRepository) GetChannelMessages(ctx context.Context, channelID int32, limit, offset int32) ([]repo.Message, error) {
	return r.queries.GetChannelMessages(ctx, repo.GetChannelMessagesParams{
		ChannelID: pgtype.Int4{
			Int32: channelID,
			Valid: true,
		},
		Limit:  limit,
		Offset: offset,
	})
}

func (r *MessageRepository) GetMessagesBefore(ctx context.Context, channelID, beforeMessageID, limit int32) ([]repo.Message, error) {
	return r.queries.GetMessagesBefore(ctx, repo.GetMessagesBeforeParams{
		ChannelID: pgtype.Int4{
			Int32: channelID,
			Valid: true,
		},
		ID:    beforeMessageID,
		Limit: limit,
	})
}

func (r *MessageRepository) GetMessagesAfter(ctx context.Context, channelID, afterMessageID, limit int32) ([]repo.Message, error) {
	return r.queries.GetMessagesAfter(ctx, repo.GetMessagesAfterParams{
		ChannelID: pgtype.Int4{
			Int32: channelID,
			Valid: true,
		},
		ID:    afterMessageID,
		Limit: limit,
	})
}

func (r *MessageRepository) UpdateMessage(ctx context.Context, messageID int32, content string) (repo.Message, error) {
	return r.queries.UpdateMessage(ctx, repo.UpdateMessageParams{
		ID:      messageID,
		Content: content,
	})
}

func (r *MessageRepository) DeleteMessage(ctx context.Context, messageID int32) error {
	_, err := r.queries.SoftDeleteMessage(ctx, messageID)
	return err
}

func (r *MessageRepository) PinMessage(ctx context.Context, messageID int32) error {
	_, err := r.queries.PinMessage(ctx, messageID)
	return err
}

func (r *MessageRepository) UnpinMessage(ctx context.Context, messageID int32) error {
	_, err := r.queries.UnpinMessage(ctx, messageID)
	return err
}

func (r *MessageRepository) GetPinnedMessages(ctx context.Context, channelID int32) ([]repo.Message, error) {
	return r.queries.GetPinnedMessages(ctx, pgtype.Int4{
		Int32: channelID,
		Valid: true,
	})
}

func (r *MessageRepository) BulkDeleteMessages(ctx context.Context, messageIDs []int32) error {
	_, err := r.queries.BulkSoftDeleteMessages(ctx, messageIDs)
	return err
}

func (r *MessageRepository) SearchMessages(ctx context.Context, channelID int32, query string, limit, offset int32) ([]repo.Message, error) {
	return r.queries.SearchMessages(ctx, repo.SearchMessagesParams{
		ChannelID: pgtype.Int4{
			Int32: channelID,
			Valid: true,
		},
		Column2: pgtype.Text{String: query, Valid: true},
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *MessageRepository) GetUserMessages(ctx context.Context, userID int32, limit, offset int32) ([]repo.Message, error) {
	return r.queries.GetUserMessages(ctx, repo.GetUserMessagesParams{
		SenderID: userID,
		Limit:    limit,
		Offset:   offset,
	})
}

func (r *MessageRepository) CreateReaction(ctx context.Context, messageID, userID int32, emoji string, emojiID *string) (repo.MessageReaction, error) {
	var emojiIDType pgtype.Text
	if emojiID != nil {
		emojiIDType = pgtype.Text{String: *emojiID, Valid: true}
	}

	return r.queries.CreateReaction(ctx, repo.CreateReactionParams{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		EmojiID:   emojiIDType,
	})
}

func (r *MessageRepository) GetMessageReactions(ctx context.Context, messageID int32) ([]repo.MessageReaction, error) {
	return r.queries.GetMessageReactions(ctx, messageID)
}

func (r *MessageRepository) GetReactionsByEmoji(ctx context.Context, messageID int32, emoji string) ([]repo.MessageReaction, error) {
	return r.queries.GetReactionsByEmoji(ctx, repo.GetReactionsByEmojiParams{
		MessageID: messageID,
		Emoji:     emoji,
	})
}

func (r *MessageRepository) DeleteReaction(ctx context.Context, messageID, userID int32, emoji string) error {
	_, err := r.queries.DeleteReaction(ctx, repo.DeleteReactionParams{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	})
	return err
}

func (r *MessageRepository) DeleteAllReactions(ctx context.Context, messageID int32) error {
	_, err := r.queries.DeleteAllReactions(ctx, messageID)
	return err
}

func (r *MessageRepository) CreateAttachment(ctx context.Context, messageID int32, fileURL, fileName, fileType string, fileSize int64, width, height *int32) (repo.MessageAttachment, error) {
	var widthType, heightType pgtype.Int4
	if width != nil {
		widthType = pgtype.Int4{Int32: *width, Valid: true}
	}
	if height != nil {
		heightType = pgtype.Int4{Int32: *height, Valid: true}
	}

	return r.queries.CreateMessageAttachment(ctx, repo.CreateMessageAttachmentParams{
		MessageID: messageID,
		FileUrl:   fileURL,
		FileName:  fileName,
		FileType:  fileType,
		FileSize:  fileSize,
		Width:     widthType,
		Height:    heightType,
	})
}

func (r *MessageRepository) GetMessageAttachments(ctx context.Context, messageID int32) ([]repo.MessageAttachment, error) {
	return r.queries.GetMessageAttachments(ctx, messageID)
}

func (r *MessageRepository) DeleteAttachment(ctx context.Context, attachmentID int32) error {
	_, err := r.queries.SoftDeleteMessageAttachment(ctx, attachmentID)
	return err
}

func (r *MessageRepository) DeleteMessageAttachments(ctx context.Context, messageID int32) error {
	_, err := r.queries.SoftDeleteMessageAttachments(ctx, messageID)
	return err
}
