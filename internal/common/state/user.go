package state

import (
	"discord/gen/proto/schema"
	"discord/pkg/watchvar"
)

type UserState struct {
	user        *watchvar.WatchVar[*schema.User]
	friendState *map[string]*watchvar.WatchVar[*schema.Friend]
	hashMapRef  *HashMap
}

func NewUserState(user *schema.User, hashmapref *HashMap) *UserState {
	return &UserState{
		user:        watchvar.NewWatcher(user),
		friendState: &map[string]*watchvar.WatchVar[*schema.Friend]{},
		hashMapRef:  hashmapref,
	}
}
