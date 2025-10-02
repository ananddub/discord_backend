package handler

import (
	"context"

	user "discord/gen/proto/service/user"
	"discord/internal/user/service"
)

type UserHandler struct {
	user.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{userService: service.NewUserService()}
}

func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	return h.userService.GetUser(ctx, req)
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	return h.userService.UpdateUser(ctx, req)
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *user.GetUserProfileRequest) (*user.GetUserProfileResponse, error) {
	return h.userService.GetUserProfile(ctx, req)
}
