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
	"time"

	"github.com/pvillela/go-rendezvous/rdv"
	"golang.org/x/sync/errgroup"
)

func f1(_ context.Context) (int, error) {
	fmt.Println(">>> f1 starting")
	time.Sleep(10 * 1 * time.Millisecond)
	fmt.Println("<<< f1 finishing")
	return 1, nil
}

func f2(_ context.Context) (int, error) {
	fmt.Println(">>> f2 starting")
	time.Sleep(20 * 1 * time.Millisecond)
	fmt.Println("<<< f2 finishing")
	panic("f2 panicked")
}

func f3(_ context.Context) (int, error) {
	fmt.Println(">>> f3 starting")
	time.Sleep(30 * 1 * time.Millisecond)
	fmt.Println("<<< f3 finishing")
	return 0, errors.New("f3 errored-out")
}

func f4(ctx context.Context) (int, error) {
	fmt.Println(">>> f4 starting")
	time.Sleep(40 * 1 * time.Millisecond)
	fmt.Println("<<< f4 finishing")
	if ctx.Err() != nil {
		return 444, errors.New("f4 aborted")
	}
	return 4, nil
}

func f5(ctx context.Context) (int, error) {
	fmt.Println(">>> f5 starting")
	valDur := int64(50 * 1 * time.Millisecond)
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
	time.Sleep(60 * 1 * time.Millisecond)
	fmt.Println("<<< f6 finishing")
	return 6, nil
}

var wastedCancel func()

func ctxTO(millis int64) context.Context {
	timeout := time.Duration(millis) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	wastedCancel = cancel
	return ctx
}

func main() {

	func() {
		fmt.Println("\n*** RdvGo 1")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f1))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvGo 2")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f2))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvGo 3")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f3))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvGo 4")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f4))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvGo 5")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f5))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvGo 6")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f6))
		fmt.Println(rv.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvGo 1 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f1))
		fmt.Println(rv.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvGo 2 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f2))
		fmt.Println(rv.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvGo 3 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f3))
		fmt.Println(rv.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvGo 4 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f3))
		fmt.Println(rv.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvGo 5 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f5))
		fmt.Println(rv.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvGo 6 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		rv := rdv.Go(rdv.CtxApply(ctx, f6))
		fmt.Println(rv.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 12")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f1))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		fmt.Println(rvA.Receive())
		fmt.Println(rvB.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 12 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f1))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		fmt.Println(rvA.ReceiveWatch(ctx))
		fmt.Println(rvB.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 24")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f4))
		fmt.Println(rvA.Receive())
		fmt.Println(rvB.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 24 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f4))
		fmt.Println(rvA.ReceiveWatch(ctx))
		fmt.Println(rvB.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 42")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f4))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		fmt.Println(rvA.Receive())
		fmt.Println(rvB.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 42 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f4))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		fmt.Println(rvA.ReceiveWatch(ctx))
		fmt.Println(rvB.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 25")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f5))
		fmt.Println(rvA.Receive())
		fmt.Println(rvB.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 25 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f5))
		fmt.Println(rvA.ReceiveWatch(ctx))
		fmt.Println(rvB.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 52")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f5))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		fmt.Println(rvA.Receive())
		fmt.Println(rvB.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 52 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f5))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f2))
		fmt.Println(rvA.ReceiveWatch(ctx))
		fmt.Println(rvB.ReceiveWatch(ctx))
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 45")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f4))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f5))
		fmt.Println(rvA.Receive())
		fmt.Println(rvB.Receive())
	}()

	func() {
		fmt.Println("\n*** RdvEgGo 45 ReceiveWatch")
		ctx := ctxTO(48 * 1)
		eg, egCtx := errgroup.WithContext(ctx)
		egF := func() interface{} { return eg } // hack to work around go2go bug
		rvA := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f4))
		rvB := rdv.XGoEg(egF, rdv.CtxApply(egCtx, f5))
		fmt.Println(rvA.ReceiveWatch(ctx))
		fmt.Println(rvB.ReceiveWatch(ctx))
	}()
}
