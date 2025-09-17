package pipelines

import (
	"context"
	"go-pipeline/internal/ports"
)

type RunnerShortCircuit[T any] struct {
	stages []ports.StageFn[T]
}

func NewRunnerShortCircuit[T any](stages ...ports.StageFn[T]) *RunnerShortCircuit[T] {
	return &RunnerShortCircuit[T]{
		stages: stages,
	}
}

func (r *RunnerShortCircuit[T]) Run(ctx context.Context, m T) (T, error) {
	cur := m
	for _, stage := range r.stages {
		next, err := stage(ctx, m)
		if err != nil {
			return cur, err
		}
		cur = next
	}
	return cur, nil
}
