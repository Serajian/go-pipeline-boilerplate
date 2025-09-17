package di

import "go-pipeline/internal/pipelines"

type Pipelines struct {
	Registry *pipelines.Registry
}

func NewPipelines(stages *Stages) *Pipelines {
	registry := pipelines.NewRegistry(
		stages.Registry.Validation,
		stages.Registry.Store,
		stages.Registry.Produce,
	)
	return &Pipelines{
		Registry: registry,
	}
}
