package stages

import (
	"context"
	"errors"
	"fmt"

	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
)

type ValidationRegistryStage struct{}

func (v *ValidationRegistryStage) Name() string {
	return "validation_registry"
}

func (v *ValidationRegistryStage) Execute(ctx context.Context, in <-chan model.UserData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case userData, ok := <-in:
		if !ok {
			return errors.New("user data channel closed")
		}
		fmt.Println(userData)
	}
	return nil
}

var _ ports.RegistryStage = (*ValidationRegistryStage)(nil)
