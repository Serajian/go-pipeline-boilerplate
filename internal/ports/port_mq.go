package ports

import (
	"context"
)

// MessageQueueProducer defines the contract for a message queue producer.
// A producer is responsible for establishing a connection to the broker,
// sending messages to a specified topic, and closing the connection.
//
// Implementations should handle retries, batching, and any broker-specific
// configuration internally.
type MessageQueueProducer interface {
	Connect(ctx context.Context) error
	Close() error
	Produce(ctx context.Context, topic string, msg interface{}) error
}

// MessageQueueConsumer defines the contract for a message queue consumer.
// A consumer is responsible for connecting to the broker, continuously
// consuming messages from one or more topics, and closing the connection.
//
// Implementations should respect the provided context to allow
// cancellation and graceful shutdown.
type MessageQueueConsumer interface {
	Connect(ctx context.Context) error
	Close() error
	Consume(ctx context.Context) error
}
