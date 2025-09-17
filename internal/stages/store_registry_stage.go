package stages

import (
	"context"

	"go-pipeline/config"

	"go-pipeline/internal/ports"

	"go-pipeline/internal/model"
)

type StoreRegistryStage struct{}

func NewStoreRegistryStage() *StoreRegistryStage {
	return &StoreRegistryStage{}
}

func (s *StoreRegistryStage) Name() string {
	return "store_registry"
}

func (s *StoreRegistryStage) Run(
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
				// TODO: store to DB/Cache
				out <- m
			}
		}
	}()
	return out, err
}

var _ ports.Stage[model.UserData] = (*StoreRegistryStage)(nil)
