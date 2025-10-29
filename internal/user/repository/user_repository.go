package repository

import (
	"context"
	"time"

	"discord/gen/repo"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
	q  *repo.Queries
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
		q:  repo.New(db),
	}
}

// User Management
func (r *UserRepository) GetUser(ctx context.Context, userID int32) (repo.User, error) {
	return r.q.GetUserByID(ctx, userID)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (repo.User, error) {
	return r.q.GetUserByUsername(ctx, username)
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID int32, fullName, profilePic, bio, colorCode, backgroundColor, backgroundPic *string) (repo.User, error) {
	params := repo.UpdateUserParams{
		ID: userID,
	}

	if fullName != nil {
		params.FullName = pgtype.Text{String: *fullName, Valid: true}
	}
	if profilePic != nil {
		params.ProfilePic = pgtype.Text{String: *profilePic, Valid: true}
	}
	if bio != nil {
		params.Bio = pgtype.Text{String: *bio, Valid: true}
	}
	if colorCode != nil {
		params.ColorCode = pgtype.Text{String: *colorCode, Valid: true}
	}
	if backgroundColor != nil {
		params.BackgroundColor = pgtype.Text{String: *backgroundColor, Valid: true}
	}
	if backgroundPic != nil {
		params.BackgroundPic = pgtype.Text{String: *backgroundPic, Valid: true}
	}

	return r.q.UpdateUser(ctx, params)
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int32) error {
	return r.q.DeleteUser(ctx, userID)
}

func (r *UserRepository) SearchUsers(ctx context.Context, query string, limit, offset int32) ([]repo.User, error) {
	return r.q.SearchUsers(ctx, repo.SearchUsersParams{
		Column1: pgtype.Text{String: query, Valid: true},
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int32) ([]repo.User, error) {
	return r.q.ListUsers(ctx, repo.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID int32, hashedPassword string) error {
	return r.q.UpdateUserPassword(ctx, repo.UpdateUserPasswordParams{
		ID:       userID,
		Password: hashedPassword,
	})
}

// User Presence
func (r *UserRepository) GetUserPresence(ctx context.Context, userID int32) (repo.UserPresence, error) {
	return r.q.GetUserPresence(ctx, userID)
}

func (r *UserRepository) UpdatePresenceStatus(ctx context.Context, userID int32, status string) error {
	return r.q.UpdatePresenceStatus(ctx, repo.UpdatePresenceStatusParams{
		UserID: userID,
		Status: pgtype.Text{String: status, Valid: true},
	})
}

func (r *UserRepository) UpsertUserPresence(ctx context.Context, userID int32, status, customStatus, customStatusEmoji, activity *string) (repo.UserPresence, error) {
	params := repo.UpsertUserPresenceParams{
		UserID: userID,
	}

	if status != nil {
		params.Status = pgtype.Text{String: *status, Valid: true}
	}
	if customStatus != nil {
		params.CustomStatus = pgtype.Text{String: *customStatus, Valid: true}
	}
	if customStatusEmoji != nil {
		params.CustomStatusEmoji = pgtype.Text{String: *customStatusEmoji, Valid: true}
	}
	if activity != nil {
		params.Activity = pgtype.Text{String: *activity, Valid: true}
	}

	return r.q.UpsertUserPresence(ctx, params)
}

func (r *UserRepository) SetCustomStatus(ctx context.Context, userID int32, customStatus, emoji *string, expiresAt *time.Time) error {
	params := repo.SetCustomStatusParams{
		UserID: userID,
	}

	if customStatus != nil {
		params.CustomStatus = pgtype.Text{String: *customStatus, Valid: true}
	}
	if emoji != nil {
		params.CustomStatusEmoji = pgtype.Text{String: *emoji, Valid: true}
	}
	if expiresAt != nil {
		params.CustomStatusExpiresAt = pgtype.Timestamp{Time: *expiresAt, Valid: true}
	}

	return r.q.SetCustomStatus(ctx, params)
}

func (r *UserRepository) GetMultipleUserPresences(ctx context.Context, userIDs []int32) ([]repo.UserPresence, error) {
	return r.q.GetMultipleUserPresences(ctx, userIDs)
}

func (r *UserRepository) ClearExpiredCustomStatuses(ctx context.Context) error {
	return r.q.ClearExpiredCustomStatuses(ctx)
}

// User Blocking (using friends table)
func (r *UserRepository) BlockUser(ctx context.Context, userID, blockedUserID int32) error {
	return r.q.BlockUser(ctx, repo.BlockUserParams{
		UserID:   userID,
		FriendID: blockedUserID,
	})
}

func (r *UserRepository) UnblockUser(ctx context.Context, userID, blockedUserID int32) error {
	// Update status back to 'accepted' or remove the relationship
	return r.q.UpdateFriendStatus(ctx, repo.UpdateFriendStatusParams{
		UserID:   userID,
		FriendID: blockedUserID,
		Status:   "none", // or "accepted" if they were friends before
	})
}

func (r *UserRepository) GetBlockedUsers(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return r.q.GetBlockedUsers(ctx, userID)
}

func (r *UserRepository) UpdateUserStatus(ctx context.Context, userID int32, status string) error {
	return r.q.UpdateUserStatus(ctx, repo.UpdateUserStatusParams{
		ID:     userID,
		Status: status,
	})
}
