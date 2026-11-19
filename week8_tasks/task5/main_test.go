package main

import (
	"runtime"
	"testing"
	"time"
)

// assertNoLeak gives pending goroutines a short grace period to finish
// (GC and the scheduler need a moment), then fails if the count went up.
func assertNoLeak(t *testing.T, before int) {
	t.Helper()
	deadline := time.Now().Add(200 * time.Millisecond)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("goroutine leak: had %d goroutines before, still %d after waiting", before, runtime.NumGoroutine())
}

func TestFetchFixedDoesNotLeak(t *testing.T) {
	before := runtime.NumGoroutine()
	for i := 0; i < 50; i++ {
		_, _ = fetchFixed(10 * time.Millisecond)
	}
	assertNoLeak(t, before)
}

// TestFetchLeakyLeaks documents the bug: it is expected to fail if run
// against fetchLeaky, proving the leak exists. Skipped by default so the
// normal test run stays green; run with -run TestFetchLeakyLeaks to see it fail.
func TestFetchLeakyLeaks(t *testing.T) {
	t.Skip("documents the leak in fetchLeaky; remove t.Skip to see it fail")
	before := runtime.NumGoroutine()
	for i := 0; i < 50; i++ {
		_, _ = fetchLeaky(10 * time.Millisecond)
	}
	assertNoLeak(t, before)
}
