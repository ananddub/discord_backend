package service

import (
	"context"

	text_channel "discord/gen/proto/service/text_channel"
	"discord/internal/text_channel/repository"
)

type TextChannelService interface {
	CreateTextGroup(ctx context.Context, req *text_channel.CreateTextGroupRequest) (*text_channel.CreateTextGroupResponse, error)
	CreateTextChannel(ctx context.Context, req *text_channel.CreateTextChannelRequest) (*text_channel.CreateTextChannelResponse, error)
	ArchiveTextChannel(ctx context.Context, req *text_channel.ArchiveTextChannelRequest) (*text_channel.ArchiveTextChannelResponse, error)
}

type textChannelService struct {
	textChannelRepo repository.TextChannelRepository
}

func NewTextChannelService() TextChannelService {
	textChannelRepo := repository.NewTextChannelRepository()
	return &textChannelService{textChannelRepo: textChannelRepo}
}

func (s *textChannelService) CreateTextGroup(ctx context.Context, req *text_channel.CreateTextGroupRequest) (*text_channel.CreateTextGroupResponse, error) {
	return s.textChannelRepo.CreateTextGroup(ctx, req)
}

func (s *textChannelService) CreateTextChannel(ctx context.Context, req *text_channel.CreateTextChannelRequest) (*text_channel.CreateTextChannelResponse, error) {
	return s.textChannelRepo.CreateTextChannel(ctx, req)
}

func (s *textChannelService) ArchiveTextChannel(ctx context.Context, req *text_channel.ArchiveTextChannelRequest) (*text_channel.ArchiveTextChannelResponse, error) {
	return s.textChannelRepo.ArchiveTextChannel(ctx, req)
}
