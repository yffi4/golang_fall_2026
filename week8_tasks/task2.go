package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Task 2: First goroutines
//
// Run:
//   go run task2.go
//   go run -race task2.go   (to see the data race in part 2)

// --- 2.1: fix "say" with a sync.WaitGroup -----------------------------------

func say(msg string) {
	for i := 0; i < 3; i++ {
		fmt.Println(msg, i)
	}
}

func demoSay() {
	fmt.Println("--- 2.1: say, fixed with WaitGroup ---")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		say("goroutine")
	}()
	say("main")
	wg.Wait() // without this, main could return before the goroutine prints anything
}

// --- 2.2: two classic mistakes, fixed two ways ------------------------------

// racyCounter reproduces both classic bugs on purpose, so `go run -race`
// can show the WARNING. Only call it under -race for demonstration; its
// printed result is intentionally unreliable.
func racyCounter() int {
	var wg sync.WaitGroup
	counter := 0
	for i := 0; i < 1000; i++ {
		wg.Add(1) // still "wrong-ish" in spirit: Add happens outside the goroutine, but see fixedCounterAtomic for the real fix target
		go func() {
			defer wg.Done()
			counter++ // DATA RACE: read-add-write is not atomic
		}()
	}
	wg.Wait()
	return counter
}

// fixedCounterAtomic: fix #1 — use sync/atomic so the increment is one
// indivisible operation.
func fixedCounterAtomic() int64 {
	var counter atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	return counter.Load()
}

// fixedCounterPerGoroutine: fix #2 — give every goroutine its own slice
// element, so nobody shares mutable state and there is nothing to race on.
func fixedCounterPerGoroutine() int {
	const n = 1000
	hits := make([]int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			hits[i] = 1
		}()
	}
	wg.Wait()
	total := 0
	for _, h := range hits {
		total += h
	}
	return total
}

func demoClassicMistakes() {
	fmt.Println("\n--- 2.2: two classic mistakes ---")
	fmt.Println("racy counter (expected 1000, run with -race to see the warning):", racyCounter())
	fmt.Println("fix 1, atomic.Int64:", fixedCounterAtomic())
	fmt.Println("fix 2, one slice element per goroutine:", fixedCounterPerGoroutine())
}

// --- 2.3: panic in a goroutine, recovered without killing the program ------

func demoPanicRecover() {
	fmt.Println("\n--- 2.3: panic + recover inside the goroutine ---")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// The defer with recover() must be INSIDE this goroutine: recover
		// only catches a panic in the goroutine that is currently
		// unwinding. main's own recover would not see this panic at all.
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered inside goroutine:", r)
			}
		}()
		var m map[string]int
		m["x"] = 1 // panics: assignment to entry in nil map
	}()
	wg.Wait()
	fmt.Println("main is still alive")
}

func main() {
	demoSay()
	demoClassicMistakes()
	demoPanicRecover()
}
