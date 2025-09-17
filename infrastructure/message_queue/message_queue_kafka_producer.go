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

// KafkaProducerConfig holds configuration settings for the Kafka producer.
type KafkaProducerConfig struct {
	Name    string // Logical name of the producer
	Port    int    // Kafka broker port
	Address string // Kafka broker host or IP
}

// KafkaProducerAdapter implements ports.MessageQueueProducer using Sarama's SyncProducer.
type KafkaProducerAdapter struct {
	Config   *KafkaProducerConfig // Producer configuration
	Producer sarama.SyncProducer  // Sarama synchronous producer instance
}

// Connect initializes the Kafka producer with production-ready settings
// and establishes a connection to the given broker(s).
func (p *KafkaProducerAdapter) Connect(ctx context.Context) error {
	// Build broker address list from config
	brokers := []string{fmt.Sprintf("%s:%d", p.Config.Address, p.Config.Port)}

	cfg := sarama.NewConfig()

	// 1. Kafka protocol version (must match the broker version, e.g., 3.9.0 here)
	cfg.Version = sarama.V3_9_0_0

	// 2. Producer acknowledgment settings
	cfg.Producer.RequiredAcks = sarama.WaitForAll // Wait for all replicas to ack

	// 3. Partition strategy
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner // Distribute evenly across partitions

	// 4. Return behavior
	cfg.Producer.Return.Successes = true // Report successful sends
	cfg.Producer.Return.Errors = true    // Report failed sends

	// 5. Retry policy
	cfg.Producer.Retry.Max = 10                         // Retry up to 10 times
	cfg.Producer.Retry.Backoff = 100 * time.Millisecond // Delay between retries

	// 6. Timeout for waiting acknowledgments
	cfg.Producer.Timeout = 20 * time.Second

	// 7. Idempotent producer (ensures no duplicate messages on retries)
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1 // Required for idempotent mode

	// 8. Compression to reduce bandwidth usage and improve throughput
	cfg.Producer.Compression = sarama.CompressionSnappy

	// Create a new SyncProducer with the above config
	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return fmt.Errorf(
			"%w: failed to connect to Kafka Producer: %w",
			apperror.ErrUnavailable,
			err,
		)
	}

	// Assign the producer to the adapter
	p.Producer = producer

	// Log successful connection
	logger.GetLogger().Info(&logger.Log{
		Event:   "connect kafka producer",
		TraceID: config.GetTraceID(ctx),
	})

	return nil
}

// Produce sends a single message to the specified Kafka topic.
func (p *KafkaProducerAdapter) Produce(
	ctx context.Context,
	topic string,
	msg interface{},
) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: buildValueEncoder(msg), // Encode message as bytes
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace-id"),
				Value: []byte(config.GetTraceID(ctx)),
			},
			{
				Key:   []byte("content-type"),
				Value: []byte("application/json"),
			},
		},
		Metadata: map[string]string{
			"source": config.Get().AppConfig.Name,
		},
		Partition: -1,
	}

	// Send the message synchronously and capture partition/offset
	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf(
			"%w: failed to produce message: %w, partiotion:%d, offset: %d",
			apperror.ErrUnavailable,
			err,
			partition,
			offset,
		)
	}
	return nil
}

// Close gracefully shuts down the Kafka producer connection.
func (p *KafkaProducerAdapter) Close() error {
	if p.Producer != nil {
		return p.Producer.Close()
	}
	return nil
}

// Ensure KafkaProducerAdapter implements the MessageQueueProducer interface.
var _ ports.MessageQueueProducer = (*KafkaProducerAdapter)(nil)

// buildValueEncoder chooses the best Kafka encoder based on the type of input.
// - If input is []byte → use ByteEncoder (efficient for JSON/raw data).
// - If input is string → use StringEncoder (efficient for plain text).
// - Otherwise → fallback: convert to string and use StringEncoder.
func buildValueEncoder(value interface{}) sarama.Encoder {
	switch v := value.(type) {
	case []byte:
		return sarama.ByteEncoder(v)
	case string:
		return sarama.StringEncoder(v)
	default:
		return sarama.StringEncoder(fmt.Sprintf("%v", v))
	}
}
