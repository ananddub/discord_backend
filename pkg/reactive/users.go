package reactive

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
	"discord/pkg/mypg"
	"discord/pkg/pubsub"
	"strconv"
	"strings"
)

type UserReactive struct{}

func (s UserReactive) ReactUser(pg mypg.QueryData) {
	s.publish(s.convertToUser(pg))
}

func (s UserReactive) publish(user *schema.User) {
	pub := pubsub.Get()
	topic := "user:" + strconv.Itoa(int(user.Id))
	pub.Publish(topic, &user)
}

func (s UserReactive) convertToUser(pg mypg.QueryData) *schema.User {
	var value repo.User
	pg.Data.Scan(
		&value.ID,
		&value.Username,
		&value.Email,
		&value.FullName,
		&value.Status,
		&value.CustomStatus,
		&value.ProfilePic,
		&value.BackgroundColor,
		&value.ColorCode,
		&value.Bio,
		&value.IsBot,
		&value.IsVerified,
		&value.CreatedAt,
	)
	return &schema.User{
		Id:              value.ID,
		Username:        value.Username,
		Email:           value.Email,
		FullName:        value.FullName.String,
		Status:          value.Status,
		CustomStatus:    value.CustomStatus.String,
		ProfilePic:      value.ProfilePic.String,
		BackgroundColor: value.BackgroundColor.String,
		ColorCode:       value.ColorCode.String,
		Bio:             value.Bio.String,
		IsBot:           value.IsBot.Bool,
		IsVerified:      value.IsVerified.Bool,
		CreatedAt:       value.CreatedAt.Time.Unix(),
		Operation:       strings.ToLower(pg.Type),
	}
}
