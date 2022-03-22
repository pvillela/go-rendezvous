/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package util

import (
	"context"
	"fmt"
	"time"
)

// Unit is a type alias representing the standard functional programming Unit type
type Unit = struct{}

// ErrorOf is a type that wraps an arbitrary value to turn the value into an error instance
type ErrorOf struct {
	Value interface{}
}

// Error implements the error interface
func (err ErrorOf) Error() string {
	return fmt.Sprintf("%v", err.Value)
}

// ToError transforms an arbitrary value x into an error. If x is an error, it does nothing.
// Otherwise, it wraps x in an ErrorOf.
func ToError(x interface{}) error {
	switch x.(type) {
	case error:
		return x.(error)
	default:
		return ErrorOf{x}
	}
}

// Tuple2 is tuple with 2 elements
type Tuple2[T1, T2 any] struct {
	X1 T1
	X2 T2
}

// Tuple3 is tuple with 3 elements
type Tuple3[T1, T2, T3 any] struct {
	X1 T1
	X2 T2
	X3 T3
}

// Tuple4 is tuple with 4 elements
type Tuple4[T1, T2, T3, T4 any] struct {
	X1 T1
	X2 T2
	X3 T3
	X4 T4
}

// SafeFunc0E returns a function that never panics.
// That function returns the same values as f if f doesn't panic and returns an error if f panics.
func SafeFunc0E[U any](f func() (U, error)) func() (U, error) {
	return func() (res U, err error) {
		defer func() {
			err0 := recover()
			if err0 != nil {
				err = ToError(err0)
			}
		}()
		return f()
	}
}

// SafeFunc0VE returns a function that never panics.
// That function returns the same value as f if f doesn't panic and returns an error if f panics.
func SafeFunc0VE(f func() error) func() error {
	fu := func() (Unit, error) { return Unit{}, f() }
	return func() error {
		_, err := SafeFunc0E(fu)()
		return err
	}
}

// SafeFunc0 returns a function that never panics.
// That function returns the same value as f if f doesn't panic and returns an error if f panics.
func SafeFunc0[U any](f func() U) func() (U, error) {
	fe := func() (U, error) { return f(), nil }
	return SafeFunc0E(fe)
}

// SafeFunc0V returns a function that never panics.
// That function returns nil if f doesn't panic and returns an error if f panics.
func SafeFunc0V(f func()) func() error {
	fe := func() error { f(); return nil }
	return SafeFunc0VE(fe)
}

// SafeFunc1E returns a function that never panics.
// That function returns the same values as f if f doesn't panic and returns an error if f panics.
func SafeFunc1E[T1, U any](f func(T1) (U, error)) func(T1) (U, error) {
	return func(t1 T1) (res U, err error) {
		defer func() {
			err0 := recover()
			if err0 != nil {
				err = ToError(err0)
			}
		}()
		return f(t1)
	}
}

// SafeFunc1VE returns a function that never panics.
// That function returns the same value as f if f doesn't panic and returns an error if f panics.
func SafeFunc1VE[T1 any](f func(T1) error) func(T1) error {
	fu := func(t1 T1) (Unit, error) { return Unit{}, f(t1) }
	return func(t1 T1) error {
		_, err := SafeFunc1E(fu)(t1)
		return err
	}
}

// SafeFunc1 returns a function that never panics.
// That function returns the same value as f if f doesn't panic and returns an error if f panics.
func SafeFunc1[T1, U any](f func(T1) U) func(T1) (U, error) {
	fe := func(t1 T1) (U, error) { return f(t1), nil }
	return SafeFunc1E(fe)
}

// SafeFunc1V returns a function that never panics.
// That function returns nil if f doesn't panic and returns an error if f panics.
func SafeFunc1V[T1 any](f func(T1)) func(T1) error {
	fe := func(t1 T1) error { f(t1); return nil }
	return SafeFunc1VE(fe)
}

// SafeFunc2E returns a function that never panics.
// That function returns the same values as f if f doesn't panic and returns an error if f panics.
func SafeFunc2E[T1, T2, U any](f func(T1, T2) (U, error)) func(T1, T2) (U, error) {
	return func(t1 T1, t2 T2) (res U, err error) {
		defer func() {
			err0 := recover()
			if err0 != nil {
				err = ToError(err0)
			}
		}()
		return f(t1, t2)
	}
}

// SafeFunc2VE returns a function that never panics.
// That function returns the same value as f if f doesn't panic and returns an error if f panics.
func SafeFunc2VE[T1, T2 any](f func(T1, T2) error) func(T1, T2) error {
	fu := func(t1 T1, t2 T2) (Unit, error) { return Unit{}, f(t1, t2) }
	return func(t1 T1, t2 T2) error {
		_, err := SafeFunc2E(fu)(t1, t2)
		return err
	}
}

// SafeFunc2 returns a function that never panics.
// That function returns the same value as f if f doesn't panic and returns an error if f panics.
func SafeFunc2[T1, T2, U any](f func(T1, T2) U) func(T1, T2) (U, error) {
	fe := func(t1 T1, t2 T2) (U, error) { return f(t1, t2), nil }
	return SafeFunc2E(fe)
}

// SafeFunc2V returns a function that never panics.
// That function returns nil if f doesn't panic and returns an error if f panics.
func SafeFunc2V[T1, T2 any](f func(T1, T2)) func(T1, T2) error {
	fe := func(t1 T1, t2 T2) error { f(t1, t2); return nil }
	return SafeFunc2VE(fe)
}

// RunWithTimeout executes function f with a context constructed from ctx with the addition
// of timeout.
func RunWithTimeout[T any](
	ctx context.Context,
	timeout time.Duration,
	f func(context.Context) (T, error),
) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return f(ctx)
}
