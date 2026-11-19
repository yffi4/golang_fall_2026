package main

import "testing"

var sinkN int

func makeIDs() []int {
	ids := make([]int, 200)
	for i := range ids {
		ids[i] = i
	}
	return ids
}

func BenchmarkRenderNew(b *testing.B) {
	ids := makeIDs()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkN = RenderNew(ids)
	}
}

func BenchmarkRenderPool(b *testing.B) {
	ids := makeIDs()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkN = RenderPool(ids)
	}
}

// Measured with `go test -bench=. -benchmem` on Apple M5 Pro, go1.27.1,
// 200 ids:
//
//   BenchmarkRenderNew-15    4946 ns/op   4146 B/op   2 allocs/op
//   BenchmarkRenderPool-15   4631 ns/op      0 B/op   0 allocs/op
//
// RenderNew allocates a new 4 KB *bytes.Buffer (2 allocs: the Buffer
// struct itself, plus its backing []byte from Grow) on every single call,
// even though the buffer is thrown away a few lines later. RenderPool
// reuses an existing buffer from the pool, so in the steady state (once
// the pool has been warmed up) it makes 0 allocations per call - Get
// returns an already-allocated buffer instead of the runtime handing out
// fresh memory.
//
// Why must you call buf.Reset() after Get?
//   sync.Pool doesn't know or care what state an object is in when you
//   Put it back - it just hands the same object to the next Get, byte-
//   for-byte as it was left. If the previous caller wrote "id=199\n" into
//   the buffer and didn't reset it, the next caller's Fprintf calls would
//   just keep appending after that leftover data instead of starting from
//   an empty buffer - you'd get garbage mixed into your "fresh" output.
//   Reset() truncates the buffer's length back to 0 while keeping its
//   already-allocated backing array, which is exactly the point: get the
//   capacity without the old contents.
//
// What could go wrong if you keep using buf after Put?
//   The moment you call Put, you no longer own that buffer - any other
//   goroutine's Get can receive the exact same pointer at any time
//   afterwards. If you keep writing to or reading from buf after Put, you
//   now have two goroutines mutating the same *bytes.Buffer concurrently
//   with no synchronization: a classic data race (corrupted/garbled
//   output, or a crash if the slice header itself gets torn mid-write).
//   The fix is discipline, not a lock: never touch a pooled object again
//   after Put, the same way you never use a slice after append has
//   possibly handed its backing array to someone else.
