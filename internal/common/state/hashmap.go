package state

type HashMap struct {
	user   map[string]*UserState
	friend map[string]*FriendState
}

func NewHashMap() *HashMap {
	return &HashMap{
		user:   make(map[string]*UserState),
		friend: make(map[string]*FriendState),
	}
}
