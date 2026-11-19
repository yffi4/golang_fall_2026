package main

import (
	"sync"
	"testing"
)

func noop() {}

// BenchmarkGoroutineStartStop measures the cost of starting a goroutine and
// waiting for it to finish.
func BenchmarkGoroutineStartStop(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			noop()
		}()
		wg.Wait()
	}
}

// BenchmarkFuncCall measures the cost of calling the same function directly,
// with no goroutine involved — the baseline to compare against.
func BenchmarkFuncCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		noop()
	}
}
