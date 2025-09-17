package di

import "go-pipeline/internal/ports"

type Stages struct {
	Registry *RegistryStages
}

func NewStagesContainer(p ports.MessageQueueProducer) *Stages {
	registry := NewRegistryStages(p)

	return &Stages{
		Registry: registry,
	}
}
