package message_queue

import (
	"context"
	"fmt"
	"time"

	"go-pipeline/config"
	"go-pipeline/internal/ports"
	"go-pipeline/pkg/apperror"
	"go-pipeline/pkg/logger"

	"github.com/IBM/sarama"
)

// KafkaConsumerConfig holds the configuration for a Kafka consumer instance.
type KafkaConsumerConfig struct {
	Name    string   // Logical name of the consumer
	Address string   // Broker address (hostname or IP)
	Port    int      // Broker port
	GroupID string   // Kafka consumer group ID
	Topics  []string // List of topics to subscribe to
}

// KafkaConsumerAdapter implements the ports.MessageQueueConsumer interface
// using Sarama's ConsumerGroup.
type KafkaConsumerAdapter struct {
	Config   *KafkaConsumerConfig // Consumer configuration
	Consumer sarama.ConsumerGroup // Underlying Sarama consumer group
	Handler  sarama.ConsumerGroupHandler
}

// Connect initializes and connects to the Kafka consumer group with the given config.
func (c *KafkaConsumerAdapter) Connect(ctx context.Context) error {
	brokers := []string{fmt.Sprintf("%s:%d", c.Config.Address, c.Config.Port)}

	cfg := sarama.NewConfig()

	// 1. Kafka protocol version (must match the broker version)
	cfg.Version = sarama.V3_9_0_0

	// 2. Consumer behavior options
	cfg.Consumer.Return.Errors = true                            // Forward errors to the Errors() channel
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest           // Start from oldest if no committed offset is found
	cfg.Consumer.Group.Rebalance.Retry.Max = 10                  // Retry rebalancing attempts up to 10 times
	cfg.Consumer.Group.Rebalance.Retry.Backoff = 2 * time.Second // Backoff duration between retries

	// 3. Offset commit settings
	cfg.Consumer.Offsets.AutoCommit.Enable = true              // Enable automatic offset commits
	cfg.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // Commit interval
	cfg.Consumer.Offsets.Retry.Max = 3                         // Retry commit if it fails

	// Create a new consumer group instance
	consumer, err := sarama.NewConsumerGroup(brokers, c.Config.GroupID, cfg)
	if err != nil {
		return fmt.Errorf("%w: failed to consume message: %w", apperror.ErrUnavailable, err)
	}

	c.Consumer = consumer

	// Log successful connection
	logger.GetLogger().Info(&logger.Log{
		Event:   "connect kafka consumer",
		TraceID: config.GetTraceID(ctx),
	})

	return nil
}

// Consume starts consuming messages for the configured topics
// using the provided ConsumerGroupHandler.
func (c *KafkaConsumerAdapter) Consume(
	ctx context.Context,
) error {
	if c.Consumer == nil {
		return fmt.Errorf("%w: consumer group not connected", apperror.ErrUnavailable)
	}
	return c.Consumer.Consume(ctx, c.Config.Topics, c.Handler)
}

// Close shuts down the consumer group connection gracefully.
func (c *KafkaConsumerAdapter) Close() error {
	if c.Consumer != nil {
		return c.Consumer.Close()
	}
	return nil
}

// Ensure KafkaConsumerAdapter implements the MessageQueueConsumer interface.
var _ ports.MessageQueueConsumer = (*KafkaConsumerAdapter)(nil)
