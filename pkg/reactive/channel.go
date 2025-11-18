package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
	"strings"
)

type ChannelReactive struct{}

func (s ChannelReactive) ReactChannel(pg mypg.QueryData) {
	pgtype := strings.ToLower(pg.Type)
	switch pgtype {
	case "insert", "update", "delete":
		channel := s.convertToChannel(pg)
		s.publish(channel)
	}
}

func (s ChannelReactive) publish(data *schema.Channel) {
	pub := pubsub.Get()
	topic := "channel:" + strconv.Itoa(int(data.Id))
	pub.Publish(topic, data)
}

func (s ChannelReactive) convertToChannel(pg mypg.QueryData) *schema.Channel {
	var value repo.Channel
	pg.Data.Scan(
		&value.ID,
		&value.ServerID,
		&value.CategoryID,
		&value.Name,
		&value.Type,
		&value.Position,
		&value.Topic,
		&value.IsNsfw,
		&value.SlowmodeDelay,
		&value.UserLimit,
		&value.Bitrate,
		&value.IsPrivate,
		&value.IsDeleted,
		&value.CreatedAt,
		&value.UpdatedAt,
	)
	return &schema.Channel{
		Id:          value.ID,
		ServerId:    value.ServerID,
		CategoryId:  value.CategoryID.Int32,
		Name:        value.Name,
		Position:    value.Position.Int32,
		Description: value.Topic.String,
		IsNsfw:      value.IsNsfw.Bool,
		IsDeleted:   value.IsDeleted.Bool,
		CreatedAt:   value.CreatedAt.Time.Unix(),
		UpdatedAt:   value.UpdatedAt.Time.Unix(),
		Operation:   &pg.Type,
	}
}
