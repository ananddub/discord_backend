package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
)

type MessageAttachmentReactive struct{}

func (s MessageAttachmentReactive) ReactMessageAttachment(pg mypg.QueryData) {
	attachment := s.convertToMessageAttachment(pg)
	s.publish(attachment)
}

func (s MessageAttachmentReactive) publish(data *schema.MessageAttachment) {
	pub := pubsub.Get()
	topic := "message_attachment:" + strconv.Itoa(int(data.MessageId))
	pub.Publish(topic, data)
}

func (s MessageAttachmentReactive) convertToMessageAttachment(pg mypg.QueryData) *schema.MessageAttachment {
	var value repo.MessageAttachment
	pg.Data.Scan(
		&value.ID,
		&value.MessageID,
		&value.FileUrl,
		&value.FileName,
		&value.FileType,
		&value.FileSize,
		&value.Width,
		&value.Height,
		&value.IsDeleted,
		&value.CreatedAt,
	)
	return &schema.MessageAttachment{
		Id:        value.ID,
		MessageId: value.MessageID,
		FileUrl:   value.FileUrl,
		FileName:  value.FileName,
		FileType:  value.FileType,
		FileSize:  value.FileSize,
		Width:     value.Width.Int32,
		Height:    value.Height.Int32,
		IsDeleted: value.IsDeleted.Bool,
		CreatedAt: value.CreatedAt.Time.Unix(),
		Operation: &pg.Type,
	}
}
