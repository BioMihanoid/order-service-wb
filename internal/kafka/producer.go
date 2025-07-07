package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.RecordDeliveryTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	return &Producer{
		client: client,
		topic:  topic,
	}, nil

}

func (p *Producer) Send(ctx context.Context, key string, value []byte) error {
	return p.client.ProduceSync(ctx, &kgo.Record{
		Topic: p.topic,
		Key:   []byte(key),
		Value: value,
	}).FirstErr()
}

func (p *Producer) Close() {
	p.client.Close()
}
