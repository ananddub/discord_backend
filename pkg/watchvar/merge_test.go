package watchvar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeWatchVars(t *testing.T) {
	wv1 := NewWatcher(10)
	wv2 := NewWatcher(20)
	defer wv1.Close()
	defer wv2.Close()

	merged := MergeWatchVars(wv1, wv2)
	defer merged.Close()
	sub := merged.Subscribe()

	received := make([]int, 0, 2)
	go func() {
		for i := 0; i < 2; i++ {
			v := <-sub
			received = append(received, v)
		}
	}()
	wv1.Send(5)
	wv2.Send(7)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 2, len(received))
	total := received[0] + received[1]
	assert.Equal(t, 12, total)
}

func TestMergeBroadCasts(t *testing.T) {
	bc1 := NewBroadCast[int]()
	bc2 := NewBroadCast[int]()
	defer bc1.Close()
	defer bc2.Close()

	merged := MergeBroadCasts(bc1, bc2)
	defer merged.Close()

	sub := merged.Subscribe()

	received := make([]int, 0, 2)
	go func() {
		for i := 0; i < 2; i++ {
			v := <-sub
			received = append(received, v)
		}
	}()
	bc1.Send(5)
	bc2.Send(7)
	time.Sleep(1 * time.Second)
	assert.Equal(t, 2, len(received))
	total := received[0] + received[1]
	assert.Equal(t, 12, total)
}
