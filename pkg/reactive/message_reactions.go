package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
)

type MessageReactionReactive struct{}

func (s MessageReactionReactive) ReactMessageReaction(pg mypg.QueryData) {
	reaction := s.convertToMessageReaction(pg)
	s.publish(reaction)
}

func (s MessageReactionReactive) publish(data *schema.MessageReaction) {
	pub := pubsub.Get()
	topic := "message_reaction:" + strconv.Itoa(int(data.MessageId))
	pub.Publish(topic, data)
}

func (s MessageReactionReactive) convertToMessageReaction(pg mypg.QueryData) *schema.MessageReaction {
	var value repo.MessageReaction
	pg.Data.Scan(
		&value.ID,
		&value.MessageID,
		&value.UserID,
		&value.Emoji,
		&value.EmojiID,
		&value.CreatedAt,
	)
	return &schema.MessageReaction{
		Id:        value.ID,
		MessageId: value.MessageID,
		UserId:    value.UserID,
		Emoji:     value.Emoji,
		EmojiId:   value.EmojiID.String,
		CreatedAt: value.CreatedAt.Time.Unix(),
		Operation: &pg.Type,
	}
}
