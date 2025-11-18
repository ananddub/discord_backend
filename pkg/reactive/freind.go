package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
	"strings"
)

type FriendReactive struct{}

func (s FriendReactive) ReactFriend(pg mypg.QueryData) {
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

func (s FriendReactive) publish(data *schema.Friend) {
	pub := pubsub.Get()
	topic := "friend:" + strconv.Itoa(int(data.FriendId))
	pub.Publish(topic, data)
	topc := "friend:" + strconv.Itoa(int(data.UserId))
	pub.Publish(topc, data)
}

func (s FriendReactive) convertToUser(pg mypg.QueryData) *schema.Friend {
	var value repo.Friend
	pg.Data.Scan(
		&value.ID,
		&value.UserID,
		&value.FriendID,
		&value.AliasName,
		&value.IsDeleted,
		&value.IsFavorite,
		&value.IsBlocked,
		&value.CreatedAt,
		&value.UpdatedAt,
		&value.IsPending,
		&value.IsAccepted,
		&value.IsMuted,
	)
	return &schema.Friend{
		Id:         value.ID,
		UserId:     value.UserID,
		FriendId:   value.FriendID,
		AliasName:  &value.AliasName.String,
		IsDeleted:  value.IsDeleted.Bool,
		IsFavorite: value.IsFavorite.Bool,
		IsBlocked:  value.IsBlocked.Bool,
		CreatedAt:  value.CreatedAt.Time.Unix(),
		UpdatedAt:  value.UpdatedAt.Time.Unix(),
		IsPending:  value.IsPending.Bool,
		IsAccepted: value.IsAccepted.Bool,
		IsMuted:    value.IsMuted.Bool,
		Operation:  &pg.Type,
	}

}
