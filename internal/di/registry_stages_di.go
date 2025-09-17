package di

import (
	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
	"go-pipeline/internal/stages"
)

type RegistryStages struct {
	Validation ports.Stage[model.UserData]
	Store      ports.Stage[model.UserData]
	Produce    ports.Stage[model.UserData]
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
