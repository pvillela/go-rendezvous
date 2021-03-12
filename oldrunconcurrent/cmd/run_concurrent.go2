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

	"github.com/pvillela/go-trygo2/oldasync"
)

func f1(context.Context) (int, error) {
	fmt.Println(">>> f1 starting")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("<<< f1 finishing")
	return 1, nil
}

func f2(context.Context) (int, error) {
	fmt.Println(">>> f2 starting")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("<<< f2 finishing")
	panic("f2 panicked")
}

func f3(context.Context) (int, error) {
	fmt.Println(">>> f3 starting")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("<<< f3 finishing")
	return 0, errors.New("f3 errored-out")
}

func f4(ctx context.Context) (int, error) {
	fmt.Println(">>> f4 starting")
	time.Sleep(400 * time.Millisecond)
	fmt.Println("<<< f4 finishing")
	if ctx.Err() != nil {
		return 999, errors.New("aborted")
	}
	return 4, nil
}

func main() {
	ctx := context.Background()

	fmt.Println("*** RunConcurrentsWg 1234")
	resultsWg, err := oldasync.RunConcurrentsWg(ctx, f1, f2, f3, f4)
	fmt.Println(resultsWg, err)

	fmt.Println("*** RunConcurrentsEg 1234")
	resultsEg, err := oldasync.RunConcurrentsEg(ctx, f1, f2, f3, f4)
	fmt.Println(resultsEg, err)

	fmt.Println("*** RunConcurrentsEg 124")
	resultsEg, err = oldasync.RunConcurrentsEg(ctx, f1, f2, f4)
	fmt.Println(resultsEg, err)

	fmt.Println("*** RunConcurrentsEg 134")
	resultsEg, err = oldasync.RunConcurrentsEg(ctx, f1, f3, f4)
	fmt.Println(resultsEg, err)

	fmt.Println("*** RunConcurrentsEg 14")
	resultsEg, err = oldasync.RunConcurrentsEg(ctx, f1, f4)
	fmt.Println(resultsEg, err)

	fmt.Println("*** RunConcurrent2Wg 12")
	result2Wg, err := oldasync.RunConcurrent2Wg(ctx, f1, f2)
	fmt.Println(result2Wg, err)

	fmt.Println("*** RunConcurrent2Wg 13")
	result2Wg, err = oldasync.RunConcurrent2Wg(ctx, f1, f3)
	fmt.Println(result2Wg, err)

	fmt.Println("*** RunConcurrent2Wg 14")
	result2Wg, err = oldasync.RunConcurrent2Wg(ctx, f1, f4)
	fmt.Println(result2Wg, err)

	fmt.Println("*** RunConcurrent2Eg 12")
	result2Eg, err := oldasync.RunConcurrent2Eg(ctx, f1, f2)
	fmt.Println(result2Eg, err)

	fmt.Println("*** RunConcurrent2Eg 13")
	result2Eg, err = oldasync.RunConcurrent2Eg(ctx, f1, f3)
	fmt.Println(result2Eg, err)

	fmt.Println("*** RunConcurrent2Eg 14")
	result2Eg, err = oldasync.RunConcurrent2Eg(ctx, f1, f4)
	fmt.Println(result2Eg, err)

	fmt.Println("*** RunConcurrent2Eg 24")
	result2Eg, err = oldasync.RunConcurrent2Eg(ctx, f2, f4)
	fmt.Println(result2Eg, err)

	fmt.Println("*** RunConcurrent2Eg 43")
	result2Eg, err = oldasync.RunConcurrent2Eg(ctx, f4, f3)
	fmt.Println(result2Eg, err)
}
