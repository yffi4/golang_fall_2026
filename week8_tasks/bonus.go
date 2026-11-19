package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Bonus question
//
// Run:
//   go run bonus.go

func main() {
	runtime.GOMAXPROCS(1)
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Print(i, " ")
		}()
	}
	wg.Wait()
	fmt.Println()
	fmt.Println(`
It usually prints "5 1 2 3 4" (goroutine 5 lands in P's runnext slot and
runs first — see task4.go, section 4.4), and on a single-core run with
GOMAXPROCS(1) it tends to be very consistent run after run.

Why relying on this order is still a bug, even if it "never changes":
the Go language specification makes NO promise about goroutine scheduling
order — "5 first" is an accident of the current runtime implementation
(the runnext fast-path), not a guaranteed behavior. It can change between
Go versions, between GOMAXPROCS settings, under load from other
goroutines, or on a different CPU/OS — and nothing in the toolchain is
obligated to keep it stable. Code that depends on it is correct only by
coincidence, and will silently misbehave the day the implementation detail
changes underneath it.`)
}
