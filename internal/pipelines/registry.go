package pipelines

import (
	"context"

	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
)

type Registry struct {
	stages []ports.RegistryStage
}

func NewRegistry(stages ...ports.RegistryStage) *Registry {
	return &Registry{stages: stages}
}

func (r *Registry) Run(ctx context.Context, in chan model.UserData) error {
	for _, stage := range r.stages {
		if err := stage.Execute(ctx, in); err != nil {
			return err
		}
		// name := stage.Name()
	}
	return nil
}

var _ ports.RegistryPipeline = (*Registry)(nil)
