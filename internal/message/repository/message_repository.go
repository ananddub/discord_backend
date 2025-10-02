package repository

import (
	"context"
	"discord/config"
	message "discord/gen/proto/service/message"
	"discord/gen/repo"
)

type MessageRepository interface {
	SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error)
	GetMessages(ctx context.Context, req *message.GetMessagesRequest) (*message.GetMessagesResponse, error)
}

type messageRepository struct {
	readdb  *repo.Queries
	writedb *repo.Queries
}

func NewMessageRepository() MessageRepository {
	queries, err := config.RepoQuieries()
	if err != nil {
		panic(err)
	}
	return &messageRepository{
		readdb:  repo.New(queries.ReadDb),
		writedb: repo.New(queries.WriteDb),
	}
}

func (r *messageRepository) SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	return &message.SendMessageResponse{}, nil
}

func (r *messageRepository) GetMessages(ctx context.Context, req *message.GetMessagesRequest) (*message.GetMessagesResponse, error) {
	return &message.GetMessagesResponse{}, nil
}
