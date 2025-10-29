package service

import (
	"context"
	"errors"

	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	friendRepo "discord/internal/friend/repository"
)

type FriendService struct {
	friendRepo *friendRepo.FriendRepository
}

func NewFriendService(friendRepo *friendRepo.FriendRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
	}
}

// SendFriendRequest sends a friend request
func (s *FriendService) SendFriendRequest(ctx context.Context, userID, friendID int32) error {
	// Validate input
	if userID == friendID {
		return commonErrors.ErrInvalidInput
	}

	// Check if friendship already exists
	existing, err := s.friendRepo.GetFriendship(ctx, userID, friendID)
	if err == nil {
		// Friendship exists
		if existing.Status == "accepted" {
			return errors.New("already friends")
		}
		if existing.Status == "pending" {
			return errors.New("friend request already sent")
		}
		if existing.Status == "blocked" {
			return errors.New("user is blocked")
		}
	}

	// Create friend request with pending status
	return s.friendRepo.CreateFriendship(ctx, userID, friendID, "pending")
}

// AcceptFriendRequest accepts a friend request
func (s *FriendService) AcceptFriendRequest(ctx context.Context, userID, requesterID int32) error {
	// Check if friend request exists
	friendship, err := s.friendRepo.GetFriendship(ctx, requesterID, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if friendship.Status != "pending" {
		return errors.New("no pending friend request")
	}

	// Update status to accepted for both sides
	err = s.friendRepo.UpdateFriendStatus(ctx, requesterID, userID, "accepted")
	if err != nil {
		return err
	}

	// Create reverse friendship
	return s.friendRepo.CreateFriendship(ctx, userID, requesterID, "accepted")
}

// RejectFriendRequest rejects a friend request
func (s *FriendService) RejectFriendRequest(ctx context.Context, userID, requesterID int32) error {
	// Check if friend request exists
	friendship, err := s.friendRepo.GetFriendship(ctx, requesterID, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if friendship.Status != "pending" {
		return errors.New("no pending friend request")
	}

	// Delete the friend request
	return s.friendRepo.DeleteFriendship(ctx, requesterID, userID)
}

// RemoveFriend removes a friend
func (s *FriendService) RemoveFriend(ctx context.Context, userID, friendID int32) error {
	// Check if friendship exists
	_, err := s.friendRepo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Delete bidirectional friendship
	return s.friendRepo.DeleteFriendship(ctx, userID, friendID)
}

// BlockUser blocks a user
func (s *FriendService) BlockUser(ctx context.Context, userID, targetID int32) error {
	// Check if already blocked
	existing, err := s.friendRepo.GetFriendship(ctx, userID, targetID)
	if err == nil && existing.Status == "blocked" {
		return errors.New("user already blocked")
	}

	// Delete any existing friendship first
	_ = s.friendRepo.DeleteFriendship(ctx, userID, targetID)

	// Create blocked status
	return s.friendRepo.CreateFriendship(ctx, userID, targetID, "blocked")
}

// UnblockUser unblocks a user
func (s *FriendService) UnblockUser(ctx context.Context, userID, targetID int32) error {
	// Check if blocked
	existing, err := s.friendRepo.GetFriendship(ctx, userID, targetID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if existing.Status != "blocked" {
		return errors.New("user is not blocked")
	}

	// Delete the block
	return s.friendRepo.DeleteFriendship(ctx, userID, targetID)
}

// UpdateAliasName updates the alias name for a friend
func (s *FriendService) UpdateAliasName(ctx context.Context, userID, friendID int32, aliasName string) error {
	// Check if friendship exists
	friendship, err := s.friendRepo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if friendship.Status != "accepted" {
		return errors.New("not friends")
	}

	return s.friendRepo.UpdateFriendAlias(ctx, userID, friendID, aliasName)
}

// GetUserFriends retrieves all friends for a user
func (s *FriendService) GetUserFriends(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return s.friendRepo.GetUserFriends(ctx, userID)
}

// GetPendingFriendRequests retrieves pending friend requests
func (s *FriendService) GetPendingFriendRequests(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return s.friendRepo.GetPendingFriendRequests(ctx, userID)
}

// GetSentFriendRequests retrieves sent friend requests
func (s *FriendService) GetSentFriendRequests(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return s.friendRepo.GetSentFriendRequests(ctx, userID)
}

// GetBlockedUsers retrieves all blocked users
func (s *FriendService) GetBlockedUsers(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return s.friendRepo.GetBlockedUsers(ctx, userID)
}
