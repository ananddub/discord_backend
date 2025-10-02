package handler

import (
	"context"

	text_channel "discord/gen/proto/service/text_channel"
	"discord/internal/text_channel/service"
)

type TextChannelHandler struct {
	text_channel.UnimplementedTextChannelServiceServer
	textChannelService service.TextChannelService
}

func NewTextChannelHandler() *TextChannelHandler {
	return &TextChannelHandler{textChannelService: service.NewTextChannelService()}
}

func (h *TextChannelHandler) CreateTextGroup(ctx context.Context, req *text_channel.CreateTextGroupRequest) (*text_channel.CreateTextGroupResponse, error) {
	return h.textChannelService.CreateTextGroup(ctx, req)
}

func (h *TextChannelHandler) CreateTextChannel(ctx context.Context, req *text_channel.CreateTextChannelRequest) (*text_channel.CreateTextChannelResponse, error) {
	return h.textChannelService.CreateTextChannel(ctx, req)
}

func (h *TextChannelHandler) ArchiveTextChannel(ctx context.Context, req *text_channel.ArchiveTextChannelRequest) (*text_channel.ArchiveTextChannelResponse, error) {
	return h.textChannelService.ArchiveTextChannel(ctx, req)
}
