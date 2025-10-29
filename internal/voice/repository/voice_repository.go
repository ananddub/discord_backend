package repository

import (
	"context"

	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoiceRepository struct {
	db *pgxpool.Pool
	q  *repo.Queries
}

func NewVoiceRepository(db *pgxpool.Pool) *VoiceRepository {
	return &VoiceRepository{
		db: db,
		q:  repo.New(db),
	}
}

// CreateVoiceState creates a new voice state for a user in a channel
func (r *VoiceRepository) CreateVoiceState(ctx context.Context, userID, channelID int32, serverID *int32, sessionID string, isMuted, isDeafened, selfMute, selfDeaf bool) (repo.VoiceState, error) {
	params := repo.CreateVoiceStateParams{
		UserID:     userID,
		ChannelID:  channelID,
		SessionID:  sessionID,
		IsMuted:    pgtype.Bool{Bool: isMuted, Valid: true},
		IsDeafened: pgtype.Bool{Bool: isDeafened, Valid: true},
		SelfMute:   pgtype.Bool{Bool: selfMute, Valid: true},
		SelfDeaf:   pgtype.Bool{Bool: selfDeaf, Valid: true},
	}

	if serverID != nil {
		params.ServerID = pgtype.Int4{Int32: *serverID, Valid: true}
	}

	return r.q.CreateVoiceState(ctx, params)
}

// GetVoiceState retrieves a voice state for a user in a channel
func (r *VoiceRepository) GetVoiceState(ctx context.Context, userID, channelID int32) (repo.VoiceState, error) {
	return r.q.GetVoiceState(ctx, repo.GetVoiceStateParams{
		UserID:    userID,
		ChannelID: channelID,
	})
}

// GetUserVoiceState retrieves the current voice state for a user
func (r *VoiceRepository) GetUserVoiceState(ctx context.Context, userID int32) (repo.VoiceState, error) {
	return r.q.GetUserVoiceState(ctx, userID)
}

// GetChannelVoiceStates retrieves all voice states in a channel
func (r *VoiceRepository) GetChannelVoiceStates(ctx context.Context, channelID int32) ([]repo.VoiceState, error) {
	return r.q.GetChannelVoiceStates(ctx, channelID)
}

// UpdateVoiceState updates voice state properties
func (r *VoiceRepository) UpdateVoiceState(ctx context.Context, userID, channelID int32, isMuted, isDeafened, selfMute, selfDeaf, selfVideo, selfStream *bool) (repo.VoiceState, error) {
	params := repo.UpdateVoiceStateParams{
		UserID:    userID,
		ChannelID: channelID,
	}

	if isMuted != nil {
		params.IsMuted = pgtype.Bool{Bool: *isMuted, Valid: true}
	}
	if isDeafened != nil {
		params.IsDeafened = pgtype.Bool{Bool: *isDeafened, Valid: true}
	}
	if selfMute != nil {
		params.SelfMute = pgtype.Bool{Bool: *selfMute, Valid: true}
	}
	if selfDeaf != nil {
		params.SelfDeaf = pgtype.Bool{Bool: *selfDeaf, Valid: true}
	}
	if selfVideo != nil {
		params.SelfVideo = pgtype.Bool{Bool: *selfVideo, Valid: true}
	}
	if selfStream != nil {
		params.SelfStream = pgtype.Bool{Bool: *selfStream, Valid: true}
	}

	return r.q.UpdateVoiceState(ctx, params)
}

// DeleteVoiceState removes a voice state
func (r *VoiceRepository) DeleteVoiceState(ctx context.Context, userID, channelID int32) error {
	return r.q.DeleteVoiceState(ctx, repo.DeleteVoiceStateParams{
		UserID:    userID,
		ChannelID: channelID,
	})
}

// DeleteUserVoiceStates removes all voice states for a user
func (r *VoiceRepository) DeleteUserVoiceStates(ctx context.Context, userID int32) error {
	return r.q.DeleteUserVoiceStates(ctx, userID)
}
