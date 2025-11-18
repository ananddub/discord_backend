package service

import (
	"context"
	"discord/gen/proto/schema"
	"discord/pkg/pubsub"
	"strconv"
)

func UserTopic(userId int32) string {
	return "user:" + strconv.Itoa(int(userId))
}

func UserFriendTopic(userId int32) string {
	return "friend_pr:" + strconv.Itoa(int(userId))
}

func (s *UserService) publishUser(ctx context.Context, user *schema.User) {
	topic := UserTopic(user.Id)
	go s.pubsub.Publish(topic, user)
	s.publishFriendUpdate(ctx, user)
}

func (s *UserService) publishFriendUpdate(ctx context.Context, user *schema.User) {
	connectedFriends, err := s.GetConnectedFriends(ctx, user.Id)
	if err != nil {
		return
	}
	go func() {
		for _, friend := range connectedFriends {
			friendTopic := UserFriendTopic(friend.ID)
			pubsub.Get().Publish(friendTopic, user)
		}
	}()
}

func StreamUser(id int32) *pubsub.Channel {
	ps := pubsub.Get()
	ch := ps.Subscribe(UserTopic(id))
	return ch
}

func StreamUserFriendUpdates(id int32) *pubsub.Channel {
	ps := pubsub.Get()
	ch := ps.Subscribe(UserFriendTopic(id))
	return ch
}
