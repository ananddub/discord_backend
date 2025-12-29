package pubsub

type RChannel[T any] struct {
	done    chan bool
	Data    *mydata
	publish *chan T
	idx     int
	ch      chan T
}

func NewRChannel[T any](ch *chan T, data *mydata) *RChannel[T] {
	return &RChannel[T]{
		ch:      make(chan T, 100),
		done:    make(chan bool),
		Data:    data,
		publish: ch,
	}
}

func (ch *RChannel[T]) Publish(data T) {
	*ch.publish <- data
}

func (ch *RChannel[T]) SetData(data any) {
	ch.Data.value = data
}

func (ch *RChannel[T]) GetData() any {
	return ch.Data.value
}

func (ch *RChannel[T]) spublish(data T) {
	ch.ch <- data
}

func (ch *RChannel[T]) Subscribe() <-chan T {
	return ch.ch
}

func (ch *RChannel[T]) Close() {
	close(ch.ch)
	ch.done <- true
}
