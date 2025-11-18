package pubsub

import "sync"

type Channel struct {
	ch           chan interface{}
	topic        string
	pubsub       *PubSub
	isBulk       bool
	bulkTopics   []string
	bulkChannels []chan interface{}
	closeOnce    sync.Once
}

func (c *Channel) Receive() <-chan interface{} {
	return c.ch
}

func (c *Channel) Close() {
	c.closeOnce.Do(func() {
		if c.isBulk {
			c.pubsub.cleanupBulkSubscription(c)
		} else {
			c.pubsub.Unsubscribe(c.topic, c.ch)
		}
	})
}

type PubSub struct {
	topics map[string][]chan interface{}
	mu     sync.RWMutex
}

var (
	instance *PubSub
	once     sync.Once
)

func Get() *PubSub {
	once.Do(func() {
		instance = &PubSub{
			topics: make(map[string][]chan interface{}),
		}
	})
	return instance
}

func (ps *PubSub) ListTopic() map[string][]chan interface{} {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	copy := make(map[string][]chan interface{})
	for topic, channels := range ps.topics {
		copy[topic] = append([]chan interface{}{}, channels...)
	}
	return copy
}

func (ps *PubSub) Subscribe(topic string) *Channel {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 100)
	ps.topics[topic] = append(ps.topics[topic], ch)

	return &Channel{
		ch:     ch,
		topic:  topic,
		pubsub: ps,
		isBulk: false,
	}
}

func (ps *PubSub) Exists(topic string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	_, exists := ps.topics[topic]
	return exists
}

func (ps *PubSub) Unsubscribe(topic string, ch chan interface{}) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	channels := ps.topics[topic]
	for i, channel := range channels {
		if channel == ch {
			close(channel)
			ps.topics[topic] = append(channels[:i], channels[i+1:]...)

			if len(ps.topics[topic]) == 0 {
				delete(ps.topics, topic)
			}
			break
		}
	}
}

func (ps *PubSub) Publish(topic string, data interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, ch := range ps.topics[topic] {
		select {
		case ch <- data:
		default:
			// Channel full, skip
		}
	}
}

func (ps *PubSub) Close(topic string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range ps.topics[topic] {
		close(ch)
	}
	delete(ps.topics, topic)
}

func (ps *PubSub) CloseAll() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for topic, channels := range ps.topics {
		for _, ch := range channels {
			close(ch)
		}
		delete(ps.topics, topic)
	}
}

func (ps *PubSub) BulkSubscribe(topics []string) *Channel {
	ps.mu.Lock()

	mergedCh := make(chan interface{}, 100)
	internalChannels := make([]chan interface{}, len(topics))

	// Create internal channels and subscribe to each topic
	for i, topic := range topics {
		ch := make(chan interface{}, 100)
		ps.topics[topic] = append(ps.topics[topic], ch)
		internalChannels[i] = ch
	}

	ps.mu.Unlock()

	// Forward messages from all internal channels to merged channel
	var wg sync.WaitGroup
	for _, ch := range internalChannels {
		wg.Add(1)
		go func(internalCh chan interface{}) {
			defer wg.Done()
			for msg := range internalCh {
				mergedCh <- msg
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(mergedCh)
	}()

	return &Channel{
		ch:           mergedCh,
		topic:        "BULK",
		pubsub:       ps,
		isBulk:       true,
		bulkTopics:   topics,
		bulkChannels: internalChannels,
	}
}

func (ps *PubSub) cleanupBulkSubscription(c *Channel) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, topic := range c.bulkTopics {
		internalCh := c.bulkChannels[i]
		channels := ps.topics[topic]

		for j, ch := range channels {
			if ch == internalCh {
				close(ch)
				ps.topics[topic] = append(channels[:j], channels[j+1:]...)

				if len(ps.topics[topic]) == 0 {
					delete(ps.topics, topic)
				}
				break
			}
		}
	}
}
