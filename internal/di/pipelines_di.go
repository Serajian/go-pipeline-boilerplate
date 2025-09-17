package di

import (
	"go-pipeline/internal/model"
	"go-pipeline/internal/pipelines"
)

type Pipelines struct {
	Registry *pipelines.Runner[model.UserData]
}

func NewPipelines(stages *Stages) *Pipelines {
	r := pipelines.NewRunner[model.UserData](
		stages.Registry.Validation,
		stages.Registry.Store,
		stages.Registry.Produce,
	)
	return &Pipelines{
		Registry: r,
	}
}
