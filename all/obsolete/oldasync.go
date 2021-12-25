/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package obsolete

import (
	"context"
	"fmt"
	"sync"

	"github.com/pvillela/go-rendezvous/util"
	"golang.org/x/sync/errgroup"
)

/////////////////////
// Unit

// Unit is a type alias
type Unit = struct{}

/////////////////////
// ResultWithError

// ResultWithError encapsulates a normal result value and an error.
type ResultWithError[T any] struct {
	Value T
	Error error
}

/////////////////////
// Waiters

// preWaiter encapsulates a context.Context and a chan error.  It supports waiting on the
// completion of an asynchronous operation while paying attention to the context's deadline.
// See its ErrWait() method.
type preWaiter struct {
	Ctx        context.Context
	WaitCh     chan error
	ActiveGate chan bool
	WaitChGate *sync.Once
	Err        error
}

// ErrWait waits on both the WaitCh channel field in the receiver and the Done() channel for the
// Ctx context field in the receiver.  If Ctx.Done() happens first then this method returns
// a DeadlineExceeded or a Canceled error.  Otherwise, it returns the error received in WaitCh.
// This function may be called multiple times and always returns the same results.
func (w preWaiter) ErrWait() error {
	// fmt.Println("%%% entered ErrWait")
	w.WaitChGate.Do(func() {
		// fmt.Println("%%% in WaitChGate")
		err := w.Err
		select {
		case errW := <-w.WaitCh:
			err = errW
		case <-w.Ctx.Done():
			err = w.Ctx.Err()
		}
		w.Err = err
	})
	return w.Err
}

// SingleWaiter supports waiting on the completion of a single goroutine.
type SingleWaiter struct {
	P       preWaiter // embedded structs not supported by go2go
	ComplCh chan Unit
}

// MakeSingleWaiter constructs a SingleWaiter.
func MakeSingleWaiter(ctx context.Context) SingleWaiter {
	complCh := make(chan Unit)
	waitCh := make(chan error)
	activeGate := make(chan bool, 1)
	activeGate <- true
	waitChGate := new(sync.Once)
	return SingleWaiter{preWaiter{ctx, waitCh, activeGate, waitChGate, nil}, complCh}
}

// Wait waits on both the WaitCh channel field in the receiver and the Done() channel for the
// Ctx context field in the receiver.  If Ctx.Done() happens first then this method returns
// a DeadlineExceeded or a Canceled error.  Otherwise, it returns the error received in WaitCh.
// This function may be called multiple times and always returns the same results.
func (w SingleWaiter) Wait() error {
	// fmt.Println("%%% entered Wait")
	if active := <-w.P.ActiveGate; active {
		// fmt.Println("%%% executing ActiveGate")
		close(w.P.ActiveGate)
		go func() {
			<-w.ComplCh
			close(w.P.WaitCh)
		}()
	}
	return w.P.ErrWait()
}

// WgWaiter supports waiting on the completion of a group of goroutines associated with a
// sync.WaitGroup.
type WgWaiter struct {
	P  preWaiter // embedded structs not supported by go2go
	Wg *sync.WaitGroup
}

// MakeWgWaiter constructs a WgWaiter.
func MakeWgWaiter(ctx context.Context) WgWaiter {
	wg := new(sync.WaitGroup)
	waitCh := make(chan error)
	activeGate := make(chan bool, 1)
	activeGate <- true
	waitChGate := new(sync.Once)
	return WgWaiter{preWaiter{ctx, waitCh, activeGate, waitChGate, nil}, wg}
}

// Wait waits on both the WaitCh channel field in the receiver and the Done() channel for the
// Ctx context field in the receiver.  If Ctx.Done() happens first then this method returns
// a DeadlineExceeded or a Canceled error.  Otherwise, it returns the error received in WaitCh.
// This function may be called multiple times and always returns the same results.
func (w WgWaiter) Wait() error {
	if active := <-w.P.ActiveGate; active {
		close(w.P.ActiveGate)
		go func() {
			w.Wg.Wait()
			close(w.P.WaitCh)
		}()
	}
	return w.P.ErrWait()
}

// EgWaiter supports waiting on the completion of a group of goroutines associated with an
// errgroup.Group.
type EgWaiter struct {
	P     preWaiter // embedded structs not supported by go2go
	Eg    *errgroup.Group
	EgCtx context.Context
}

// MakeEgWaiter constructs an EgWaiter.
func MakeEgWaiter(parentCtx context.Context) EgWaiter {
	eg, egCtx := errgroup.WithContext(parentCtx)
	waitCh := make(chan error)
	activeGate := make(chan bool, 1)
	activeGate <- true
	waitChGate := new(sync.Once)
	return EgWaiter{preWaiter{parentCtx, waitCh, activeGate, waitChGate, nil}, eg, egCtx}
}

// Wait waits on both the WaitCh channel field in the receiver and the Done() channel for the
// Ctx context field in the receiver.  If Ctx.Done() happens first then this method returns
// a DeadlineExceeded or a Canceled error.  Otherwise, it returns the error received in WaitCh.
// This function may be called multiple times and always returns the same results.
func (w EgWaiter) Wait() error {
	if active := <-w.P.ActiveGate; active {
		close(w.P.ActiveGate)
		go func() {
			err := w.Eg.Wait()
			w.P.WaitCh <- err
			close(w.P.WaitCh)
		}()
	}
	return w.P.ErrWait()
}

/////////////////////
// SafeGo

// Go launches f as a saparate goroutine and puts its result in the struct pointed to by
// pResult.
// Before using the result, a caller must wait on the Waiter returned by this function.
// If f panics, the panic value is converted to an error and set in the result.
func Go[T any](
	ctx context.Context,
	f func(context.Context) (T, error),
	pResult *ResultWithError[T],
) SingleWaiter {
	w := MakeSingleWaiter(ctx)
	// fmt.Println("%%%% entered SafeGoCh")
	if active := <-w.P.ActiveGate; active {
		// fmt.Println("%%%% in active portion of SafeGoCh")
		defer func() { w.P.ActiveGate <- true }() // keep it active
		fs := util.SafeFunc1E(f)
		go func() {
			res, err := fs(w.P.Ctx)
			// fmt.Println("%%%% res, err :=", res, err)
			*pResult = ResultWithError[T]{res, err}
			fmt.Println("%%%% *pResult", *pResult)
			close(w.ComplCh)
		}()
	} else {
		panic("attemp to add goroutine to inactive waiter")
	}
	return w
}

// GoWg launches f as a saparate goroutine and puts its result in the struct pointed to by
// pResult.
// Before using the result, a caller must wait on the Waiter passed into this function.
// If f panics, the panic value is converted to an error and set in the result.
func GoWg[T any](
	w WgWaiter,
	f func(context.Context) (T, error),
	pResult *ResultWithError[T],
) {
	if active := <-w.P.ActiveGate; active {
		defer func() { w.P.ActiveGate <- true }() // keep it active
		w.Wg.Add(1)
		fs := util.SafeFunc1E(f)
		go func() {
			defer w.Wg.Done()
			res, err := fs(w.P.Ctx)
			*pResult = ResultWithError[T]{res, err}
		}()
	} else {
		panic("attemp to add goroutine to inactive waiter")
	}
}

// GoEg launches f as a saparate goroutine and puts its non-error result in the struct
// pointed to by pResult.
// Before using the result, a caller must wait on the Waiter passed into this function.
// If f panics, the panic value is converted to an error and set in the result.
func GoEg[T any](
	w EgWaiter,
	f func(context.Context) (T, error),
	pResult *T,
) {
	if active := <-w.P.ActiveGate; active {
		defer func() { w.P.ActiveGate <- true }() // keep it active
		fs := util.SafeFunc1E(f)
		w.Eg.Go(func() error {
			res, err := fs(w.EgCtx)
			*pResult = res
			return err
		})
	} else {
		panic("attemp to add goroutine to inactive waiter")
	}
}

/////////////////////
// RunConcurrent

// RunConcurrentsWg runs funcs concurrently and returns a slice containing the results of
// the function executions once all functions complete normaly, with an error, or with a panic.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError.
func RunConcurrentsWg[T any](
	ctx context.Context,
	funcs ...func(context.Context) (T, error),
) ([]ResultWithError[T], error) {
	results := make([]ResultWithError[T], len(funcs))
	waiter := MakeWgWaiter(ctx)
	for index, f := range funcs {
		GoWg(waiter, f, &results[index])
	}
	err := waiter.Wait()
	fmt.Println("%%% results =", results, "err = ", err)
	return results, err
}

// RunConcurrent2Wg runs funcs concurrently and returns a slice containing the results of
// the function executions once all functions complete normaly, with an error, or with a panic.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError.
func RunConcurrent2Wg[T1, T2 any](
	ctx context.Context,
	f1 func(context.Context) (T1, error),
	f2 func(context.Context) (T2, error),
) (util.Tuple2[ResultWithError[T1], ResultWithError[T2]], error) {
	results := util.Tuple2[ResultWithError[T1], ResultWithError[T2]]{}
	waiter := MakeWgWaiter(ctx)
	GoWg(waiter, f1, &results.X1)
	GoWg(waiter, f2, &results.X2)
	err := waiter.Wait()
	return results, err
}

// RunConcurrentsEg runs funcs concurrently and returns a slice containing the non-error results
// of the function executions if all functions complete normaly.  If any of the functions
// returns an error or panics, this function returns early, with the first error encountered.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError.
func RunConcurrentsEg[T any](
	ctx context.Context,
	funcs ...func(context.Context) (T, error),
) ([]T, error) {
	results := make([]T, len(funcs))
	waiter := MakeEgWaiter(ctx)
	for index, f := range funcs {
		GoEg(waiter, f, &results[index])
	}
	err := waiter.Wait()
	return results, err
}

// RunConcurrent2Eg runs funcs concurrently and returns a slice containing the non-error results
// of the function executions if all functions complete normaly.  If any of the functions
// returns an error or panics, this function returns early, with the first error encountered.
// Panics in function executions are converted to errors.
// In case of a context timeout or cancellation, this functionn returns early with a
// TimeoutError or CancellationError.
func RunConcurrent2Eg[T1, T2 any](
	ctx context.Context,
	f1 func(context.Context) (T1, error),
	f2 func(context.Context) (T2, error),
) (util.Tuple2[T1, T2], error) {
	res := util.Tuple2[T1, T2]{}
	waiter := MakeEgWaiter(ctx)
	GoEg(waiter, f1, &res.X1)
	GoEg(waiter, f2, &res.X2)
	err := waiter.Wait()
	return res, err
}

/////////////////////
// Promise

// Promise supports awaiting on and receiving the result of an asynchronous computation.
type Promise[T any] interface {
	Await() (T, error)
}

// promisImpl is the concrete implementation of Promise.
type promiseImpl[T any] struct {
	ResWE  ResultWithError[T] // can't use embedded struct with go2go
	Waiter SingleWaiter
}

// Await waits on the asychronous computation associated with the receiver and returns the
// result of the computation.
// This function may be called multiple times and always returns the same results.
func (P *promiseImpl[T]) Await() (T, error) {
	// fmt.Println("%%% about to wait")
	errW := P.Waiter.Wait()
	value := P.ResWE.Value
	err := P.ResWE.Error
	if errW != nil {
		err = errW
	}
	return value, err
}

// Async constructs a Promise for the asynchronous execution of f.
func Async[T any](
	ctx context.Context,
	f func(context.Context) (T, error),
) Promise[T] {
	promImpl := promiseImpl[T]{}
	pResWE := &promImpl.ResWE
	promImpl.Waiter = Go(ctx, f, pResWE)
	return &promImpl
}

// Async2Eg returns a Promise for the concurrent execution of the functions f1 and f2
// in an errgroup.Group.
// The Promise completes when the error group Wait method returns.
func Async2Eg[T1 any, T2 any](
	ctx context.Context,
	f1 func(ctx context.Context) (T1, error),
	f2 func(ctx context.Context) (T2, error),
) Promise[util.Tuple2[T1, T2]] {
	f := func(ctx context.Context) (util.Tuple2[T1, T2], error) {
		return RunConcurrent2Eg[T1, T2](ctx, f1, f2)
	}
	return Async(ctx, f)
}

// AsyncsEg returns a Promise for the concurrent execution of the functions funcs
// in an errgroup.Group.
// The Promise completes when the error group Wait method returns.
func AsyncsEg[T any](
	ctx context.Context,
	funcs ...func(ctx context.Context) (T, error),
) Promise[[]T] {
	f := func(ctx context.Context) ([]T, error) {
		return RunConcurrentsEg[T](ctx, funcs...)
	}
	return Async(ctx, f)
}

// Async2Wg returns a Promise for the concurrent execution of the functions f1 and f2
// in a sync.WaitGroup.
// The Promise completes when the wait group Wait method returns.
func Async2Wg[T1 any, T2 any](
	ctx context.Context,
	f1 func(ctx context.Context) (T1, error),
	f2 func(ctx context.Context) (T2, error),
) Promise[util.Tuple2[ResultWithError[T1], ResultWithError[T2]]] {
	f := func(ctx context.Context) (util.Tuple2[ResultWithError[T1], ResultWithError[T2]], error) {
		return RunConcurrent2Wg[T1, T2](ctx, f1, f2)
	}
	return Async(ctx, f)
}

// AsyncsWg returns a Promise for the concurrent execution of the functions funcs
// in a sync.WaitGroup.
// The Promise completes when the wait group Wait method returns.
func AsyncsWg[T any](
	ctx context.Context,
	funcs ...func(ctx context.Context) (T, error),
) Promise[[]ResultWithError[T]] {
	f := func(ctx context.Context) ([]ResultWithError[T], error) {
		return RunConcurrentsWg[T](ctx, funcs...)
	}
	return Async(ctx, f)
}

/////////////////////
// RunInBackground

// RunInBackground runs a function in the background as a goroutine.  This function returns
// early if the context ctx is cancelled or times-out.  In all cases, errorHandler is executed
// if f returns an error, even if f completes after the context ctx is cancelled.
func RunInBackground[T any](
	ctx context.Context,
	f func(context.Context) (T, error),
	errorHandler func(error),
) (T, error) {
	f1 := func(ctx context.Context) (T, error) {
		res, err := f(ctx)
		// is cancelled or times-out and f continues executing thene
		if err != nil {
			errorHandler(err)
		}
		return res, err
	}
	promise := Async(ctx, f1)
	return promise.Await()
}
