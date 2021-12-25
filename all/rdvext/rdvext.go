/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Simple extensions to rdv package to support running groups of functions concurrently.
package rdvext

import (
	"context"

	"github.com/pvillela/go-rendezvous/rdv"
	"github.com/pvillela/go-rendezvous/util"
	"golang.org/x/sync/errgroup"
)

/////////////////////
// ResultWithError

// ResultWithError encapsulates a normal result value and an error.
type ResultWithError[T any] struct {
	Value T
	Error error
}

/////////////////////
// Run multiple

// RunSlice runs funcs concurrently and returns a slice containing the results of
// the function executions once all functions complete normaly, with an error, or with a panic.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError for each of the funcs that had not aready returned.
// If there are any errors, the returned error is the one associated with the first function
// in the list of aguments that has an error response (not necessarily the first function to
// return an error).
func RunSlice[T any](
	ctx context.Context,
	funcs ...func(context.Context) (T, error),
) ([]ResultWithError[T], error) {
	rvs := make([]rdv.Rdv[T], len(funcs))
	for i, f := range funcs {
		rvs[i] = rdv.Go(rdv.CtxApply(ctx, f))
	}

	results := make([]ResultWithError[T], len(funcs))
	for i := 0; i < len(rvs); i++ {
		results[i].Value, results[i].Error = rvs[i].ReceiveWatch(ctx)
	}

	var err error = nil
	for _, res := range results {
		if res.Error != nil {
			err = res.Error
			break
		}
	}

	return results, err
}

// Run2 runs funcs concurrently and returns a tuple containing the results of
// the function executions once all functions complete normaly, with an error, or with a panic.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError for each of the funcs that had not aready returned.
// If there are any errors, the returned error is the one associated with the first function
// in the list of aguments that has an error response (not necessarily the first function to
// return an error).
func Run2[T1, T2 any](
	ctx context.Context,
	f1 func(context.Context) (T1, error),
	f2 func(context.Context) (T2, error),
) (util.Tuple2[ResultWithError[T1], ResultWithError[T2]], error) {
	rv1 := rdv.Go(rdv.CtxApply(ctx, f1))
	rv2 := rdv.Go(rdv.CtxApply(ctx, f2))

	results := util.Tuple2[ResultWithError[T1], ResultWithError[T2]]{}
	results.X1.Value, results.X1.Error = rv1.ReceiveWatch(ctx)
	results.X2.Value, results.X2.Error = rv2.ReceiveWatch(ctx)

	var err error = nil
	errs := []error{results.X1.Error, results.X2.Error}
	for _, e := range errs {
		if e != nil {
			err = e
			break
		}
	}

	return results, err
}

// RunSliceEg runs funcs concurrently and returns a slice containing the non-error results
// of the function executions if all functions complete normaly.  If any of the functions
// returns an error or panics, this function returns early, with the first error encountered.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError.
func RunSliceEg[T any](
	ctx context.Context,
	funcs ...func(context.Context) (T, error),
) ([]T, error) {
	eg, egCtx := errgroup.WithContext(ctx)
	rvs := make([]rdv.Rdv[T], len(funcs))
	for i, f := range funcs {
		rvs[i] = rdv.GoEg(eg, rdv.CtxApplyWatch(egCtx, f))
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	results := make([]T, len(funcs))
	for i := 0; i < len(rvs); i++ {
		results[i], _ = rvs[i].Receive()
	}

	return results, err
}

// Run2Eg runs funcs concurrently and returns a tuple containing the non-error results
// of the function executions if all functions complete normaly.  If any of the functions
// returns an error or panics, this function returns early, with the first error encountered.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError.
func Run2Eg[T1, T2 any](
	ctx context.Context,
	f1 func(context.Context) (T1, error),
	f2 func(context.Context) (T2, error),
) (util.Tuple2[T1, T2], error) {
	eg, egCtx := errgroup.WithContext(ctx)
	rv1 := rdv.GoEg(eg, rdv.CtxApplyWatch(egCtx, f1))
	rv2 := rdv.GoEg(eg, rdv.CtxApplyWatch(egCtx, f2))

	results := util.Tuple2[T1, T2]{}

	err := eg.Wait()
	if err != nil {
		return results, err
	}

	results.X1, _ = rv1.ReceiveWatch(ctx)
	results.X2, _ = rv2.ReceiveWatch(ctx)

	return results, err
}

/////////////////////
// Go multiple

// GoSlice returns an rdv.Rdv for the concurrent execution of the functions funcs
// in a sync.WaitGroup.
// The rdv.Rdv encapsulates a slice containing the results of
// the function executions once all functions complete normaly, with an error, or with a panic.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, the rdv.Rdv completes early with a
// TimeoutError or CancellationError.
func GoSlice[T any](
	ctx context.Context,
	funcs ...func(ctx context.Context) (T, error),
) rdv.Rdv[[]ResultWithError[T]] {
	f := func() ([]ResultWithError[T], error) {
		return RunSlice[T](ctx, funcs...)
	}
	return rdv.Go(f)
}

// Go2 returns an rdv.Rdv for the concurrent execution of the functions f1 and f2.
// The rdv.Rdv encapsulates a tuple containing the results of
// the function executions once all functions complete normaly, with an error, or with a panic.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, the rdv.Rdv completes early with a
// TimeoutError or CancellationError.
func Go2[T1 any, T2 any](
	ctx context.Context,
	f1 func(ctx context.Context) (T1, error),
	f2 func(ctx context.Context) (T2, error),
) rdv.Rdv[util.Tuple2[ResultWithError[T1], ResultWithError[T2]]] {
	f := func() (util.Tuple2[ResultWithError[T1], ResultWithError[T2]], error) {
		return Run2[T1, T2](ctx, f1, f2)
	}
	return rdv.Go(f)
}

// GoSliceEg returns an rdv.Rdv for the concurrent execution of the functions funcs
// in an errgroup.Group.
// The rdv.Rdv encapsulates a slice containing the non-error results
// of the function executions if all functions complete normaly.  If any of the functions
// returns an error or panics, this function returns early, with the first error encountered.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, the rdv.Rdv completes early with a
// TimeoutError or CancellationError.
func GoSliceEg[T any](
	ctx context.Context,
	funcs ...func(ctx context.Context) (T, error),
) rdv.Rdv[[]T] {
	f := func() ([]T, error) {
		return RunSliceEg[T](ctx, funcs...)
	}
	return rdv.Go(f)
}

// Go2Eg returns an rdv.Rdv for the concurrent execution of the functions f1 and f2
// in an errgroup.Group.
// The rdv.Rdv encapsulates a tuple containing the non-error results
// of the function executions if all functions complete normaly.  If any of the functions
// returns an error or panics, this function returns early, with the first error encountered.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, the rdv.Rdv completes early with a
// TimeoutError or CancellationError.
func Go2Eg[T1 any, T2 any](
	ctx context.Context,
	f1 func(ctx context.Context) (T1, error),
	f2 func(ctx context.Context) (T2, error),
) rdv.Rdv[util.Tuple2[T1, T2]] {
	f := func() (util.Tuple2[T1, T2], error) {
		return Run2Eg[T1, T2](ctx, f1, f2)
	}
	return rdv.Go(f)
}
