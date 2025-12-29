package watchvar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBroadCastBasic(t *testing.T) {
	bc := NewBroadCast[int]()
	defer bc.Close()

	ch := bc.Subscribe()

	bc.Send(42)
	v := <-ch
	assert.Equal(t, 42, v)
}

func TestBroadCastMultipleSubscribers(t *testing.T) {
	bc := NewBroadCast[int]()
	defer bc.Close()

	ch1 := bc.Subscribe()
	ch2 := bc.Subscribe()

	bc.Send(100)

	v1 := <-ch1
	v2 := <-ch2

	assert.Equal(t, 100, v1)
	assert.Equal(t, 100, v2)
}

func TestBroadCastNonBlocking(t *testing.T) {
	bc := NewBroadCast[int]()
	defer bc.Close()

	bc.Send(1)
	bc.Send(2)
	bc.Send(3)

	ch := bc.Subscribe()
	bc.Send(4)

	v := <-ch
	assert.Equal(t, 1, v)
}
