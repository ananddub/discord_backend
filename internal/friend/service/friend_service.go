package service

import (
	"context"
	"errors"

	"discord/gen/proto/schema"
	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	friendRepo "discord/internal/friend/repository"
	"discord/pkg/pubsub"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type FriendService struct {
	friendRepo *friendRepo.FriendRepository
	pubsub     *pubsub.PubSub
	grpcClient *grpc.ClientConn
}

func NewFriendService(friendRepo *friendRepo.FriendRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		pubsub:     pubsub.Get(),
	}
}

func (s *FriendService) SendFriendRequest(ctx context.Context, userID, friendID int32) error {
	if userID == friendID {
		return commonErrors.ErrInvalidInput
	}

	existing, err := s.friendRepo.GetFriendship(ctx, userID, friendID)
	if err != nil && err != pgx.ErrNoRows {
		return errors.New("user not found")
	}
	if existing.IsAccepted.Bool {
		return errors.New("already friends")
	}
	if existing.IsPending.Bool {
		return errors.New("friend request already sent")
	}
	if existing.IsBlocked.Bool {
		return errors.New("user is blocked")
	}

	err = s.friendRepo.CreateFriendship(ctx, userID, friendID, "pending")
	if err != nil {
		return err
	}
	user, err := s.friendRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	userschema := schema.User{
		Id:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		ProfilePic:      user.ProfilePic.String,
		BackgroundPic:   user.BackgroundPic.String,
		ColorCode:       user.ColorCode.String,
		BackgroundColor: user.BackgroundColor.String,
		Bio:             user.Bio.String,
		Status:          user.Status,
		CustomStatus:    user.CustomStatus.String,
		IsBot:           user.IsBot.Bool,
		IsVerified:      user.IsVerified.Valid,
		IsDeleted:       user.IsDeleted.Bool,
	}
	friendRequest := &schema.Friend{
		UserId:     userID,
		FriendId:   friendID,
		IsPending:  true,
		IsAccepted: false,
		Status:     "pending",
		User:       &userschema,
	}

	//reactive.ChatServiceClient.CreateUser(ctx,&schema2.User{..userschema}):
	s.publishToFriend(friendID, friendRequest)

	return nil
}

func (s *FriendService) AcceptFriendRequest(ctx context.Context, userID, requesterID int32) error {
	// Check if friend request exists
	friendship, err := s.friendRepo.GetFriendship(ctx, requesterID, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if friendship.IsPending.Bool != true {
		return errors.New("no pending friend request")
	}

	// Update status to accepted for both sides
	friendship, err = s.friendRepo.AcceptFriendRequest(ctx, userID, requesterID)
	if err != nil {
		return err
	}
	user, err := s.friendRepo.GetUserByID(ctx, requesterID)
	if err != nil {
		return err
	}

	acceptedFriend := &schema.Friend{
		UserId:     requesterID,
		FriendId:   userID,
		IsPending:  false,
		IsAccepted: true,
		Status:     "accepted",
		User: &schema.User{
			Id:              user.ID,
			Username:        user.Username,
			Email:           user.Email,
			ProfilePic:      user.ProfilePic.String,
			BackgroundPic:   user.BackgroundPic.String,
			ColorCode:       user.ColorCode.String,
			BackgroundColor: user.BackgroundColor.String,
			Bio:             user.Bio.String,
			Status:          user.Status,
			CustomStatus:    user.CustomStatus.String,
			IsBot:           user.IsBot.Bool,
			IsVerified:      user.IsVerified.Valid,
			IsDeleted:       user.IsDeleted.Bool,
		},
	}

	s.publishToFriend(requesterID, acceptedFriend)
	user, err = s.friendRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	reverseFriend := &schema.Friend{
		UserId:     userID,
		FriendId:   requesterID,
		IsPending:  false,
		IsAccepted: true,
		Status:     "accepted",
		User: &schema.User{
			Id:              user.ID,
			Username:        user.Username,
			Email:           user.Email,
			ProfilePic:      user.ProfilePic.String,
			BackgroundPic:   user.BackgroundPic.String,
			ColorCode:       user.ColorCode.String,
			BackgroundColor: user.BackgroundColor.String,
			Bio:             user.Bio.String,
			Status:          user.Status,
			CustomStatus:    user.CustomStatus.String,
			IsBot:           user.IsBot.Bool,
			IsVerified:      user.IsVerified.Valid,
			IsDeleted:       user.IsDeleted.Bool,
		},
	}
	s.publishToFriend(userID, reverseFriend)

	return nil
}

func (s *FriendService) RejectFriendRequest(ctx context.Context, userID, requesterID int32) error {
	// Check if friend request exists
	friendship, err := s.friendRepo.GetFriendship(ctx, requesterID, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if friendship.IsPending.Bool != true {
		return errors.New("no pending friend request")
	}

	err = s.friendRepo.DeleteFriendship(ctx, requesterID, userID)
	if err != nil {
		return err
	}

	rejectedFriend := &schema.Friend{
		UserId:     requesterID,
		FriendId:   userID,
		IsPending:  false,
		IsAccepted: false,
		IsRejected: true,
		Status:     "rejected",
	}

	s.publishToFriend(requesterID, rejectedFriend)

	return nil
}

func (s *FriendService) RemoveFriend(ctx context.Context, userID, friendID int32) error {
	// Check if friendship exists
	_, err := s.friendRepo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	err = s.friendRepo.DeleteFriendship(ctx, userID, friendID)
	if err != nil {
		return err
	}

	err = s.friendRepo.DeleteFriendship(ctx, friendID, userID)
	if err != nil {
		return err
	}

	removedFriend := &schema.Friend{
		UserId:    userID,
		FriendId:  friendID,
		IsDeleted: true,
		Status:    "removed",
	}

	s.publishToFriend(friendID, removedFriend)
	return nil
}

func (s *FriendService) BlockUser(ctx context.Context, userID, targetID int32) error {

	existing, err := s.friendRepo.GetFriendship(ctx, userID, targetID)
	if err == nil && existing.IsBlocked.Bool == true {
		return errors.New("user already blocked")
	}

	// Delete any existing friendship first
	_ = s.friendRepo.DeleteFriendship(ctx, userID, targetID)
	// s.publishToFriend(ctx, targetID, &schema.Friend{
	// 	UserId:    userID,
	// 	FriendId:  targetID,
	// 	IsDeleted: true,
	// 	Status:    "removed",
	// })
	s.publishToFriend(targetID, &schema.Friend{
		UserId:    userID,
		FriendId:  targetID,
		IsDeleted: true,
		Status:    "removed",
	})
	return s.friendRepo.CreateFriendship(ctx, userID, targetID, "blocked")
}

// UnblockUser unblocks a user
func (s *FriendService) UnblockUser(ctx context.Context, userID, targetID int32) error {
	// Check if blocked
	existing, err := s.friendRepo.GetFriendship(ctx, userID, targetID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if existing.IsBlocked.Bool != true {
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

	if friendship.IsAccepted.Bool != true {
		return errors.New("not friends")
	}

	return s.friendRepo.UpdateFriendAlias(ctx, userID, friendID, aliasName)
}

// GetUserFriends retrieves all friends for a user
func (s *FriendService) GetAcceptedFriends(ctx context.Context, userID int32) ([]repo.Friend, error) {
	return s.friendRepo.GetAcceptedFriends(ctx, userID)
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
