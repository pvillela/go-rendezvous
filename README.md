Rendezvous
==========

The Rendezvous library supports the safe and convenient execution of asynchronous computations with goroutines and provides facilities for the safe retrieval of the computation results. 

It provides safety in the sense that panics in asynchronous computations are transformed into error results and its methods and functions prevent resource leaks, race conditions, and deadlocks for the channels used to pass data between the parent and child goroutines.

This library uses Golang generics introduced in Go v1.18.  go1.18beta1 or higher must be used with this library.

## Documentation

Run godoc at the root directory to browse the package documentation.
