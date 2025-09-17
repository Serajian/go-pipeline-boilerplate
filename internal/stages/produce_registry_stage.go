package stages

import (
	"context"

	"go-pipeline/config"

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

func (p *ProduceRegistryStage) Run(
	ctx context.Context,
	in <-chan model.UserData,
) (<-chan model.UserData, <-chan error) {
	out := make(chan model.UserData, config.BuffData)
	err := make(chan error, config.BuffErr)

	go func() {
		defer close(out)
		defer close(err)
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-in:
				if !ok {
					return
				}
				if errP := p.producer.Produce(ctx, "users", m); errP != nil {
					err <- errP
					continue
				}
				out <- m
			}
		}
	}()
	return out, err
}

var _ ports.Stage[model.UserData] = (*ProduceRegistryStage)(nil)
