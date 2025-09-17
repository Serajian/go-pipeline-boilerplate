package ports

import (
	"context"
)

type Stage[T any] interface {
	Name() string
	Run(ctx context.Context, in <-chan T) (out <-chan T, err <-chan error)
}
