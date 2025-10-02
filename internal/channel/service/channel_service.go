package service

import (
	"context"

	channel "discord/gen/proto/service/channel"
	"discord/internal/channel/repository"
)

type ChannelService interface {
	CreateChannel(ctx context.Context, req *channel.CreateChannelRequest) (*channel.CreateChannelResponse, error)
	GetChannel(ctx context.Context, req *channel.GetChannelRequest) (*channel.GetChannelResponse, error)
}

type channelService struct {
	channelRepo repository.ChannelRepository
}

func NewChannelService() ChannelService {
	channelRepo := repository.NewChannelRepository()
	return &channelService{channelRepo: channelRepo}
}

func (s *channelService) CreateChannel(ctx context.Context, req *channel.CreateChannelRequest) (*channel.CreateChannelResponse, error) {
	return s.channelRepo.CreateChannel(ctx, req)
}

func (s *channelService) GetChannel(ctx context.Context, req *channel.GetChannelRequest) (*channel.GetChannelResponse, error) {
	return s.channelRepo.GetChannel(ctx, req)
}
