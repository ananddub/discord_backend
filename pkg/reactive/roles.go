package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
)

type RoleReactive struct{}

func (s RoleReactive) ReactRole(pg mypg.QueryData) {
	role := s.convertToRole(pg)
	s.publish(role)
}

func (s RoleReactive) publish(data *schema.Role) {
	pub := pubsub.Get()
	topic := "role:" + strconv.Itoa(int(data.Id))
	pub.Publish(topic, data)
}

func (s RoleReactive) convertToRole(pg mypg.QueryData) *schema.Role {
	var value repo.Role
	pg.Data.Scan(
		&value.ID,
		&value.ServerID,
		&value.Name,
		&value.Color,
		&value.Hoist,
		&value.Position,
		&value.Permissions,
		&value.Mentionable,
		&value.Icon,
		&value.Description,
		&value.IsDefault,
		&value.IsDeleted,
		&value.CreatedAt,
		&value.UpdatedAt,
	)
	return &schema.Role{
		Id:          value.ID,
		ServerId:    value.ServerID,
		Name:        value.Name,
		Color:       value.Color.String,
		Hoist:       value.Hoist.Bool,
		Position:    value.Position.Int32,
		Permissions: value.Permissions.Int64,
		Mentionable: value.Mentionable.Bool,
		Icon:        value.Icon.String,
		Description: value.Description.String,
		IsDeleted:   value.IsDeleted.Bool,
		CreatedAt:   value.CreatedAt.Time.Unix(),
		UpdatedAt:   value.UpdatedAt.Time.Unix(),
		Operation:   &pg.Type,
	}
}
