package watchvar

import (
	"sync/atomic"
)

func MergeWatchVars[T any](watchvars ...*WatchVar[T]) *BroadCast[T] {
	broadcast := NewBroadCast[T]()
	for _, wv := range watchvars {
		wv.customSubscribe(broadcast.in)
	}
	return broadcast
}
func MergeBroadCasts[T any](broadcasts ...*BroadCast[T]) *BroadCast[T] {
	broadcast := NewBroadCast[T]()
	var counter int32
	atomic.StoreInt32(&counter, int32(len(broadcasts)))
	for _, bc := range broadcasts {
		bc.customSubscribe(broadcast.in)
	}
	return broadcast
}
