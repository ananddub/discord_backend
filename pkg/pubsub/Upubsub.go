package pubsub

type UPubSub struct {
	value map[string]*ArrPub[any]
	done  chan bool
}

func NewUPubSub() *UPubSub {
	return &UPubSub{
		value: make(map[string]*ArrPub[any]),
		done:  make(chan bool),
	}
}

func (upub *UPubSub) Publish(key string, data any) {
	if upub.value[key] == nil {
		upub.value[key] = NewArrPub[any]()
	}
	upub.value[key].Publish(data)
}

func (upub *UPubSub) Subscribe(key string) *RChannel[any] {
	if upub.value[key] == nil {
		upub.value[key] = NewArrPub[any]()
	}
	return upub.value[key].Subscribe()
}

func (upub *UPubSub) Wait() {
	for _, ch := range upub.value {
		ch.Wait()
	}
	upub.done <- true
}

func (upub *UPubSub) Close() {
	for _, ch := range upub.value {
		ch.Close()
	}
	close(upub.done)
}
