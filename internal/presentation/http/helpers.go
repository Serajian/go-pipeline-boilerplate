package http

import (
	"context"
	"go-pipeline/internal/model"
)

func drainAll(
	ctx context.Context,
	out <-chan model.UserData,
	errCh <-chan error,
) (items []model.UserData, errs []error, canceled bool) {

	for out != nil || errCh != nil {
		select {
		case <-ctx.Done():
			return items, errs, true
		case m, ok := <-out:
			if !ok {
				out = nil
				continue
			}
			items = append(items, m)
		case e, ok := <-errCh:
			if !ok {
				errCh = nil
				continue
			}
			if e != nil {
				errs = append(errs, e)
			}
		}
	}
	return items, errs, false
}
