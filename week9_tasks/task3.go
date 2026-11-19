package main

import (
	"fmt"
	"sync"
)

// Task 3: Ownership
//
// Run:
//   go run task3.go
//   go run -race task3.go
//
// Part 1a (WRONG: every writer closes the channel) lives in its own file,
// task3_wrong.go, because it is meant to crash the process with a panic —
// see that file's comment for why it can't share a process with the rest
// of this demo. Run it 3 times separately:
//   go run task3_wrong.go

// --- 1b. fix: a closer goroutine waits for all writers, then closes ---

func collect(workers int) <-chan int {
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				results <- w*10 + i
			}
		}(w)
	}
	go func() {
		wg.Wait()      // all writers are finished...
		close(results) // ...so now it's safe to close
	}()
	return results
}

func demoCollect() {
	fmt.Println("=== 1b. fix: collect() with a closer goroutine ===")
	sum := 0
	for v := range collect(3) {
		sum += v
	}
	fmt.Println("sum:", sum, "(correct: 189)")
	fmt.Println(`
wg.Wait() and close(results) must run in their own goroutine because
wg.Wait() blocks until all writers call wg.Done() — and the writers can
only call wg.Done() AFTER they finish sending, which requires something to
be receiving from "results" at the same time (it's unbuffered). If
wg.Wait() ran directly in collect() before "return results", collect
itself would block forever waiting for writers that are in turn blocked
waiting for a receiver that never shows up (collect hasn't returned the
channel to anyone yet) — a deadlock. Running it in a separate goroutine
lets collect return the channel immediately, so the caller's range can
start receiving while wg.Wait() watches in the background.`)
	fmt.Println()
}

// --- 2 & 3. pipeline: generator -> square -> filterEven -> sum ---

// generator owns its output channel: creates it, writes to it, closes it.
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square owns its output channel too.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// filterEven: stage 4, inserted between square and sum. Same shape as
// square: read from in, write to its own out, defer close(out).
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

// sum is a terminal stage: it consumes a channel and returns a plain value,
// no channel of its own to own/close.
func sum(in <-chan int) int {
	total := 0
	for n := range in {
		total += n
	}
	return total
}

func demoPipeline() {
	fmt.Println("=== 2. pipeline: generator -> square -> sum ===")
	total := sum(square(generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)))
	fmt.Println("sum of squares 1..10:", total)

	fmt.Println("=== 3. pipeline with filterEven: generator -> square -> filterEven -> sum ===")
	total2 := sum(filterEven(square(generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))))
	fmt.Println("sum of EVEN squares 1..10:", total2)
	fmt.Println(`
Lines changed in generator/square/sum to add filterEven: zero. Each stage
only knows "read from my <-chan int in, write to my own out, close out
when done" — it has no idea what comes before or after it. filterEven
just slots into the chain with the exact same shape as square. That's the
whole point of the ownership pattern: each stage is independently
composable, so the pipeline's length changes without touching existing
stages.`)
	fmt.Println()
}

func main() {
	fmt.Println(`Part 1a (WRONG, many writers closing) is in task3_wrong.go — it is
meant to crash with "panic: send on closed channel". Run it separately:
  go run task3_wrong.go   (repeat 3x, the sum and panic point vary)`)
	demoCollect()
	demoPipeline()
}
