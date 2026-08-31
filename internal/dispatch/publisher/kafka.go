package publisher

/*
================================================================================
FILE: internal/dispatch/publisher/kafka.go
================================================================================

PURPOSE:
Kafka Producer wrapper using `segmentio/kafka-go` or `confluent-kafka-go`.
Publishes JSON event payloads (`TripCreatedEvent`, `TripAssignedEvent`, `TripCompletedEvent`)
to Kafka topics.

LEARNING GO CONCEPTS:
- Producing messages to Kafka (`kafka.Writer`).
- JSON serialization of event payloads (`json.Marshal`).
- Asynchronous vs Synchronous message delivery.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `KafkaPublisher` struct wrapping `*kafka.Writer`.
2. `PublishEvent(ctx context.Context, topic string, key string, payload any) error`:
   - Marshal payload to JSON bytes.
   - Write message with `Writer.WriteMessages(...)`.
================================================================================
*/

import (
	"context"
	// TODO: Uncomment when adding kafka-go:
	// "github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	// writer *kafka.Writer
}

func NewKafkaPublisher(broker, topic string) *KafkaPublisher {
	return &KafkaPublisher{}
}

func (p *KafkaPublisher) PublishEvent(ctx context.Context, key string, payload any) error {
	// TODO: Marshal payload to JSON and send to Kafka topic
	return nil
}
