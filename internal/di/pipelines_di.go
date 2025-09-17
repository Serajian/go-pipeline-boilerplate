package di

import (
	"go-pipeline/internal/model"
	"go-pipeline/internal/pipelines"
)

type Pipelines struct {
	// 1) parallel (channels)
	Parallel *pipelines.Runner[model.UserData]
	// 2) barrier (channels without overlap)
	Barrier *pipelines.RunnerBarrier[model.UserData]
	// 3) fn (short)
	Short *pipelines.RunnerShortCircuit[model.UserData]
}

func NewPipelines(st *Stages) *Pipelines {
	registry := pipelines.NewRunner[model.UserData](
		st.Registry.Validation,
		st.Registry.Store,
		st.Registry.Produce,
	)

	b := pipelines.NewRunnerBarrier[model.UserData](
		16,
		st.Registry.Validation,
		st.Registry.Store,
		st.Registry.Produce,
	)

	rfn := pipelines.NewRunnerShortCircuit[model.UserData](
		st.ShortCircuits.Validation,
		st.ShortCircuits.Transform,
		st.ShortCircuits.Sink,
	)

	return &Pipelines{
		Parallel: registry,
		Barrier:  b,
		Short:    rfn,
	}
}
