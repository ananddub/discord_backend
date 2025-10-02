package repository

import (
	"context"
	"discord/config"
	friend "discord/gen/proto/service/friend"
	"discord/gen/repo"
)

type FriendRepository interface {
	SendFriendRequest(ctx context.Context, req *friend.SendFriendRequestRequest) (*friend.SendFriendRequestResponse, error)
	AcceptFriendRequest(ctx context.Context, req *friend.AcceptFriendRequestRequest) (*friend.AcceptFriendRequestResponse, error)
}

type friendRepository struct {
	readdb  *repo.Queries
	writedb *repo.Queries
}

func NewFriendRepository() FriendRepository {
	queries, err := config.RepoQuieries()
	if err != nil {
		panic(err)
	}
	return &friendRepository{
		readdb:  repo.New(queries.ReadDb),
		writedb: repo.New(queries.WriteDb),
	}
}

func (r *friendRepository) SendFriendRequest(ctx context.Context, req *friend.SendFriendRequestRequest) (*friend.SendFriendRequestResponse, error) {
	return &friend.SendFriendRequestResponse{}, nil
}

func (r *friendRepository) AcceptFriendRequest(ctx context.Context, req *friend.AcceptFriendRequestRequest) (*friend.AcceptFriendRequestResponse, error) {
	return &friend.AcceptFriendRequestResponse{}, nil
}
