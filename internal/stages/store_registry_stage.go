package stages

import (
	"context"
	"errors"

	"go-pipeline/internal/model"
)

type StoreRegistryStage struct{}

func NewStoreRegistryStage() *StoreRegistryStage {
	return &StoreRegistryStage{}
}

func (s *StoreRegistryStage) Name() string {
	return "store_registry"
}

func (s *StoreRegistryStage) Execute(ctx context.Context, in chan model.UserData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case userData, ok := <-in:
		if !ok {
			return errors.New("user data channel closed")
		}
		in <- userData
	default:
		return nil
	}
	return nil
}
