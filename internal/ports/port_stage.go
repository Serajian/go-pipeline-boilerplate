package ports

import (
	"context"
)

// Stage represents a processing step used in Parallel or Barrier pipelines.
// Each stage has a unique name and consumes values from the `in` channel.
// The processed results are sent to the `out` channel, while any errors
// are reported through the `err` channel. Both output channels must be
// closed by the stage when processing is complete.
type Stage[T any] interface {
	Name() string
	Run(ctx context.Context, in <-chan T) (out <-chan T, err <-chan error)
}

// StageFn represents a function-based stage for short-circuit pipelines.
// Unlike Stage, this processes a single value at a time rather than streams
// from a channel. If an error is returned, the pipeline is immediately
// stopped (short-circuited) and the error is propagated.
type StageFn[T any] func(ctx context.Context, m T) (T, error)
