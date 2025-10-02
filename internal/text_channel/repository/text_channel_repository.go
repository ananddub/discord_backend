package repository

import (
	"context"
	"discord/config"
	text_channel "discord/gen/proto/service/text_channel"
	"discord/gen/repo"
)

type TextChannelRepository interface {
	CreateTextGroup(ctx context.Context, req *text_channel.CreateTextGroupRequest) (*text_channel.CreateTextGroupResponse, error)
	CreateTextChannel(ctx context.Context, req *text_channel.CreateTextChannelRequest) (*text_channel.CreateTextChannelResponse, error)
	ArchiveTextChannel(ctx context.Context, req *text_channel.ArchiveTextChannelRequest) (*text_channel.ArchiveTextChannelResponse, error)
}

type textChannelRepository struct {
	readdb  *repo.Queries
	writedb *repo.Queries
}

func NewTextChannelRepository() TextChannelRepository {
	queries, err := config.RepoQuieries()
	if err != nil {
		panic(err)
	}
	return &textChannelRepository{
		readdb:  repo.New(queries.ReadDb),
		writedb: repo.New(queries.WriteDb),
	}
}

func (r *textChannelRepository) CreateTextGroup(ctx context.Context, req *text_channel.CreateTextGroupRequest) (*text_channel.CreateTextGroupResponse, error) {
	return &text_channel.CreateTextGroupResponse{}, nil
}

func (r *textChannelRepository) CreateTextChannel(ctx context.Context, req *text_channel.CreateTextChannelRequest) (*text_channel.CreateTextChannelResponse, error) {
	return &text_channel.CreateTextChannelResponse{}, nil
}

func (r *textChannelRepository) ArchiveTextChannel(ctx context.Context, req *text_channel.ArchiveTextChannelRequest) (*text_channel.ArchiveTextChannelResponse, error) {
	return &text_channel.ArchiveTextChannelResponse{}, nil
}
