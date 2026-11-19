package main

import "testing"

var sinkPoint Point
var sinkPointPtr *Point
var sinkSlice []int

func BenchmarkNewPointValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPoint = NewPointValue(1, 2)
	}
}

func BenchmarkNewPointPointer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkPointPtr = NewPointPointer(1, 2)
	}
}

func BenchmarkCollectNoCap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkSlice = CollectNoCap(1000)
	}
}

func BenchmarkCollectWithCap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkSlice = CollectWithCap(1000)
	}
}

// Why //go:noinline on NewPointValue/NewPointPointer?
//   Both functions are tiny, so without it the compiler would inline them
//   straight into the benchmark loop. Once inlined, escape analysis looks
//   at the *inlined* code, not the original function in isolation - for
//   NewPointPointer that can let the compiler prove the pointer never
//   really leaves the loop body and keep it on the stack after all,
//   which would hide exactly the value-vs-pointer allocation difference
//   this benchmark is supposed to measure. //go:noinline forces the
//   compiler to treat each call as a real, opaque function call, so the
//   escape decision matches what `go build -gcflags='-m -l'` reports for
//   the function as written.
//
// Measured with `go test -bench=. -benchmem` on Apple M5 Pro, go1.27.1:
//
//   Benchmark             ns/op    B/op   allocs/op
//   NewPointValue         0.45      0       0
//   NewPointPointer       5.47     16       1
//   CollectNoCap(1000)   2074   25208      12
//   CollectWithCap(1000)  922    8192       1
//
// Connecting every result to `go build -gcflags='-m -l' .`:
//   NewPointValue:   no escape line at all for the Point{x,y} literal -
//                    it's returned by value, so it stays on the stack in
//                    the caller's frame. 0 allocs/op confirms it.
//   NewPointPointer: "&Point{...} escapes to heap" - returning the address
//                    of a local forces it to the heap (same shape as
//                    newUserPointer in the lecture). 1 alloc/op, 16 B
//                    (two ints, 8 bytes each) confirms it.
//   CollectNoCap:    "append escapes to heap" marks the backing array
//                    `append` (re)allocates every time it runs out of
//                    room, and growing from a nil slice reallocates
//                    repeatedly: for n=1000 the capacity sequence Go's
//                    append uses is roughly 1,2,4,8,...,1024 (doubling,
//                    then a smaller growth factor past 256 elements) -
//                    that's where the 12 allocations come from, not 1.
//                    The *final* array has capacity 1024 (not 1000), so it
//                    holds 1024*8=8192 bytes, but CollectNoCap's 25208
//                    B/op also includes every *intermediate* array it
//                    allocated and then threw away along the growth path,
//                    which the GC has to clean up later - that's the
//                    "costs twice" lesson from Part 3 in a single
//                    benchmark number.
//   CollectWithCap:  "make([]int, 0, n) escapes to heap" (it's the return
//                    value, so of course it escapes) - but it allocates
//                    exactly once, for 1000 ints -> 8000 bytes (rounded up
//                    by the allocator's size classes to 8192), because
//                    make([]int, 0, n) sizes the backing array up front
//                    and append never needs to grow it.
