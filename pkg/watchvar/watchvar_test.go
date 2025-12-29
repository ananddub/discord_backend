package watchvar

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWatchVarBasic(t *testing.T) {
	wv := NewWatcher(10)
	defer wv.Close()

	ch := wv.Subscribe()

	v := <-ch
	assert.Equal(t, v, 10)

	// Send new value
	wv.Send(20)
	v = <-ch
	assert.Equal(t, 20, v)
}

func TestWatchVarModify(t *testing.T) {
	wv := NewWatcher(5)
	defer wv.Close()
	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()
		wg.Add(1)
		v := <-wv.Subscribe()
		assert.Equal(t, 15, v)
	}()
	go func() {
		defer wg.Done()
		wg.Add(1)
		v := <-wv.Subscribe()
		assert.Equal(t, 15, v)
	}()
	time.Sleep(time.Second)
	wv.Send(20)

	wv.SendModify(func(data *int) {
		*data += 10
	})
	wg.Wait()
}

func TestSendModifyIf(t *testing.T) {
	wv := NewWatcher(5)
	defer wv.Close()
	wg := sync.WaitGroup{}
	go func() {
		defer wg.Done()
		wg.Add(1)
		v := <-wv.Subscribe()
		assert.Equal(t, 15, v)
	}()
	go func() {
		defer wg.Done()
		wg.Add(1)
		v := <-wv.Subscribe()
		assert.Equal(t, 15, v)
	}()
	wv.SendModifyIf(func(data *int) bool {
		*data += 10
		return true
	})
	time.Sleep(time.Second)
	wv.SendModify(func(data *int) {
		*data += 10
	})
	value := wv.Get()
	fmt.Println(value)
	wg.Wait()
}

func TestWatchVarGet(t *testing.T) {
	wv := NewWatcher(42)
	defer wv.Close()

	val := wv.Get()
	if *val != 42 {
		t.Fatalf("Expected 42, got %d", *val)
	}
}
