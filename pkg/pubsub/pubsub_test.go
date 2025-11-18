package pubsub

import (
	"sync"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	ps := Get()
	received := make(chan string)

	// Subscribe in goroutine
	go func() {
		ch := ps.Subscribe("test")
		data := <-ch.Receive()
		received <- data.(string)
	}()

	// Small delay to ensure subscription
	time.Sleep(10 * time.Millisecond)

	// Publish
	ps.Publish("test", "hello")

	// Check
	select {
	case msg := <-received:
		if msg != "hello" {
			t.Errorf("Expected hello, got %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Should receive message")
	}
}

func TestExists(t *testing.T) {
	ps := Get()

	// Topic should not exist initially
	if ps.Exists("exists-test") {
		t.Error("Topic should not exist initially")
	}

	// Subscribe to create topic
	ch := ps.Subscribe("exists-test")

	// Topic should exist now
	if !ps.Exists("exists-test") {
		t.Error("Topic should exist after subscription")
	}

	// Close channel
	ch.Close()

	// Topic should not exist after cleanup
	if ps.Exists("exists-test") {
		t.Error("Topic should not exist after cleanup")
	}
}

func TestChannelClose(t *testing.T) {
	ps := Get()

	// Subscribe
	ch := ps.Subscribe("channel-close-test")

	// Close using channel method
	ch.Close()

	// Channel should be closed
	select {
	case _, ok := <-ch.Receive():
		if ok {
			t.Error("Channel should be closed after Close()")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Channel should be closed immediately")
	}
}

func TestCloseCh(t *testing.T) {
	ps := Get()
	topic := "close-ch-test"
	closeCh := make(chan int)
	ch := ps.Subscribe(topic)
	go func() {
		defer func() {
			t.Logf("closed channel")
			closeCh <- 1
		}()
		for data := range ch.Receive() {
			t.Logf("Received: %v", data)
		}
	}()
	ps.Publish(topic, "hello")
	ps.Publish(topic, "world")

	time.Sleep(100 * time.Millisecond)
	ch.Close()
	<-closeCh
}

func TestClose(t *testing.T) {
	ps := Get()

	// Subscribe
	ch := ps.Subscribe("close-test")

	// Close topic
	ch.Close()

	// Channel should be closed
	select {
	case _, ok := <-ch.Receive():
		if ok {
			t.Error("Channel should be closed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Channel should be closed immediately")
	}
}

func TestBulkSubscribe(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	topics := []string{"bulk1", "bulk2", "bulk3"}
	bulkCh := ps.BulkSubscribe(topics)
	defer bulkCh.Close()

	// Publish to different topics
	ps.Publish("bulk1", "message1")
	ps.Publish("bulk2", "message2")
	ps.Publish("bulk3", "message3")

	received := make(map[string]bool)

	for i := 0; i < 3; i++ {
		select {
		case msg := <-bulkCh.Receive():
			received[msg.(string)] = true
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for messages")
		}
	}

	if !received["message1"] || !received["message2"] || !received["message3"] {
		t.Error("Not all messages received from bulk subscription")
	}
}

func TestBulkSubscribeClose(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	topics := []string{"cleanup1", "cleanup2", "cleanup3"}
	bulkCh := ps.BulkSubscribe(topics)

	// Verify all topics exist
	for _, topic := range topics {
		if !ps.Exists(topic) {
			t.Errorf("Topic %s should exist", topic)
		}
	}
	v := ps.ListTopic()
	t.Logf("Current topics: %v", v)
	// Close bulk subscription
	bulkCh.Close()

	// Give time for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify all topics are cleaned up
	for _, topic := range topics {
		if ps.Exists(topic) {
			t.Errorf("Topic %s should not exist after bulk close", topic)
		}
	}
}

func TestConcurrentPublish(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	ch := ps.Subscribe("concurrent")
	defer ch.Close()

	var wg sync.WaitGroup
	messageCount := 100

	// Publish concurrently
	for i := 0; i < messageCount; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			ps.Publish("concurrent", num)
		}(i)
	}

	wg.Wait()

	// Receive messages
	received := 0
	timeout := time.After(2 * time.Second)

loop:
	for {
		select {
		case <-ch.Receive():
			received++
			if received == messageCount {
				break loop
			}
		case <-timeout:
			break loop
		}
	}

	if received != messageCount {
		t.Errorf("Expected %d messages, received %d", messageCount, received)
	}
}

func TestCloseTopic(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	ch1 := ps.Subscribe("close-topic")
	ch2 := ps.Subscribe("close-topic")

	ps.Close("close-topic")

	// Channels should be closed
	_, ok1 := <-ch1.Receive()
	_, ok2 := <-ch2.Receive()

	if ok1 || ok2 {
		t.Error("Channels should be closed")
	}

	if ps.Exists("close-topic") {
		t.Error("Topic should not exist after Close")
	}
}

func TestCloseAll(t *testing.T) {
	ps := Get()

	ch1 := ps.Subscribe("topic1")
	ch2 := ps.Subscribe("topic2")
	ch3 := ps.Subscribe("topic3")

	ps.CloseAll()

	// All channels should be closed
	_, ok1 := <-ch1.Receive()
	_, ok2 := <-ch2.Receive()
	_, ok3 := <-ch3.Receive()

	if ok1 || ok2 || ok3 {
		t.Error("All channels should be closed")
	}

	topics := ps.ListTopic()
	if len(topics) != 0 {
		t.Error("All topics should be removed")
	}
}

func TestFullChannelBuffer(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	ch := ps.Subscribe("buffer-test")
	defer ch.Close()

	// Fill the buffer (100 capacity)
	for i := 0; i < 150; i++ {
		ps.Publish("buffer-test", i)
	}

	// Should not block or panic
	received := 0
	timeout := time.After(1 * time.Second)

loop:
	for {
		select {
		case <-ch.Receive():
			received++
		case <-timeout:
			break loop
		}
	}

	// Should receive at most 100 messages (buffer size)
	if received > 100 {
		t.Errorf("Received more than buffer capacity: %d", received)
	}
}

func TestListTopic(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	ch1 := ps.Subscribe("list1")
	ch2 := ps.Subscribe("list2")
	defer ch1.Close()
	defer ch2.Close()

	topics := ps.ListTopic()

	if len(topics) != 2 {
		t.Errorf("Expected 2 topics, got %d", len(topics))
	}

	if _, ok := topics["list1"]; !ok {
		t.Error("list1 should be in topics")
	}

	if _, ok := topics["list2"]; !ok {
		t.Error("list2 should be in topics")
	}
}

func TestDoubleClose(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	ch := ps.Subscribe("double-close")

	// Close twice should not panic
	ch.Close()
	ch.Close()

	// Should not panic
}

func TestBulkSubscribeMixedPublish(t *testing.T) {
	ps := Get()
	defer ps.CloseAll()

	bulkCh := ps.BulkSubscribe([]string{"mix1", "mix2"})
	defer bulkCh.Close()

	singleCh := ps.Subscribe("mix1")
	defer singleCh.Close()

	ps.Publish("mix1", "msg1")
	ps.Publish("mix2", "msg2")

	// Bulk channel should receive both
	received := make(map[string]bool)
	for i := 0; i < 2; i++ {
		select {
		case msg := <-bulkCh.Receive():
			received[msg.(string)] = true
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout on bulk channel")
		}
	}

	if !received["msg1"] || !received["msg2"] {
		t.Error("Bulk channel should receive both messages")
	}

	// Single channel should receive only mix1
	select {
	case msg := <-singleCh.Receive():
		if msg != "msg1" {
			t.Errorf("Single channel should receive msg1, got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout on single channel")
	}
}
