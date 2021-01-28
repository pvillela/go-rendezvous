package runconcurrent

import (
	"context"
	"golang.org/x/sync/errgroup"
)
import "golang.org/x/sync/semaphore"

// This is modeled on my Kotlin Promixe
type DeferredSliceInt struct {
	Values              []int
	Error               error
	awaitSemaphore      *semaphore.Weighted
	completionSemaphore *semaphore.Weighted  // only needed to prevent multiple completion of the promise
}

func AsyncInt(callerCtx context.Context, funcs ...func(ctx context.Context) (int, error)) DeferredSliceInt {
	eg, ctx := errgroup.WithContext(callerCtx)
	var deferred DeferredSliceInt
	for i, f := range funcs {
		eg.Go(func() error {
			res, err := f(ctx)
			if err != nil {
				deferred.Error = err
				return err
			}
			deferred.Values[i] = res
			return nil
		})
	}
	return deferred
}
