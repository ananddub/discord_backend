package reactive

import (
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
)

type ChannelPermissionReactive struct{}

func (s ChannelPermissionReactive) ReactChannelPermission(pg mypg.QueryData) {
	permission := s.convertToChannelPermission(pg)
	s.publish(permission)
}

func (s ChannelPermissionReactive) publish(data *repo.ChannelPermission) {
	pub := pubsub.Get()
	topic := "channel_permission:" + strconv.Itoa(int(data.ChannelID))
	pub.Publish(topic, data)
}

func (s ChannelPermissionReactive) convertToChannelPermission(pg mypg.QueryData) *repo.ChannelPermission {
	var value repo.ChannelPermission
	pg.Data.Scan(
		&value.ID,
		&value.ChannelID,
		&value.RoleID,
		&value.UserID,
		&value.AllowPermissions,
		&value.DenyPermissions,
		&value.CreatedAt,
		&value.UpdatedAt,
	)
	return &value
}
