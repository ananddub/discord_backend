package pubsub

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNPubSub(t *testing.T) {
	arr := NewArrPub[string]()
	ch := arr.Subscribe()
	ch1 := arr.Subscribe()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case data, ok := <-ch.Subscribe():
				if !ok {
					return
				}
				fmt.Println("ch1 :", data)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case data, ok := <-ch1.Subscribe():
				if !ok {
					return
				}
				fmt.Println("ch2 :", data)
			}
		}
	}()
	arr.Publish("test")
	ch.Publish("hello")
	ch1.Publish("world")
	time.Sleep(1 * time.Second)
	assert.Equal(t, arr.Size(), 2)
	arr.Close()
	wg.Wait()
	assert.Equal(t, arr.Size(), 0)
}
func TestNPubSubClose(t *testing.T) {
	arr := NewArrPub[string]()
	ch := arr.Subscribe()
	ch1 := arr.Subscribe()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case data, ok := <-ch.Subscribe():
				if !ok {
					return
				}
				fmt.Println("ch1 :", data)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case data, ok := <-ch1.Subscribe():
				if !ok {
					return
				}
				fmt.Println("ch2 :", data)
			}
		}
	}()
	arr.Publish("test")
	ch.Publish("hello")
	ch1.Publish("world")
	time.Sleep(1 * time.Second)
	assert.Equal(t, arr.Size(), 2)
	ch.Close()

	time.Sleep(1 * time.Second)
	assert.Equal(t, arr.Size(), 1)
	ch1.Close()

	time.Sleep(1 * time.Second)
	assert.Equal(t, arr.Size(), 0)
	wg.Wait()
}

func TestNDataCheck(t *testing.T) {
	arr := NewArrPub[string]()
	v := "hello"
	arr.SetData(v)
	ch := arr.Subscribe()
	fmt.Println(ch.GetData())
	assert.Equal(t, ch.GetData(), "hello")
	v = "world"
	ch.SetData(v)
	fmt.Println(ch.GetData())
	assert.Equal(t, ch.GetData(), "world")
	assert.Equal(t, arr.GetData(), "world")
}
