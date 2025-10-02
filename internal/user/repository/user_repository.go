package repository

import (
	"context"
	"discord/config"
	user "discord/gen/proto/service/user"
	"discord/gen/repo"
)

type UserRepository interface {
	GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error)
	UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error)
	GetUserProfile(ctx context.Context, req *user.GetUserProfileRequest) (*user.GetUserProfileResponse, error)
}

type userRepository struct {
	readdb  *repo.Queries
	writedb *repo.Queries
}

func NewUserRepository() UserRepository {
	queries, err := config.RepoQuieries()
	if err != nil {
		panic(err)
	}
	return &userRepository{
		readdb:  repo.New(queries.ReadDb),
		writedb: repo.New(queries.WriteDb),
	}
}

func (r *userRepository) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	return &user.GetUserResponse{}, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	return &user.UpdateUserResponse{}, nil
}

func (r *userRepository) GetUserProfile(ctx context.Context, req *user.GetUserProfileRequest) (*user.GetUserProfileResponse, error) {
	return &user.GetUserProfileResponse{}, nil
}
