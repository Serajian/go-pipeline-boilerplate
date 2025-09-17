package pipelines

import (
	"context"
	"sync"

	"go-pipeline/internal/ports"
)

type Runner[T any] struct {
	stages []ports.Stage[T]
}

func NewRunner[T any](stages ...ports.Stage[T]) *Runner[T] { return &Runner[T]{stages: stages} }

func (r *Runner[T]) Chain(ctx context.Context, in <-chan T) (out <-chan T, errMerged <-chan error) {
	cur := in
	errs := make([]<-chan error, len(r.stages))
	for i, s := range r.stages {
		o, e := s.Run(ctx, cur)
		cur = o
		errs[i] = e
	}
	return cur, mergeErrors(errs...)
}

var _ ports.ChainPipeline[any] = (*Runner[any])(nil)

func mergeErrors(chs ...<-chan error) <-chan error {
	out := make(chan error, 64)
	var wg sync.WaitGroup
	wg.Add(len(chs))
	for _, c := range chs {
		go func(c <-chan error) {
			defer wg.Done()
			for err := range c {
				if err != nil {
					out <- err
				}
			}
		}(c)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
