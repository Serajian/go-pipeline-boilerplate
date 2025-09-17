package stages

import (
	"context"
	"errors"

	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
)

type ProduceRegistryStage struct {
	producer ports.MessageQueueProducer
}

func NewProduceRegistryStage(producer ports.MessageQueueProducer) *ProduceRegistryStage {
	return &ProduceRegistryStage{producer: producer}
}

func (p *ProduceRegistryStage) Name() string { return "produce-registry" }

func (p *ProduceRegistryStage) Execute(ctx context.Context, in chan model.UserData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case userData, ok := <-in:
		if !ok {
			return errors.New("input channel closed unexpectedly")
		}
		return p.producer.Produce(ctx, "users", userData)
	default:
		return nil
	}
}

var _ ports.RegistryStage = (*ProduceRegistryStage)(nil)
