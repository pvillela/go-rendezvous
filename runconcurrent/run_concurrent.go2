package runconcurrent

import "sync"

type ResultWithError struct {
	Result int
	Error  interface{}
}

func safeGo(pWg *sync.WaitGroup, f func() int, index int, pResults *[]ResultWithError) {
	pWg.Add(1)
	safeF := func() {
		defer func() {
			err := recover()
			if err != nil {
				(*pResults)[index] = ResultWithError{0, err}
			}
			pWg.Done()
		}()
		(*pResults)[index] = ResultWithError{f(), nil}
	}
	go safeF()
}

func RunConcurrent(funcs ...func() int) []ResultWithError {
	results := make([]ResultWithError, len(funcs))
	var wg sync.WaitGroup
	for index, f := range funcs {
		safeGo(&wg, f, index, &results)
	}
	wg.Wait()
	return results
}
