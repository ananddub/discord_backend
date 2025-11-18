package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
)

type ServerMemberReactive struct{}

func (s ServerMemberReactive) ReactServerMember(pg mypg.QueryData) {
	member := s.convertToServerMember(pg)
	s.publish(member)
}

func (s ServerMemberReactive) publish(data *schema.ServerMember) {
	pub := pubsub.Get()
	topic := "server_member:" + strconv.Itoa(int(data.ServerId))
	pub.Publish(topic, data)
}

func (s ServerMemberReactive) convertToServerMember(pg mypg.QueryData) *schema.ServerMember {
	var value repo.ServerMember
	pg.Data.Scan(
		&value.ID,
		&value.ServerID,
		&value.UserID,
		&value.Nickname,
		&value.JoinedAt,
		&value.IsMuted,
		&value.IsDeafened,
		&value.UpdatedAt,
	)
	return &schema.ServerMember{
		Id:         value.ID,
		ServerId:   value.ServerID,
		UserId:     value.UserID,
		Nickname:   value.Nickname.String,
		JoinedAt:   value.JoinedAt.Time.Unix(),
		IsMuted:    value.IsMuted.Bool,
		IsDeafened: value.IsDeafened.Bool,
	}
}
