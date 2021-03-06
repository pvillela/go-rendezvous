/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/pvillela/go-rendezvous/obsolete/async"
	"time"

	"github.com/pvillela/go-rendezvous/util"
)

func f1(_ context.Context) (int, error) {
	fmt.Println(">>> f1 starting")
	time.Sleep(10 * time.Millisecond)
	fmt.Println("<<< f1 finishing")
	return 1, nil
}

func f2(_ context.Context) (int, error) {
	fmt.Println(">>> f2 starting")
	time.Sleep(20 * time.Millisecond)
	fmt.Println("<<< f2 finishing")
	panic("f2 panicked")
}

func f3(_ context.Context) (int, error) {
	fmt.Println(">>> f3 starting")
	time.Sleep(30 * time.Millisecond)
	fmt.Println("<<< f3 finishing")
	return 0, errors.New("f3 errored-out")
}

func f4(ctx context.Context) (int, error) {
	fmt.Println(">>> f4 starting")
	time.Sleep(40 * time.Millisecond)
	fmt.Println("<<< f4 finishing")
	if ctx.Err() != nil {
		return 444, errors.New("f4 aborted")
	}
	return 4, nil
}

func f5(ctx context.Context) (int, error) {
	fmt.Println(">>> f5 starting")
	valDur := int64(50 * time.Millisecond)
	shorDur := time.Duration(valDur / 10)
	for i := 0; i < 10; i++ {
		if ctx.Err() != nil {
			fmt.Println("<<< f5 aborting")
			return 555, errors.New(fmt.Sprintf("f5 aborted on iteration %v", i))
		}
		time.Sleep(shorDur)
	}
	fmt.Println("<<< f5 finishing")
	return 5, nil
}

func f6(_ context.Context) (int, error) {
	fmt.Println(">>> f6 starting")
	time.Sleep(60 * time.Millisecond)
	fmt.Println("<<< f6 finishing")
	return 6, nil
}

var wastedCancel func()

func toCtx(millis int64) context.Context {
	timeout := time.Duration(millis) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	wastedCancel = cancel
	return ctx
}

func main() {
	func() {
		fmt.Println("\n*** SafeFunc1E(f1)")
		f := util.SafeFunc1E(f1)
		fmt.Println(f(toCtx(49)))
	}()

	func() {
		fmt.Println("\n*** SafeFunc1E(f2)")
		f := util.SafeFunc1E(f2)
		fmt.Println(f(toCtx(49)))
	}()

	func() {
		fmt.Println("\n*** Async 1")
		promise := async.Async(toCtx(49), f1)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async 2")
		promise := async.Async(toCtx(49), f2)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async 3")
		promise := async.Async(toCtx(49), f3)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async2Eg 12")
		promise := async.Async2Eg(toCtx(49), f1, f2)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async2Eg 14")
		promise := async.Async2Eg(toCtx(49), f1, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async2Eg 15")
		promise := async.Async2Eg(toCtx(49), f1, f5)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsEg 1234")
		promise := async.AsyncsEg(toCtx(49), f1, f2, f3, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsEg 12345")
		promise := async.AsyncsEg(toCtx(49), f1, f2, f3, f4, f5)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsEg 123456")
		promise := async.AsyncsEg(toCtx(49), f1, f2, f3, f4, f5, f6)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsEg 14")
		promise := async.AsyncsEg(toCtx(49), f1, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsEg 24")
		promise := async.AsyncsEg(toCtx(49), f2, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsEg 145")
		promise := async.AsyncsEg(toCtx(49), f1, f4, f5)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async2Wg 12")
		promise := async.Async2Wg(toCtx(49), f1, f2)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async2Wg 14")
		promise := async.Async2Wg(toCtx(49), f1, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** Async2Wg 15")
		promise := async.Async2Wg(toCtx(49), f1, f5)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsWg 1234")
		promise := async.AsyncsWg(toCtx(49), f1, f2, f3, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsWg 12345")
		promise := async.AsyncsWg(toCtx(49), f1, f2, f3, f4, f5)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsWg 14")
		promise := async.AsyncsWg(toCtx(49), f1, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsWg 24")
		promise := async.AsyncsWg(toCtx(49), f2, f4)
		fmt.Println(promise.Await())
	}()

	func() {
		fmt.Println("\n*** AsyncsWg 145")
		promise := async.AsyncsWg(toCtx(49), f1, f4, f5)
		fmt.Println(promise.Await())
	}()
}
