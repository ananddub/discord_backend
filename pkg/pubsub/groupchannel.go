package pubsub

import "sync"

type mydata struct {
	value any
}
type ArrPub[T any] struct {
	Data    mydata
	ch      []*RChannel[T]
	publish *chan T
	done    chan bool
	lock    sync.RWMutex
}

func NewArrPub[T any]() *ArrPub[T] {

	ch := make(chan T, 100)
	v := &ArrPub[T]{
		Data:    mydata{},
		ch:      make([]*RChannel[T], 0, 100),
		lock:    sync.RWMutex{},
		done:    make(chan bool),
		publish: &ch,
	}

	go func() {
		for {
			select {
			case <-v.done:
				return
			case data := <-ch:
				v.Publish(data)
			}
		}
	}()
	return v
}

func (arr *ArrPub[T]) Publish(data T) {
	arr.lock.RLock()
	defer arr.lock.RUnlock()
	for _, ch := range arr.ch {
		ch.spublish(data)
	}
}

func (arr *ArrPub[T]) Subscribe() *RChannel[T] {
	ch := NewRChannel[T](arr.publish, &arr.Data)
	arr.lock.Lock()
	arr.ch = append(arr.ch, ch)
	ch.idx = len(arr.ch) - 1
	arr.lock.Unlock()
	go func() {
		_, ok := <-ch.done
		if ok {
			arr.remove(ch)
			if len(arr.ch) == 0 && !ok {
				arr.done <- true
			}
		}

	}()

	return ch
}

func (arr *ArrPub[T]) remove(ch *RChannel[T]) {
	arr.lock.Lock()
	defer arr.lock.Unlock()

	idx := ch.idx
	if idx < 0 {
		return
	}

	last := len(arr.ch) - 1

	if idx != last {
		arr.ch[idx] = arr.ch[last]
		arr.ch[idx].idx = idx
	}
	arr.ch = arr.ch[:last]
	ch.idx = -1
}

func (arr *ArrPub[T]) Close() {
	for _, ch := range arr.ch {
		arr.remove(ch)
		ch.Close()
	}
	close(arr.done)
}

func (arr *ArrPub[T]) Size() int {
	return len(arr.ch)
}

func (arr *ArrPub[T]) Wait() {
	<-arr.done
}

func (arr *ArrPub[T]) SetData(data any) {
	arr.Data.value = data
}

func (arr *ArrPub[T]) GetData() any {
	return arr.Data.value
}
