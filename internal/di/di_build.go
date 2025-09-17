package di

import (
	"go-pipeline/internal/pipelines"
	"go-pipeline/internal/ports"
)

type Container struct {
	Registry ports.RegistryPipeline
}

func NewContainer() *Container {
	registry := pipelines.NewRegistry()
	return &Container{Registry: registry}
}
