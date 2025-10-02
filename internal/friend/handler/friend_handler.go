package handler

import (
	"context"

	friend "discord/gen/proto/service/friend"
	"discord/internal/friend/service"
)

type FriendHandler struct {
	friend.UnimplementedFriendServiceServer
	friendService service.FriendService
}

func NewFriendHandler() *FriendHandler {
	return &FriendHandler{friendService: service.NewFriendService()}
}

func (h *FriendHandler) SendFriendRequest(ctx context.Context, req *friend.SendFriendRequestRequest) (*friend.SendFriendRequestResponse, error) {
	return h.friendService.SendFriendRequest(ctx, req)
}

func (h *FriendHandler) AcceptFriendRequest(ctx context.Context, req *friend.AcceptFriendRequestRequest) (*friend.AcceptFriendRequestResponse, error) {
	return h.friendService.AcceptFriendRequest(ctx, req)
}
