package registry

import (
	"context"
	"fmt"

	"go-pipeline/config"
	"go-pipeline/infrastructure/message_queue"
	"go-pipeline/internal/ports"
	"go-pipeline/pkg/apperror"
	"go-pipeline/pkg/logger"

	"github.com/IBM/sarama"
)

type MQRegistry struct {
	producers map[string]ports.MessageQueueProducer
	consumers map[string]ports.MessageQueueConsumer
	traceID   string
}

func NewMQRegistry(ctx context.Context, handler sarama.ConsumerGroupHandler) (*MQRegistry, error) {
	mqRegistry := &MQRegistry{
		producers: make(map[string]ports.MessageQueueProducer),
		consumers: make(map[string]ports.MessageQueueConsumer),
		traceID:   config.GetTraceID(ctx),
	}

	for _, mqCFG := range config.Get().MQConfig {
		switch mqCFG.Type {
		case "producer":
			producer := &message_queue.KafkaProducerAdapter{
				Config: &message_queue.KafkaProducerConfig{
					Name:    mqCFG.Name,
					Port:    mqCFG.Port,
					Address: mqCFG.Address,
				},
			}
			if err := producer.Connect(ctx); err != nil {
				return nil, err
			}
			mqRegistry.producers[mqCFG.Name] = producer
		case "consumer":
			consumer := &message_queue.KafkaConsumerAdapter{
				Config: &message_queue.KafkaConsumerConfig{
					Name:    mqCFG.Name,
					Port:    mqCFG.Port,
					Address: mqCFG.Address,
					GroupID: mqCFG.GroupID,
					Topics:  mqCFG.Topics,
				},
				Handler: handler,
			}
			if err := consumer.Connect(ctx); err != nil {
				return nil, err
			}
			mqRegistry.consumers[mqCFG.Name] = consumer
		default:
			return nil, fmt.Errorf(`unknown MQ configuration "%s"`, mqCFG.Name)
		}
	}

	return mqRegistry, nil
}

func (r *MQRegistry) GetKafkaProducer() *message_queue.KafkaProducerAdapter {
	kafka, err := r.getProducer("kafka-producer")
	if err != nil {
		logger.GetLogger().Fatal(&logger.Log{
			Event:   "mq register producer",
			Error:   err,
			TraceID: r.traceID,
		})
	}
	return kafka.(*message_queue.KafkaProducerAdapter)
}

func (r *MQRegistry) GetKafkaConsumer() *message_queue.KafkaConsumerAdapter {
	kafka, err := r.getConsumer("kafka-consumer")
	if err != nil {
		logger.GetLogger().Fatal(&logger.Log{
			Event:   "mq register consumer",
			Error:   err,
			TraceID: r.traceID,
		})
	}
	return kafka.(*message_queue.KafkaConsumerAdapter)
}

func (r *MQRegistry) Close() error {
	for _, p := range r.producers {
		if err := p.Close(); err != nil {
			return err
		}
	}
	for _, c := range r.consumers {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *MQRegistry) getProducer(name string) (ports.MessageQueueProducer, error) {
	producer, exists := r.producers[name]
	if !exists {
		return nil, fmt.Errorf("%w: no producer found for name %s", apperror.ErrInvalidInput, name)
	}
	return producer, nil
}

func (r *MQRegistry) getConsumer(name string) (ports.MessageQueueConsumer, error) {
	consumer, exists := r.consumers[name]
	if !exists {
		return nil, fmt.Errorf("%w: no consumer found for name %s", apperror.ErrInvalidInput, name)
	}
	return consumer, nil
}
