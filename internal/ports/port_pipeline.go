package ports

import "context"

type ChainPipeline[T any] interface {
	Chain(ctx context.Context, in <-chan T) (out <-chan T, errMerged <-chan error)
}

type ShortCircuitPipeLine[T any] interface {
	Run(ctx context.Context, m T) (T, error)
}

type BarrierPipeLine[T any] interface {
	Run(ctx context.Context, in <-chan T) (finalOut <-chan T, mergedErr <-chan error)
}
