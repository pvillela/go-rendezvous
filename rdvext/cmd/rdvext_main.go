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
	"github.com/pvillela/go-rendezvous"
	"time"

	"github.com/pvillela/go-rendezvous/rdvext"
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
		rv := rendezvous.Go(rendezvous.CtxApply(toCtx(49), f1))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Async 2")
		rv := rendezvous.Go(rendezvous.CtxApply(toCtx(49), f2))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Async 3")
		rv := rendezvous.Go(rendezvous.CtxApply(toCtx(49), f3))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Go2Eg 12")
		rv := rdvext.Go2Eg(toCtx(49), f1, f2)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Go2Eg 14")
		rv := rdvext.Go2Eg(toCtx(49), f1, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Go2Eg 15")
		rv := rdvext.Go2Eg(toCtx(49), f1, f5)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSliceEg 1234")
		rv := rdvext.GoSliceEg(toCtx(49), f1, f2, f3, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSliceEg 12345")
		rv := rdvext.GoSliceEg(toCtx(49), f1, f2, f3, f4, f5)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSliceEg 123456")
		rv := rdvext.GoSliceEg(toCtx(49), f1, f2, f3, f4, f5, f6)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSliceEg 14")
		rv := rdvext.GoSliceEg(toCtx(49), f1, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSliceEg 24")
		rv := rdvext.GoSliceEg(toCtx(49), f2, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSliceEg 145")
		rv := rdvext.GoSliceEg(toCtx(49), f1, f4, f5)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Go2 12")
		rv := rdvext.Go2(toCtx(49), f1, f2)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Go2 14")
		rv := rdvext.Go2(toCtx(49), f1, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** Go2 15")
		rv := rdvext.Go2(toCtx(49), f1, f5)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSlice 1234")
		rv := rdvext.GoSlice(toCtx(49), f1, f2, f3, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSlice 12345")
		rv := rdvext.GoSlice(toCtx(49), f1, f2, f3, f4, f5)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSlice 14")
		rv := rdvext.GoSlice(toCtx(49), f1, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSlice 24")
		rv := rdvext.GoSlice(toCtx(49), f2, f4)
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** GoSlice 145")
		rv := rdvext.GoSlice(toCtx(49), f1, f4, f5)
		fmt.Println(rv.Receive())
	}()
}
