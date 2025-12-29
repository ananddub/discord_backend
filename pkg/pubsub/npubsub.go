package pubsub

type NPubSub struct {
	User map[string]*UPubSub
}

func NewNPubSub() *NPubSub {
	return &NPubSub{
		User: make(map[string]*UPubSub),
	}
}

func (NPubSub *NPubSub) Create(key string) *UPubSub {
	if NPubSub.User[key] == nil {
		NPubSub.User[key] = NewUPubSub()
	}
	return NPubSub.User[key]
}

func (NPubSub *NPubSub) Get(key string) *UPubSub {
	return NPubSub.User[key]
}

func (NPubSub *NPubSub) Has(key string) bool {
	return NPubSub.User[key] != nil
}

func (NPubSub *NPubSub) Size() int {
	return len(NPubSub.User)
}

func (NPubSub *NPubSub) Remove(key string) {
	if NPubSub.User[key] == nil {
		return
	}
	NPubSub.User[key].Close()
	NPubSub.User[key] = nil
	delete(NPubSub.User, key)
}

func (NPubSub *NPubSub) Close() {
	for _, upub := range NPubSub.User {
		upub.Close()
	}
}
