package ports

import (
	"context"
)

// Stage for Parallelism and Barrier
type Stage[T any] interface {
	Name() string
	Run(ctx context.Context, in <-chan T) (out <-chan T, err <-chan error)
}

// StageFn for short-circuit
type StageFn[T any] func(ctx context.Context, m T) (T, error)
