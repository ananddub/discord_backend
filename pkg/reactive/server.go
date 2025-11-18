package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
	"strings"
)

type ServerReactive struct{}

func (s ServerReactive) ReactServer(pg mypg.QueryData) {
	pgtype := strings.ToLower(pg.Type)
	switch pgtype {
	case "insert", "update", "delete":
		server := s.convertToServer(pg)
		s.publish(server)
	}
}

func (s ServerReactive) publish(data *schema.Server) {
	pub := pubsub.Get()
	topic := "server:" + strconv.Itoa(int(data.Id))
	pub.Publish(topic, data)
}

func (s ServerReactive) convertToServer(pg mypg.QueryData) *schema.Server {
	var value repo.Server
	pg.Data.Scan(
		&value.ID,
		&value.Name,
		&value.Icon,
		&value.Banner,
		&value.Description,
		&value.OwnerID,
		&value.Region,
		&value.MemberCount,
		&value.IsVerified,
		&value.VanityUrl,
		&value.IsDeleted,
		&value.CreatedAt,
		&value.UpdatedAt,
	)
	return &schema.Server{
		Id:          value.ID,
		Name:        value.Name,
		Icon:        value.Icon.String,
		Banner:      value.Banner.String,
		Description: value.Description.String,
		OwnerId:     value.OwnerID,
		Region:      value.Region.String,
		MemberCount: value.MemberCount.Int32,
		IsVerified:  value.IsVerified.Bool,
		IsDeleted:   value.IsDeleted.Bool,
		CreatedAt:   value.CreatedAt.Time.Unix(),
		UpdatedAt:   value.UpdatedAt.Time.Unix(),
		Operation:   &pg.Type,
	}
}
