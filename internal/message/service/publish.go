package service

import (
	"discord/gen/proto/schema"
	"discord/pkg/pubsub"
	"strconv"
)

func Topic(userId int32) string {
	return "friend:" + strconv.Itoa(int(userId))
}

func (s *MessageService) publishToFriend(id int32, friend *schema.Friend) {
	topic := Topic(id)
	s.pubsub.Publish(topic, friend)
}

func Stream(id int32) *pubsub.Channel {
	ps := pubsub.Get()
	ch := ps.Subscribe(Topic(id))
	return ch
}
