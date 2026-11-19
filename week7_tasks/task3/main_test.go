package main

import (
	"fmt"
	"testing"
)

var sinkS string

func ids(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func TestJoinIDsSlowAllocs(t *testing.T) {
	in := ids(10)
	allocs := testing.AllocsPerRun(1000, func() {
		sinkS = JoinIDsSlow(in)
	})
	fmt.Printf("JoinIDsSlow: %v allocs/op (10 ids)\n", allocs)
}

func TestJoinIDsFastAllocs(t *testing.T) {
	in := ids(10)
	allocs := testing.AllocsPerRun(1000, func() {
		sinkS = JoinIDsFast(in)
	})
	if allocs > 1 {
		t.Errorf("JoinIDsFast: %v allocs, want at most 1", allocs)
	}
}

func BenchmarkJoinIDsSlow(b *testing.B) {
	in := ids(10)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkS = JoinIDsSlow(in)
	}
}

func BenchmarkJoinIDsFast(b *testing.B) {
	in := ids(10)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkS = JoinIDsFast(in)
	}
}

// Measured with `go test -run Allocs -v` and `go test -bench=. -benchmem`
// on Apple M5 Pro, go1.27.1, 10 ids:
//
//   JoinIDsSlow: 18 allocs/op
//   JoinIDsFast: 1  allocs/op
//
//   BenchmarkJoinIDsSlow-15   384.8 ns/op   232 B/op   18 allocs/op
//   BenchmarkJoinIDsFast-15    66.7 ns/op    80 B/op    1 allocs/op
//
// (The handout's own number - 31 allocs/op - is for 10 IDs *and* 1000
// appended numbers combined into one scenario; isolating just the 10-ID
// join here gives a smaller but same-shaped result: roughly 2 allocs per
// loop iteration for Slow - one for the "+=" growth, one for boxing id -
// versus Fast's single allocation overall.)
//
// `go build -gcflags='-m -l' .` on JoinIDsSlow shows:
//   main.go:18:18: ... argument does not escape
//   main.go:18:19: id escapes to heap
// Line 18 is `s += fmt.Sprint(id)`: fmt.Sprint has signature
// `Sprint(a ...any) string`, so every `id` passed to it must be boxed into
// an `any` - that's the escaping line. (The `s += ","` and `s +=
// fmt.Sprint(id)` concatenations themselves also force a fresh allocation
// every iteration, because Go strings are immutable: `s += x` always
// allocates a brand-new backing array sized len(s)+len(x) and copies both
// into it - that's the volume of allocations we see, on top of the one
// escape-to-heap per id.)
//
// Question: JoinIDsFast still makes 1 allocation - which one, and why
// can't it be avoided?
//   It's strings.Builder.String(): a strings.Builder grows its internal
//   []byte buffer (via sb.Grow / WriteByte/Write, which may themselves
//   allocate once if the buffer needs to grow past our Grow(len(ids)*8)
//   guess, but won't here since 10*8=80 bytes is plenty) and then String()
//   must produce an immutable string from that mutable buffer. Go strings
//   and []byte are different, incompatible memory layouts (a string
//   forbids writes), so converting the final buffer into a string requires
//   one allocation+copy at the very end. strings.Builder.String() actually
//   avoids *that* copy with an unsafe cast - but the buffer itself was
//   still heap-allocated once by Grow/append, which is where the 1
//   allocation comes from. You cannot avoid it while still returning a
//   `string` you fully control and own.
//
// Question: why does CollectNoCap(1000) make exactly 12 allocs? See task4.
