package watchvar

import (
	"sync/atomic"
	"unsafe"
)

type BroadCast[T any] struct {
	in          chan T
	subscribers unsafe.Pointer
	started     int64
}

func NewBroadCast[T any]() *BroadCast[T] {
	b := &BroadCast[T]{
		in: make(chan T, 1),
	}
	atomic.StorePointer(&b.subscribers, unsafe.Pointer(&[]chan T{}))
	b.start()
	return b
}

func (b *BroadCast[T]) start() {
	if atomic.CompareAndSwapInt64(&b.started, 0, 1) {
		go func() {
			for v := range b.in {
				subsPtr := atomic.LoadPointer(&b.subscribers)
				subs := *(*[]chan T)(subsPtr)
				writeIdx := 0
				for i := 0; i < len(subs); i++ {
					func() {
						defer func() {
							recover()
						}()
						select {
						case subs[i] <- v:
							if writeIdx != i {
								subs[writeIdx] = subs[i]
							}
							writeIdx++
						default:
						}
					}()
				}
				if writeIdx != len(subs) {
					newSubs := subs[:writeIdx]
					atomic.StorePointer(&b.subscribers, unsafe.Pointer(&newSubs))
				}
			}
		}()
	}
}

func (b *BroadCast[T]) Send(data T) {
	select {
	case b.in <- data:
	default:
	}
}

func (b *BroadCast[T]) Subscribe() <-chan T {
	ch := make(chan T, 1)
	for {
		oldPtr := atomic.LoadPointer(&b.subscribers)
		oldSubs := *(*[]chan T)(oldPtr)
		newSubs := make([]chan T, len(oldSubs)+1)
		copy(newSubs, oldSubs)
		newSubs[len(oldSubs)] = ch
		if atomic.CompareAndSwapPointer(&b.subscribers, oldPtr, unsafe.Pointer(&newSubs)) {
			break
		}
	}
	return ch

}
func (b *BroadCast[T]) customSubscribe(ch chan T) <-chan T {
	for {
		oldPtr := atomic.LoadPointer(&b.subscribers)
		oldSubs := *(*[]chan T)(oldPtr)
		newSubs := make([]chan T, len(oldSubs)+1)
		copy(newSubs, oldSubs)
		newSubs[len(oldSubs)] = ch
		if atomic.CompareAndSwapPointer(&b.subscribers, oldPtr, unsafe.Pointer(&newSubs)) {
			break
		}
	}
	return ch
}
func (b *BroadCast[T]) subscribe() chan T {
	ch := make(chan T, 1)
	for {
		oldPtr := atomic.LoadPointer(&b.subscribers)
		oldSubs := *(*[]chan T)(oldPtr)
		newSubs := make([]chan T, len(oldSubs)+1)
		copy(newSubs, oldSubs)
		newSubs[len(oldSubs)] = ch
		if atomic.CompareAndSwapPointer(&b.subscribers, oldPtr, unsafe.Pointer(&newSubs)) {
			break
		}
	}
	return ch
}

func (b *BroadCast[T]) Close() {
	defer func() { recover() }()
	close(b.in)
	subsPtr := atomic.LoadPointer(&b.subscribers)
	subs := *(*[]chan T)(subsPtr)
	for i := range subs {
		func() {
			defer func() { recover() }()
			close(subs[i])
		}()
	}
}
