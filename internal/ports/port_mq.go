package ports

import (
	"context"
)

// MessageQueueProducer defines the interface for produce message queue operations.
type MessageQueueProducer interface {
	Connect(ctx context.Context) error
	Close() error
	Produce(ctx context.Context, topic string, msg interface{}) error
}

// MessageQueueConsumer defines the interface for consume message queue operations.
type MessageQueueConsumer interface {
	Connect(ctx context.Context) error
	Close() error
	Consume(ctx context.Context) error
}
