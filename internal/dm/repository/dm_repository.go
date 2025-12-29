package repository

import (
	"context"
	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DMRepository struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

func NewDMRepository(db *pgxpool.Pool) *DMRepository {
	return &DMRepository{
		queries: repo.New(db),
		db:      db,
	}
}

func (r *DMRepository) CreateDMMessage(ctx context.Context, receiverID, senderID int32, content, messageType string, replyToMessageID *int32, mentionEveryone bool) (repo.Message, error) {
	var replyTo pgtype.Int4
	if replyToMessageID != nil {
		replyTo = pgtype.Int4{Int32: *replyToMessageID, Valid: true}
	}

	return r.queries.CreateChatMessage(ctx, repo.CreateChatMessageParams{
		ReceiverID:       pgtype.Int4{Int32: receiverID, Valid: true},
		SenderID:         senderID,
		Content:          content,
		MessageType:      pgtype.Text{String: messageType, Valid: true},
		ReplyToMessageID: replyTo,
		MentionEveryone:  pgtype.Bool{Bool: false, Valid: true},
	})
}

func (r *DMRepository) GetDMMessageByID(ctx context.Context, messageID int32) (repo.Message, error) {
	return r.queries.GetChatMessageByID(ctx, messageID)
}

func (r *DMRepository) GetDMMessages(ctx context.Context, userID1, userID2, limit, offset int32) ([]repo.Message, error) {
	return r.queries.GetChatMessages(ctx, repo.GetChatMessagesParams{
		SenderID:   userID1,
		ReceiverID: pgtype.Int4{Int32: userID2, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
}

func (r *DMRepository) GetDMMessagesBefore(ctx context.Context, userID1, userID2, beforeMessageID, limit int32) ([]repo.Message, error) {
	return r.queries.GetChatMessagesBefore(ctx, repo.GetChatMessagesBeforeParams{
		SenderID:   userID1,
		ReceiverID: pgtype.Int4{Int32: userID2, Valid: true},
		ID:         beforeMessageID,
		Limit:      limit,
	})
}

func (r *DMRepository) GetDMMessagesAfter(ctx context.Context, userID1, userID2, afterMessageID, limit int32) ([]repo.Message, error) {
	return r.queries.GetChatMessagesAfter(ctx, repo.GetChatMessagesAfterParams{
		SenderID:   userID1,
		ReceiverID: pgtype.Int4{Int32: userID2, Valid: true},
		ID:         afterMessageID,
		Limit:      limit,
	})
}

func (r *DMRepository) UpdateMessage(ctx context.Context, messageID int32, content string) (repo.Message, error) {
	return r.queries.UpdateChatMessage(ctx, repo.UpdateChatMessageParams{
		ID:      messageID,
		Content: content,
	})
}

func (r *DMRepository) DeleteMessage(ctx context.Context, messageID int32) error {
	_, err := r.queries.SoftDeleteChatMessage(ctx, messageID)
	return err
}

func (r *DMRepository) SearchDMMessages(ctx context.Context, userID1, userID2 int32, query string, limit, offset int32) ([]repo.Message, error) {
	return r.queries.SearchChatMessages(ctx, repo.SearchChatMessagesParams{
		SenderID:   userID1,
		ReceiverID: pgtype.Int4{Int32: userID2, Valid: true},
		Column3:    pgtype.Text{String: query, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
}

func (r *DMRepository) PinMessage(ctx context.Context, messageID int32) error {
	_, err := r.queries.PinMessage(ctx, messageID)
	return err
}

func (r *DMRepository) UnpinMessage(ctx context.Context, messageID int32) error {
	_, err := r.queries.UnpinMessage(ctx, messageID)
	return err
}

func (r *DMRepository) CreateReaction(ctx context.Context, messageID, userID int32, emoji string, emojiID *string) (repo.MessageReaction, error) {
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

func (r *DMRepository) GetMessageReactions(ctx context.Context, messageID int32) ([]repo.MessageReaction, error) {
	return r.queries.GetMessageReactions(ctx, messageID)
}

func (r *DMRepository) GetReactionsByEmoji(ctx context.Context, messageID int32, emoji string) ([]repo.MessageReaction, error) {
	return r.queries.GetReactionsByEmoji(ctx, repo.GetReactionsByEmojiParams{
		MessageID: messageID,
		Emoji:     emoji,
	})
}

func (r *DMRepository) DeleteReaction(ctx context.Context, messageID, userID int32, emoji string) error {
	_, err := r.queries.DeleteReaction(ctx, repo.DeleteReactionParams{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	})
	return err
}

func (r *DMRepository) DeleteAllReactions(ctx context.Context, messageID int32) error {
	_, err := r.queries.DeleteAllReactions(ctx, messageID)
	return err
}

func (r *DMRepository) CreateAttachment(ctx context.Context, messageID int32, fileURL, fileName, fileType string, fileSize int64, width, height *int32) (repo.MessageAttachment, error) {
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

func (r *DMRepository) GetMessageAttachments(ctx context.Context, messageID int32) ([]repo.MessageAttachment, error) {
	return r.queries.GetMessageAttachments(ctx, messageID)
}

func (r *DMRepository) DeleteAttachment(ctx context.Context, attachmentID int32) error {
	_, err := r.queries.SoftDeleteMessageAttachment(ctx, attachmentID)
	return err
}

func (r *DMRepository) DeleteMessageAttachments(ctx context.Context, messageID int32) error {
	_, err := r.queries.SoftDeleteMessageAttachments(ctx, messageID)
	return err
}

func (r *DMRepository) GetChatMessageByID(ctx context.Context, messageID int32) (repo.Message, error) {
	return r.queries.GetChatMessageByID(ctx, messageID)
}

func (r *DMRepository) GetMessageByID(ctx context.Context, messageID int32) (repo.Message, error) {
	return r.queries.GetChatMessageByID(ctx, messageID)
}

func (r *DMRepository) GetChatMessages(ctx context.Context, receiverID, limit, offset int32) ([]repo.Message, error) {
	return r.queries.GetChatMessages(ctx, repo.GetChatMessagesParams{
		ReceiverID: pgtype.Int4{Int32: receiverID, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
}

func (r *DMRepository) GetMessagesBefore(ctx context.Context, receiverID, beforeMessageID, limit int32) ([]repo.Message, error) {
	return r.queries.GetChatMessagesBefore(ctx, repo.GetChatMessagesBeforeParams{
		ReceiverID: pgtype.Int4{Int32: receiverID, Valid: true},
		ID:         beforeMessageID,
		Limit:      limit,
	})
}

func (r *DMRepository) GetMessagesAfter(ctx context.Context, receiverID, afterMessageID, limit int32) ([]repo.Message, error) {
	return r.queries.GetChatMessagesAfter(ctx, repo.GetChatMessagesAfterParams{
		ReceiverID: pgtype.Int4{Int32: receiverID, Valid: true},
		ID:         afterMessageID,
		Limit:      limit,
	})
}

func (r *DMRepository) GetPinnedMessages(ctx context.Context, channelID int32) ([]repo.Message, error) {
	return r.queries.GetPinnedMessages(ctx, pgtype.Int4{Int32: channelID, Valid: true})
}

func (r *DMRepository) BulkDeleteMessages(ctx context.Context, messageIDs []int32) error {
	_, err := r.queries.BulkSoftDeleteMessages(ctx, messageIDs)
	return err
}

func (r *DMRepository) SearchMessages(ctx context.Context, channelID int32, query string, limit, offset int32) ([]repo.Message, error) {
	return r.queries.SearchChatMessages(ctx, repo.SearchChatMessagesParams{
		ReceiverID: pgtype.Int4{Int32: channelID, Valid: true},
		Column3:    pgtype.Text{String: query, Valid: true},
		Limit:      limit,
		Offset:     offset,
	})
}
