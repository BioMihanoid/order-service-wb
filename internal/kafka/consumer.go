package kafka

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Consumer struct {
	client *kgo.Client
}

func NewConsumer(brokers []string, group, topic string) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{client: client}, nil
}

func (c *Consumer) Run(ctx context.Context, handler func(msg *kgo.Record) error) {
	for {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("kafka fetch error: %v", err)
			}
			continue
		}
		fetches.EachRecord(func(record *kgo.Record) {
			if err := handler(record); err == nil {
				if commitErr := c.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit error: %v", commitErr)
				}
			} else {
				log.Printf("handler error: %v", err)
			}
		})
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
