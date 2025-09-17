package ports

import "context"

type Pipeline[T any] interface {
	Chain(ctx context.Context, in <-chan T) (out <-chan T, errMerged <-chan error)
}

type PipeLineFn[T any] interface {
	Run(ctx context.Context, m T) (T, error)
}

type PipeLineBarrier[T any] interface {
	Run(ctx context.Context, in <-chan T) (finalOut <-chan T, mergedErr <-chan error)
}
