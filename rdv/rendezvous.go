/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package rdv supports the safe and convenient execution of asynchronous computations with
// goroutines and provides facilities for the safe retrieval of the computation results.
// It provides safety in the sense that panics in asynchronous computations are transformed
// into error results and its methods and functions prevent resource leaks, race conditions,
// and deadlocks for the channels used to pass data between the parent and child goroutines.
package rdv

import (
	"context"

	"github.com/pvillela/go-rendezvous/util"
	"golang.org/x/sync/errgroup"
)

/////////////////////
// Rdv

// rdvData is the data structure used by Rdv channels.
type rdvData[T any] struct {
	value    T
	err      error
	chanOpen bool
}

// Rdv encapsulates a channel used for a function launched as a goroutine to rendezvous
// with the user of the function's results.
type Rdv[T any] struct {
	ch chan rdvData[T]
}

// Receive waits on the receiver and returns the results of the asynchronous computation for
// which the receiver was created (see Go and GoEg).
// For this method and ReceiveWatch, altogether at most one invocation is allowed for a given
// receiver.
func (rv Rdv[T]) Receive() (T, error) {
	data := <-rv.ch
	if !data.chanOpen {
		panic("attempt to get data from closed rendezvous channel")
	}
	return data.value, data.err
}

// ReceiveWatch waits on the receiver and watches the context ctx for cancellation or timeout.
// If ctx is not cancelled or times-out, this function returns the results of the asynchronous
// computation for which the receiver was created (see Go and GoEg).
// Otherwise, this function returns early with a TimeoutError or CancellationError.
// For this method and Receive, altogether at most one invocation is allowed for a given
// receiver.
func (rv Rdv[T]) ReceiveWatch(ctx context.Context) (T, error) {
	data := rdvData[T]{}
	select {
	case data = <-rv.ch:
		if !data.chanOpen {
			panic("attempt to get data from closed rendezvous channel")
		}
	case <-ctx.Done():
		data.err = ctx.Err()
	}
	return data.value, data.err
}

// Go launches f as an asynchronous computation in a goroutine and returns an Rdv instance
// to be used to retrieve the results of the computation.
func Go[T any](f func() (T, error)) Rdv[T] {
	rv := Rdv[T]{make(chan rdvData[T], 1)}
	go func() {
		defer close(rv.ch)
		fs := util.SafeFunc0E(f)
		res, err := fs()
		data := rdvData[T]{res, err, true}
		rv.ch <- data
	}()
	return rv
}

// GoEg launches f as an asynchronous computation in a goroutine associated with the
// errgroup.Group eg and returns an Rdv instance to be used to retrieve the results of
// the computation.
func GoEg[T any](eg *errgroup.Group, f func() (T, error)) Rdv[T] {
	rv := Rdv[T]{make(chan rdvData[T], 1)}
	eg.Go(func() error {
		defer close(rv.ch)
		fs := util.SafeFunc0E(f)
		res, err := fs()
		data := rdvData[T]{res, err, true}
		rv.ch <- data
		return err
	})
	return rv
}

// CtxApply closes function f over the ctx argument to return a nulladic function.
func CtxApply[T any](
	ctx context.Context,
	f func(context.Context) (T, error),
) func() (T, error) {
	return func() (T, error) {
		return f(ctx)
	}
}

// CtxApplyWatch closes function f over the ctx argument to return a nulladic function and watches
// ctx for deadline expiration or cancellation.
// If ctx is not cancelled or times-out, the resulting function returns the results of f.
// Otherwise, the resulting function returns early with a TimeoutError or CancellationError.
func CtxApplyWatch[T any](
	ctx context.Context,
	f func(context.Context) (T, error),
) func() (T, error) {
	fc := CtxApply(ctx, f)
	return func() (T, error) {
		rv := Go(fc)
		return rv.ReceiveWatch(ctx)
	}
}
