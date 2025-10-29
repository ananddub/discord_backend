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

// SendFriendRequest sends a friend request
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

// AcceptFriendRequest accepts a friend request
func (c *FriendController) AcceptFriendRequest(ctx context.Context, req *friendPb.AcceptFriendRequestRequest) (*friendPb.AcceptFriendRequestResponse, error) {
	// Get user ID from context
	userId := ctx.Value("user_id").(int32)

	if req.GetFriendRequestId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// friend_request_id is the requester's user_id who sent the friend request
	requesterID := req.GetFriendRequestId()

	err := c.friendService.AcceptFriendRequest(ctx, userId, requesterID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.AcceptFriendRequestResponse{
		Success: true,
	}, nil
}

// RejectFriendRequest rejects a friend request
func (c *FriendController) RejectFriendRequest(ctx context.Context, req *friendPb.RejectFriendRequestRequest) (*friendPb.RejectFriendRequestResponse, error) {
	// Get user ID from context
	userId := ctx.Value("user_id").(int32)

	if req.GetFriendRequestId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// friend_request_id is the requester's user_id who sent the friend request
	requesterID := req.GetFriendRequestId()

	err := c.friendService.RejectFriendRequest(ctx, userId, requesterID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.RejectFriendRequestResponse{
		Success: true,
	}, nil
}

// BlockFriend blocks a user
func (c *FriendController) BlockFriend(ctx context.Context, req *friendPb.BlockRequest) (*friendPb.BlockResponse, error) {
	if req.GetUserId() == 0 || req.GetFriendId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.friendService.BlockUser(ctx, req.GetUserId(), req.GetFriendId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.BlockResponse{
		Success: true,
	}, nil
}

// UnblockFriend unblocks a user
func (c *FriendController) UnblockFriend(ctx context.Context, req *friendPb.UnblockRequest) (*friendPb.UnblockResponse, error) {
	if req.GetUserId() == 0 || req.GetFriendId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.friendService.UnblockUser(ctx, req.GetUserId(), req.GetFriendId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.UnblockResponse{
		Success: true,
	}, nil
}

// UpdateAliasName updates the alias name for a friend
func (c *FriendController) UpdateAliasName(ctx context.Context, req *friendPb.UpdateAliasNameRequest) (*friendPb.UpdateAliasNameResponse, error) {
	if req.GetUserId() == 0 || req.GetFriendId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.friendService.UpdateAliasName(ctx, req.GetUserId(), req.GetFriendId(), req.GetAliasName())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &friendPb.UpdateAliasNameResponse{
		Success: true,
	}, nil
}

// ListFreinds retrieves all friends for a user (streaming response)
func (c *FriendController) ListFreinds(req *friendPb.GetFriendsRequest, stream friendPb.FriendService_ListFreindsServer) error {
	ctx := stream.Context()

	if req.GetUserId() == 0 {
		return commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	friends, err := c.friendService.GetUserFriends(ctx, req.GetUserId())
	if err != nil {
		return commonErrors.ToGRPCError(err)
	}

	// Convert repo.Friend to proto.Friend and stream
	for _, friend := range friends {
		pbFriend := &schema.Friend{
			Id:         friend.ID,
			UserId:     friend.UserID,
			FriendId:   friend.FriendID,
			Status:     friend.Status,
			IsAccepted: friend.Status == "accepted",
			IsBlocked:  friend.Status == "blocked",
			IsPending:  friend.Status == "pending",
			IsFavorite: friend.IsFavorite.Bool,
			CreatedAt:  friend.CreatedAt.Time.Unix(),
			UpdatedAt:  friend.UpdatedAt.Time.Unix(),
		}

		if friend.AliasName.Valid {
			pbFriend.AliasName = &friend.AliasName.String
		}

		if err := stream.Send(&friendPb.GetFriendsResponse{
			Friends: []*schema.Friend{pbFriend},
		}); err != nil {
			return err
		}
	}

	return nil
}
