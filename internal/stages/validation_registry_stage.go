package stages

import (
	"context"
	"errors"

	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
)

type ValidationRegistryStage struct{}

func NewValidationRegistryStage() *ValidationRegistryStage {
	return &ValidationRegistryStage{}
}

func (v *ValidationRegistryStage) Name() string {
	return "validation_registry"
}

func (v *ValidationRegistryStage) Run(
	ctx context.Context,
	in <-chan model.UserData,
) (<-chan model.UserData, <-chan error) {
	// TODO: from config
	out := make(chan model.UserData, 64)
	err := make(chan error, 64)

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
				if m.Email == "" {
					err <- errors.New("email address is required")
					continue
				}
				if !containsAt(m.Email) {
					err <- errors.New("email address is invalid")
					continue
				}
				out <- m
			}
		}
	}()
	return out, err
}

var _ ports.Stage[model.UserData] = (*ValidationRegistryStage)(nil)

func containsAt(s string) bool {
	for i := range s {
		if s[i] == '@' {
			return true
		}
	}
	return false
}
