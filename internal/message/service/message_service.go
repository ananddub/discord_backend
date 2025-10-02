package service

import (
	"context"

	message "discord/gen/proto/service/message"
	"discord/internal/message/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error)
	GetMessages(ctx context.Context, req *message.GetMessagesRequest) (*message.GetMessagesResponse, error)
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService() MessageService {
	messageRepo := repository.NewMessageRepository()
	return &messageService{messageRepo: messageRepo}
}

func (s *messageService) SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	return s.messageRepo.SendMessage(ctx, req)
}

func (s *messageService) GetMessages(ctx context.Context, req *message.GetMessagesRequest) (*message.GetMessagesResponse, error) {
	return s.messageRepo.GetMessages(ctx, req)
}
