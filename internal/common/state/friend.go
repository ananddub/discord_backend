package state

import (
	"discord/gen/proto/schema"
	"discord/pkg/watchvar"
)

type FriendState struct {
	friends *watchvar.WatchVar[*schema.Friend]
}

func NewFriendState(friend *schema.Friend) *FriendState {
	return &FriendState{
		friends: watchvar.NewWatcher(friend),
	}
}

func (fs *FriendState) Subscribe() <-chan *schema.Friend {
	return fs.friends.Subscribe()
}
