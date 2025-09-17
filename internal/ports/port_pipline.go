package ports

import "context"

type Pipeline[T any] interface {
	Chain(ctx context.Context, in <-chan T) (out <-chan T, errMerged <-chan error)
}
