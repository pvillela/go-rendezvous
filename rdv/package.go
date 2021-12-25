/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

// Package rdv supports the safe and convenient execution of asynchronous computations with
// goroutines and provides facilities for the safe retrieval of the computation results.
// It provides safety in the sense that panics in asynchronous computations are transformed
// into error results and its methods and functions prevent resource leaks, race conditions,
// and deadlocks for the channels used to pass data between the parent and child goroutines.
package rdv
