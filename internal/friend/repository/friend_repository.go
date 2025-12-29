package repository

import (
	"context"

	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FriendRepository struct {
	db      *pgxpool.Pool
	queries *repo.Queries
}

func NewFriendRepository(db *pgxpool.Pool) *FriendRepository {
	return &FriendRepository{
		db:      db,
		queries: repo.New(db),
	}
}

// CreateFriendship creates a bidirectional friendship
func (r *FriendRepository) CreateFriendship(ctx context.Context, userID, friendID int32, status string) error {
	// Create friendship from user to friend
	_, err := r.queries.CreateFriend(ctx, repo.CreateFriendParams{
		UserID:   userID,
		FriendID: friendID,
	})
	if err != nil {
		return err
	}

	// Create reverse friendship (friend to user)
	_, err = r.queries.CreateFriend(ctx, repo.CreateFriendParams{
		UserID:   friendID,
		FriendID: userID,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetFriendship retrieves a friendship
func (r *FriendRepository) GetFriendship(ctx context.Context, userID, friendID int32) (repo.Friend, error) {
	return r.queries.GetFriendship(ctx, repo.GetFriendshipParams{
		UserID:   userID,
		FriendID: friendID,
	})
}

// GetUserFriends retrieves all accepted friends for a user
func (r *FriendRepository) GetUserFriends(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return r.queries.GetFriendsWithFlags(ctx, repo.GetFriendsWithFlagsParams{
		UserID: userID,
	})
}

func (r *FriendRepository) GetAcceptedFriends(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return r.queries.GetAcceptedFriends(ctx, userID)
}

// GetPendingFriendRequests retrieves pending friend requests received by user
func (r *FriendRepository) GetPendingFriendRequests(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return r.queries.GetPendingFriendRequests(ctx, userID)
}

// GetSentFriendRequests retrieves friend requests sent by user
func (r *FriendRepository) GetSentFriendRequests(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return r.queries.GetSentFriendRequests(ctx, userID)
}

// UpdateFriendStatus updates the status of a friendship
func (r *FriendRepository) UpdateFriendStatus(ctx context.Context, userID, friendID int32, status string) error {
	// _, err := r.queries.UpdateFriendStatus(ctx, repo.UpdateFriendStatusParams{
	// 	UserID:   userID,
	// 	FriendID: friendID,
	// 	IsAccepted: pgtype.Bool{
	// 		Bool:  true,
	// 		Valid: true,
	// 	},
	// })
	return nil
}

// UpdateFriendAlias updates the alias name for a friend
func (r *FriendRepository) UpdateFriendAlias(ctx context.Context, userID, friendID int32, aliasName string) error {
	_, err := r.queries.UpdateFriendAlias(ctx, repo.UpdateFriendAliasParams{
		UserID:    userID,
		FriendID:  friendID,
		AliasName: pgtype.Text{String: aliasName, Valid: true},
	})
	return err
}

// ToggleFavorite toggles the favorite status of a friend
func (r *FriendRepository) ToggleFavorite(ctx context.Context, userID, friendID int32, isFavorite bool) error {
	_, err := r.queries.ToggleFavoriteFlag(ctx, repo.ToggleFavoriteFlagParams{
		UserID:   userID,
		FriendID: friendID,
	})
	return err
}

// DeleteFriendship deletes a bidirectional friendship
func (r *FriendRepository) DeleteFriendship(ctx context.Context, userID, friendID int32) error {
	_, err := r.queries.SoftDeleteFriendship(ctx, repo.SoftDeleteFriendshipParams{
		UserID:   userID,
		FriendID: friendID,
	})
	return err
}

// BlockUser blocks a user
func (r *FriendRepository) BlockUser(ctx context.Context, userID, friendID int32) error {
	_, err := r.queries.BlockUser(ctx, repo.BlockUserParams{
		UserID:   userID,
		FriendID: friendID,
	})
	return err
}

// GetBlockedUsers retrieves all blocked users
func (r *FriendRepository) GetBlockedUsers(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return r.queries.GetBlockedUsers(ctx, userID)
}

func (r *FriendRepository) GetUserByID(ctx context.Context, userID int32) (repo.User, error) {
	return r.queries.GetUserByID(ctx, userID)
}

func (r *FriendRepository) AcceptFriendRequest(ctx context.Context, userID, friendID int32) (repo.Friend, error) {
	return r.queries.AcceptFriendRequest(ctx, repo.AcceptFriendRequestParams{
		UserID:   userID,
		FriendID: friendID,
	})
}

func (r *FriendRepository) RejectFriendRequest(ctx context.Context, userID, friendID int32) error {
	_, err := r.queries.RejectFriendRequest(ctx, repo.RejectFriendRequestParams{
		UserID:   userID,
		FriendID: friendID,
	})
	return err
}
