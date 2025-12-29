package redpanda

import (
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Config struct {
	Brokers        []string
	GroupID        string
	Topics         []string
	AutoCreate     bool
	OffsetEarliest bool
	OffsetLatest   bool
}

func New(cfg Config) *kgo.Client {
	opts := []kgo.Opt{kgo.SeedBrokers(cfg.Brokers...)}
	if cfg.AutoCreate {
		opts = append(opts, kgo.AllowAutoTopicCreation())
	}

	if cfg.GroupID != "" {
		opts = append(opts, kgo.ConsumerGroup(cfg.GroupID))
	}

	if len(cfg.Topics) > 0 {
		opts = append(opts, kgo.ConsumeTopics(cfg.Topics...))
	}

	if cfg.OffsetEarliest {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	}
	if cfg.OffsetLatest {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	}

	opts = append(opts,
		kgo.FetchMaxBytes(5*1024*1024),
		kgo.ClientID("redpanda-go-client"),
	)

	opts = append(opts,
		kgo.RecordRetries(5),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
	)

	// --- AUTO COMMIT ---
	opts = append(opts,
		kgo.AutoCommitInterval(2*time.Second),
	)

	client, err := kgo.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	return client
}
