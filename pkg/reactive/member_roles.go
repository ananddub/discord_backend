package reactive

import (
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
)

type MemberRoleReactive struct{}

func (s MemberRoleReactive) ReactMemberRole(pg mypg.QueryData) {
	memberRole := s.convertToMemberRole(pg)
	s.publish(memberRole)
}

func (s MemberRoleReactive) publish(data *repo.MemberRole) {
	pub := pubsub.Get()
	topic := "member_role:" + strconv.Itoa(int(data.MemberID))
	pub.Publish(topic, data)
}

func (s MemberRoleReactive) convertToMemberRole(pg mypg.QueryData) *repo.MemberRole {
	var value repo.MemberRole
	pg.Data.Scan(
		&value.ID,
		&value.MemberID,
		&value.RoleID,
		&value.AssignedAt,
		&value.UpdatedAt,
	)
	return &value
}
