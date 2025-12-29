package watchvar

import (
	"sync/atomic"
	"unsafe"
)

type WatchVar[T any] struct {
	data      unsafe.Pointer
	broadcast *BroadCast[T]
}

func NewWatcher[T any](data T) *WatchVar[T] {
	v := &WatchVar[T]{
		broadcast: NewBroadCast[T](),
	}
	atomic.StorePointer(&v.data, unsafe.Pointer(&data))
	return v
}

func (wv *WatchVar[T]) Get() *T {
	return (*T)(atomic.LoadPointer(&wv.data))
}

func (wv *WatchVar[T]) Send(data T) {
	atomic.StorePointer(&wv.data, unsafe.Pointer(&data))
	wv.broadcast.Send(data)
}

func (wv *WatchVar[T]) SendModify(modify func(data *T)) {
	for {
		oldPtr := atomic.LoadPointer(&wv.data)
		newData := *(*T)(oldPtr)
		modify(&newData)
		if atomic.CompareAndSwapPointer(&wv.data, oldPtr, unsafe.Pointer(&newData)) {
			wv.broadcast.Send(newData)
			break
		}
	}
}

func (wv *WatchVar[T]) SendModifyIf(modify func(data *T) bool) {
	for {
		oldPtr := atomic.LoadPointer(&wv.data)
		oldData := *(*T)(oldPtr)
		newData := oldData
		if modify(&newData) {
			if atomic.CompareAndSwapPointer(&wv.data, oldPtr, unsafe.Pointer(&newData)) {
				wv.broadcast.Send(newData)
				break
			}
		} else {
			break
		}
	}
}

func (wv *WatchVar[T]) customSubscribe(ch chan T) <-chan T {
	wv.broadcast.customSubscribe(ch)
	return ch
}

func (wv *WatchVar[T]) Subscribe() <-chan T {
	ch := wv.broadcast.subscribe()
	return ch
}

func (wv *WatchVar[T]) Close() {
	wv.broadcast.Close()
}
