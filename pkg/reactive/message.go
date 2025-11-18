package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
	"strings"
)

type MessageRective struct{}

func (s MessageRective) ReactMessage(pg mypg.QueryData) {
	pgtype := strings.ToLower(pg.Table)
	switch pgtype {
	case "insert":
		user := s.convertToUser(pg)
		s.publish(user)
	case "update":
		user := s.convertToUser(pg)
		s.publish(user)
	case "delete":
		user := s.convertToUser(pg)
		s.publish(user)
	}
}

func (s MessageRective) publish(data *schema.Message) {
	pub := pubsub.Get()
	topic := "message:" + strconv.Itoa(int(data.ChannelId))
	pub.Publish(topic, data)
	// topc := "message:" + strconv.Itoa(int(data.UserId))
	// pub.Publish(topc, data)
}

func (s MessageRective) convertToUser(pg mypg.QueryData) *schema.Message {
	var value repo.Message
	pg.Data.Scan(
		&value.ID,
		&value.ReceiverID,
		&value.SenderID,
		&value.Content,
		&value.MessageType,
		&value.ReplyToMessageID,
		&value.MentionEveryone,
		&value.CreatedAt,
		&value.UpdatedAt,
	)
	return &schema.Message{
		ChannelId:  value.ID,
		ReceiverId: value.ReceiverID.Int32,
		SenderId:   value.SenderID,
		Content:    value.Content,
		ReplyToMessageId: func() int32 {
			if value.ReplyToMessageID.Valid {
				return value.ReplyToMessageID.Int32
			}
			return -1
		}(),
		MentionEveryone: value.MentionEveryone.Bool,
		CreatedAt:       value.CreatedAt.Time.Unix(),
		UpdatedAt:       value.UpdatedAt.Time.Unix(),
		Operation:       &pg.Type,
	}

}
