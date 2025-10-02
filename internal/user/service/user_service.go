package service

import (
	"context"

	user "discord/gen/proto/service/user"
	"discord/internal/user/repository"
)

type UserService interface {
	GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error)
	UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error)
	GetUserProfile(ctx context.Context, req *user.GetUserProfileRequest) (*user.GetUserProfileResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService() UserService {
	userRepo := repository.NewUserRepository()
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	return s.userRepo.GetUser(ctx, req)
}

func (s *userService) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	return s.userRepo.UpdateUser(ctx, req)
}

func (s *userService) GetUserProfile(ctx context.Context, req *user.GetUserProfileRequest) (*user.GetUserProfileResponse, error) {
	return s.userRepo.GetUserProfile(ctx, req)
}
