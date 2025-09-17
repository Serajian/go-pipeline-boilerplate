package di

import (
	"go-pipeline/internal/model"
	"go-pipeline/internal/ports"
	"go-pipeline/internal/stages"
)

type ShortCircuits struct {
	Validation ports.StageFn[model.UserData]
	Transform  ports.StageFn[model.UserData]
	Sink       ports.StageFn[model.UserData]
}

func NewShortCircuits(p ports.MessageQueueProducer) *ShortCircuits {
	return &ShortCircuits{
		Validation: stages.ValidationFn(),
		Transform:  stages.TransformFn(),
		Sink:       stages.SinkFn(p),
	}
}
