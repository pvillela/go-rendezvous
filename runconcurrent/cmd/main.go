package main

import (
	"fmt"
	"github.com/pvillela/go-tryout/runconcurrent"
	"time"
)

func f1() int {
	fmt.Println(">>> f1 starting")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("<<< f1 finishing")
	return 1
}

func f2() int {
	fmt.Println(">>> f2 starting")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("<<< f2 finishing")
	panic("f2 panicked")
}

func f3() int {
	fmt.Println(">>> f3 starting")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("<<< f3 finishing")
	return 3
}

func main() {
	results := runconcurrent.RunConcurrent(f1, f2, f3)
	fmt.Println(results)
}
