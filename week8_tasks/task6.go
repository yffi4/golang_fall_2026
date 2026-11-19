package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// Task 6: A leaking "server"
//
// Run:
//   go run task6.go

type Report struct {
	Data [64 * 1024]byte // 64 KB, so the leak is easy to see in the heap
}

func buildReport() *Report {
	time.Sleep(20 * time.Millisecond)
	return &Report{}
}

// handleLeaky: an unbuffered channel send blocks forever once the timeout
// has already made handle return — the goroutine, and the 64 KB report it
// holds onto, both leak.
func handleLeaky(timeout time.Duration) error {
	ch := make(chan *Report)
	go func() { ch <- buildReport() }()
	select {
	case r := <-ch:
		_ = r
		return nil
	case <-time.After(timeout):
		return errors.New("timeout")
	}
}

// handleBuffered: step (a) — a buffered channel lets the goroutine finish
// and exit even if nobody is listening anymore. This stops the GOROUTINE
// leak, but buildReport still always runs to completion and still
// allocates its 64 KB, even though the caller gave up on it.
func handleBuffered(timeout time.Duration) error {
	ch := make(chan *Report, 1)
	go func() { ch <- buildReport() }()
	select {
	case r := <-ch:
		_ = r
		return nil
	case <-time.After(timeout):
		return errors.New("timeout")
	}
}

func buildReportCtx(ctx context.Context) *Report {
	select {
	case <-time.After(20 * time.Millisecond):
		return &Report{}
	case <-ctx.Done():
		return nil // give up early: never allocates the report
	}
}

// handleCtxTimeout: step (b) — pass a context.WithTimeout into buildReport
// so that it also stops DOING WORK once nobody needs the result, instead of
// just quietly finishing work that gets thrown away.
func handleCtxTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ch := make(chan *Report, 1)
	go func() { ch <- buildReportCtx(ctx) }()
	select {
	case r := <-ch:
		if r == nil {
			return errors.New("timeout")
		}
		return nil
	case <-ctx.Done():
		return errors.New("timeout")
	}
}

func heapMB() float64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.HeapAlloc) / 1024 / 1024
}

func runRounds(label string, handle func(time.Duration) error) {
	fmt.Println("===", label, "===")
	fmt.Printf("start            goroutines=%-6d heap=%.1f MB\n", runtime.NumGoroutine(), heapMB())
	for round := 1; round <= 3; round++ {
		for i := 0; i < 1000; i++ {
			_ = handle(5 * time.Millisecond) // buildReport takes 20ms > 5ms timeout: (almost) always "times out"
		}
		time.Sleep(30 * time.Millisecond) // give goroutines that CAN finish a chance to do so before we measure
		fmt.Printf("after round %d    goroutines=%-6d heap=%.1f MB\n", round, runtime.NumGoroutine(), heapMB())
	}
	fmt.Println()
}

func main() {
	runRounds("leaky: unbuffered channel", handleLeaky)
	runRounds("fix (a): buffered channel only", handleBuffered)
	runRounds("fix (b): buffered channel + context.WithTimeout", handleCtxTimeout)

	fmt.Println(`Why (b) is better than (a) alone:
(a) stops the GOROUTINE from leaking, but buildReport still runs to
completion every single time, burning CPU and allocating a 64 KB report
that is immediately thrown away. (b) additionally gives buildReport a
context, so it can notice ctx.Done() and stop WORKING as soon as nobody
needs the result anymore — no wasted allocation, no wasted CPU time.`)
}
