package handler

import (
	"context"

	channel "discord/gen/proto/service/channel"
	"discord/internal/channel/service"
)

type ChannelHandler struct {
	channel.UnimplementedChannelServiceServer
	channelService service.ChannelService
}

func NewChannelHandler() *ChannelHandler {
	return &ChannelHandler{channelService: service.NewChannelService()}
}

func (h *ChannelHandler) CreateChannel(ctx context.Context, req *channel.CreateChannelRequest) (*channel.CreateChannelResponse, error) {
	return h.channelService.CreateChannel(ctx, req)
}

func (h *ChannelHandler) GetChannel(ctx context.Context, req *channel.GetChannelRequest) (*channel.GetChannelResponse, error) {
	return h.channelService.GetChannel(ctx, req)
}
