package repository

import (
	"context"
	"discord/config"
	channel "discord/gen/proto/service/channel"
	"discord/gen/repo"
)

type ChannelRepository interface {
	CreateChannel(ctx context.Context, req *channel.CreateChannelRequest) (*channel.CreateChannelResponse, error)
	GetChannel(ctx context.Context, req *channel.GetChannelRequest) (*channel.GetChannelResponse, error)
}

type channelRepository struct {
	readdb  *repo.Queries
	writedb *repo.Queries
}

func NewChannelRepository() ChannelRepository {
	queries, err := config.RepoQuieries()
	if err != nil {
		panic(err)
	}
	return &channelRepository{
		readdb:  repo.New(queries.ReadDb),
		writedb: repo.New(queries.WriteDb),
	}
}

func (r *channelRepository) CreateChannel(ctx context.Context, req *channel.CreateChannelRequest) (*channel.CreateChannelResponse, error) {
	return &channel.CreateChannelResponse{}, nil
}

func (r *channelRepository) GetChannel(ctx context.Context, req *channel.GetChannelRequest) (*channel.GetChannelResponse, error) {
	return &channel.GetChannelResponse{}, nil
}
