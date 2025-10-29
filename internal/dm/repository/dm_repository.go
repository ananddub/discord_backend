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

// CreateDMChannel creates a new DM channel
func (r *DMRepository) CreateDMChannel(ctx context.Context, name, icon string, ownerID *int32, isGroup bool) (*repo.DmChannel, error) {
	params := repo.CreateDMChannelParams{
		IsGroup: pgtype.Bool{Bool: isGroup, Valid: true},
	}

	if name != "" {
		params.Name = pgtype.Text{String: name, Valid: true}
	}

	if icon != "" {
		params.Icon = pgtype.Text{String: icon, Valid: true}
	}

	if ownerID != nil {
		params.OwnerID = pgtype.Int4{Int32: *ownerID, Valid: true}
	}

	channel, err := r.queries.CreateDMChannel(ctx, params)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// GetDMChannelByID retrieves DM channel by ID
func (r *DMRepository) GetDMChannelByID(ctx context.Context, id int32) (*repo.DmChannel, error) {
	channel, err := r.queries.GetDMChannelByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// GetUserDMChannels retrieves all DM channels for a user
func (r *DMRepository) GetUserDMChannels(ctx context.Context, userID int32) ([]repo.DmChannel, error) {
	return r.queries.GetUserDMChannels(ctx, userID)
}

// GetDMChannelForUsers retrieves existing DM channel between two users
func (r *DMRepository) GetDMChannelForUsers(ctx context.Context, userID1, userID2 int32) (*repo.DmChannel, error) {
	channel, err := r.queries.GetDMChannelForUsers(ctx, repo.GetDMChannelForUsersParams{
		UserID:   userID1,
		UserID_2: userID2,
	})
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// UpdateDMChannel updates DM channel information
func (r *DMRepository) UpdateDMChannel(ctx context.Context, id int32, name, icon *string, lastMessageID *int32, lastMessageAt *pgtype.Timestamp) (*repo.DmChannel, error) {
	params := repo.UpdateDMChannelParams{
		ID: id,
	}

	if name != nil {
		params.Name = pgtype.Text{String: *name, Valid: true}
	}

	if icon != nil {
		params.Icon = pgtype.Text{String: *icon, Valid: true}
	}

	if lastMessageID != nil {
		params.LastMessageID = pgtype.Int4{Int32: *lastMessageID, Valid: true}
	}

	if lastMessageAt != nil {
		params.LastMessageAt = *lastMessageAt
	}

	channel, err := r.queries.UpdateDMChannel(ctx, params)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// DeleteDMChannel deletes a DM channel
func (r *DMRepository) DeleteDMChannel(ctx context.Context, id int32) error {
	return r.queries.DeleteDMChannel(ctx, id)
}

// AddDMParticipant adds a participant to DM channel
func (r *DMRepository) AddDMParticipant(ctx context.Context, dmChannelID, userID int32) (*repo.DmParticipant, error) {
	participant, err := r.queries.AddDMParticipant(ctx, repo.AddDMParticipantParams{
		DmChannelID: dmChannelID,
		UserID:      userID,
	})
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// GetDMParticipants retrieves all participants in a DM channel
func (r *DMRepository) GetDMParticipants(ctx context.Context, dmChannelID int32) ([]repo.DmParticipant, error) {
	return r.queries.GetDMParticipants(ctx, dmChannelID)
}

// RemoveDMParticipant removes a participant from DM channel
func (r *DMRepository) RemoveDMParticipant(ctx context.Context, dmChannelID, userID int32) error {
	return r.queries.RemoveDMParticipant(ctx, repo.RemoveDMParticipantParams{
		DmChannelID: dmChannelID,
		UserID:      userID,
	})
}

// UpdateLastReadMessage updates the last read message for a participant
func (r *DMRepository) UpdateLastReadMessage(ctx context.Context, dmChannelID, userID, messageID int32) error {
	return r.queries.UpdateLastReadMessage(ctx, repo.UpdateLastReadMessageParams{
		DmChannelID:       dmChannelID,
		UserID:            userID,
		LastReadMessageID: pgtype.Int4{Int32: messageID, Valid: true},
	})
}
