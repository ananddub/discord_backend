package service

import (
	"context"

	friend "discord/gen/proto/service/friend"
	"discord/internal/friend/repository"
)

type FriendService interface {
	SendFriendRequest(ctx context.Context, req *friend.SendFriendRequestRequest) (*friend.SendFriendRequestResponse, error)
	AcceptFriendRequest(ctx context.Context, req *friend.AcceptFriendRequestRequest) (*friend.AcceptFriendRequestResponse, error)
}

type friendService struct {
	friendRepo repository.FriendRepository
}

func NewFriendService() FriendService {
	friendRepo := repository.NewFriendRepository()
	return &friendService{friendRepo: friendRepo}
}

func (s *friendService) SendFriendRequest(ctx context.Context, req *friend.SendFriendRequestRequest) (*friend.SendFriendRequestResponse, error) {
	return s.friendRepo.SendFriendRequest(ctx, req)
}

func (s *friendService) AcceptFriendRequest(ctx context.Context, req *friend.AcceptFriendRequestRequest) (*friend.AcceptFriendRequestResponse, error) {
	return s.friendRepo.AcceptFriendRequest(ctx, req)
}
