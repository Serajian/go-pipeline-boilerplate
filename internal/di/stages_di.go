package di

import (
	"go-pipeline/internal/ports"
)

type Stages struct {
	Registry      *RegistryStages
	ShortCircuits *ShortCircuits
}

func NewStagesContainer(p ports.MessageQueueProducer) *Stages {
	registry := NewRegistryStages(p)
	sc := NewShortCircuits(p)

	return &Stages{
		Registry:      registry,
		ShortCircuits: sc,
	}
}
