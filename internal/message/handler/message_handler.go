package handler

import (
	"context"

	message "discord/gen/proto/service/message"
	"discord/internal/message/service"
)

type MessageHandler struct {
	message.UnimplementedMessageServiceServer
	messageService service.MessageService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{messageService: service.NewMessageService()}
}

func (h *MessageHandler) SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	return h.messageService.SendMessage(ctx, req)
}

func (h *MessageHandler) GetMessages(ctx context.Context, req *message.GetMessagesRequest) (*message.GetMessagesResponse, error) {
	return h.messageService.GetMessages(ctx, req)
}
