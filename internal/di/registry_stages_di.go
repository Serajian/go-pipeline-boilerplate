package di

import (
	"go-pipeline/internal/ports"
	"go-pipeline/internal/stages"
)

type RegistryStages struct {
	Validation *stages.ValidationRegistryStage
	Store      *stages.StoreRegistryStage
	Produce    *stages.ProduceRegistryStage
}

func NewRegistryStages(p ports.MessageQueueProducer) *RegistryStages {
	validation := stages.NewValidationRegistryStage()
	store := stages.NewStoreRegistryStage()
	producer := stages.NewProduceRegistryStage(p)

	return &RegistryStages{
		Validation: validation,
		Store:      store,
		Produce:    producer,
	}
}
