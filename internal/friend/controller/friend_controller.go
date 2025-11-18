package controller

import (
	"context"

	"discord/gen/proto/schema"
	friendPb "discord/gen/proto/service/friend"
	commonErrors "discord/internal/common/errors"
	friendService "discord/internal/friend/service"
)

type FriendController struct {
	friendPb.UnimplementedFriendServiceServer
	friendService *friendService.FriendService
}

func NewFriendController(friendService *friendService.FriendService) *friendPb.FriendServiceServer {
	controller := &FriendController{
		friendService: friendService,
	}
	var grpcController friendPb.FriendServiceServer = controller
	return &grpcController
}

func (c *FriendController) SendFriendRequest(ctx context.Context, req *friendPb.SendFriendRequestRequest) (*friendPb.SendFriendRequestResponse, error) {
	if req.GetUserId() == 0 || req.GetFriendId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.friendService.SendFriendRequest(ctx, req.GetUserId(), req.GetFriendId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.SendFriendRequestResponse{
		Success: true,
	}, nil
}

func (c *FriendController) AcceptFriendRequest(ctx context.Context, req *friendPb.AcceptFriendRequestRequest) (*friendPb.AcceptFriendRequestResponse, error) {

	userId := ctx.Value("user_id").(int32)

	if req.GetRequesterId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	requesterID := req.GetRequesterId()

	err := c.friendService.AcceptFriendRequest(ctx, userId, requesterID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.AcceptFriendRequestResponse{
		Success: true,
	}, nil
}

func (c *FriendController) RejectFriendRequest(ctx context.Context, req *friendPb.RejectFriendRequestRequest) (*friendPb.RejectFriendRequestResponse, error) {

	userId := ctx.Value("user_id").(int32)

	if req.GetRequesterId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	requesterID := req.GetRequesterId()

	err := c.friendService.RejectFriendRequest(ctx, userId, requesterID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.RejectFriendRequestResponse{
		Success: true,
	}, nil
}

func (c *FriendController) RemoveFriend(ctx context.Context, req *friendPb.RemoveFriendRequest) (*friendPb.RemoveFriendResponse, error) {
	err := c.friendService.RemoveFriend(ctx, req.GetUserId(), req.GetFriendId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}
	return &friendPb.RemoveFriendResponse{Success: true}, nil
}

func (c *FriendController) GetFriends(ctx context.Context, req *friendPb.GetFriendsRequest) (*friendPb.GetFriendsResponse, error) {
	friends, err := c.friendService.GetUserFriends(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	var pbFriends []*schema.Friend
	for _, friend := range friends {
		pbFriends = append(pbFriends, &schema.Friend{
			Id:       friend.ID,
			UserId:   friend.UserID,
			FriendId: friend.FriendID,
			Status:   friend.Status,
		})
	}

	return &friendPb.GetFriendsResponse{Friends: pbFriends}, nil
}

func (c *FriendController) GetPendingRequests(ctx context.Context, req *friendPb.GetPendingRequestsRequest) (*friendPb.GetPendingRequestsResponse, error) {
	requests, err := c.friendService.GetPendingFriendRequests(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	var pbRequests []*schema.Friend
	for _, request := range requests {
		pbRequests = append(pbRequests, &schema.Friend{
			Id:       request.ID,
			UserId:   request.UserID,
			FriendId: request.FriendID,
			Status:   request.Status,
		})
	}

	return &friendPb.GetPendingRequestsResponse{PendingRequests: pbRequests}, nil
}

func (c *FriendController) GetBlockedUsers(ctx context.Context, req *friendPb.GetBlockedUsersRequest) (*friendPb.GetBlockedUsersResponse, error) {
	blocked, err := c.friendService.GetBlockedUsers(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	var pbBlocked []*schema.Friend
	for _, user := range blocked {
		pbBlocked = append(pbBlocked, &schema.Friend{
			Id:       user.ID,
			UserId:   user.UserID,
			FriendId: user.FriendID,
			Status:   user.Status,
		})
	}

	return &friendPb.GetBlockedUsersResponse{BlockedUsers: pbBlocked}, nil
}

func (c *FriendController) SearchFriends(ctx context.Context, req *friendPb.SearchFriendsRequest) (*friendPb.SearchFriendsResponse, error) {
	return &friendPb.SearchFriendsResponse{Friends: []*schema.Friend{}}, nil
}

func (c *FriendController) StreamFriendUpdates(req *friendPb.StreamFriendUpdatesRequest, stream friendPb.FriendService_StreamFriendUpdatesServer) error {
	ch := friendService.Stream(req.GetUserId())
	defer ch.Close()
	for data := range ch.Receive() {
		update, ok := data.(*schema.Friend)
		if !ok {
			continue
		}
		err := stream.Send(update)
		if err != nil {
			return commonErrors.ToGRPCError(err)
		}
	}
	return nil
}
