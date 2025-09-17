package ports

import "context"

// ChainPipeline defines a parallel (concurrent) pipeline.
// Each stage processes items concurrently, and all results
// are sent downstream through the `out` channel.
// Any errors from the stages are collected and merged
// into the `errMerged` channel. The entire execution
// is controlled by the provided context.
type ChainPipeline[T any] interface {
	Chain(ctx context.Context, in <-chan T) (out <-chan T, errMerged <-chan error)
}

// ShortCircuitPipeLine defines a sequential pipeline that
// processes a single value step by step. If any stage
// returns an error, execution is immediately stopped
// (short-circuited) and the error is returned.
// This is useful for validation or scenarios where
// failure should prevent further processing.
type ShortCircuitPipeLine[T any] interface {
	Run(ctx context.Context, m T) (T, error)
}

// BarrierPipeLine defines a pipeline that processes inputs
// concurrently but waits at a synchronization point (barrier)
// until all stages have finished. Only when all processing
// is complete are the final results sent through the `finalOut`
// channel, and all errors are merged into the `mergedErr` channel.
// This pattern is useful when you need all results before proceeding.
type BarrierPipeLine[T any] interface {
	Run(ctx context.Context, in <-chan T) (finalOut <-chan T, mergedErr <-chan error)
}
