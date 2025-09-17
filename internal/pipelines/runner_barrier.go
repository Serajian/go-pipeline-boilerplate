package pipelines

import (
	"context"
	"go-pipeline/internal/ports"
)

type RunnerBarrier[T any] struct {
	stages  []ports.Stage[T]
	buffCap int
}

func NewRunnerBarrier[T any](buffCap int, st ...ports.Stage[T]) *RunnerBarrier[T] {
	return &RunnerBarrier[T]{
		stages:  st,
		buffCap: buffCap,
	}
}

func (r *RunnerBarrier[T]) Run(ctx context.Context, in <-chan T) (<-chan T, <-chan error) {
	curIn := in
	allErrs := make([]<-chan error, len(r.stages))

	for _, stage := range r.stages {
		out, errChan := stage.Run(ctx, curIn)
		allErrs = append(allErrs, errChan)

		buffer := make([]T, 0, r.buffCap)
		drainDone := false

		for !drainDone {
			select {
			case <-ctx.Done():
				ch := make(chan T)
				close(ch)
				errCh := make(chan error)
				close(errCh)
				return ch, errCh
			case m, ok := <-out:
				if !ok {
					drainDone = true
					continue
				}
				buffer = append(buffer, m)
			case e, ok := <-errChan:
				if !ok {
					errChan = nil
					continue
				}
				_ = e
			}
		}
		nextIn := make(chan T, len(buffer))
		for _, v := range buffer {
			nextIn <- v
		}
		close(nextIn)
		curIn = nextIn
	}
	return curIn, mergeErrors(allErrs...)
}

var _ ports.PipeLineBarrier[any] = (*RunnerBarrier[any])(nil)
