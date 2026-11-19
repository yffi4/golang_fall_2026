package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// Task 5: Find the leak

func readFile(i int) []byte {
	data := make([]byte, 1<<20) // pretend we read a 1 MB file
	copy(data, fmt.Sprintf("HEADER-%04d", i))
	return data
}

// keepHeadersSlice keeps a subslice of every 1 MB "file" - len 11, but the
// underlying array (cap) is still the full 1 MB.
func keepHeadersSlice(n int) [][]byte {
	var headers [][]byte
	for i := 0; i < n; i++ {
		data := readFile(i)
		headers = append(headers, data[:11])
	}
	return headers
}

// keepHeadersCopy clones just the 11 bytes it needs into a fresh, small
// allocation, letting the 1 MB array become garbage immediately.
func keepHeadersCopy(n int) [][]byte {
	var headers [][]byte
	for i := 0; i < n; i++ {
		data := readFile(i)
		headers = append(headers, bytes.Clone(data[:11]))
	}
	return headers
}

func bigString() string {
	return strings.Repeat("x", 1<<20) // 1 MB string
}

func keepSubstrings(n int) []string {
	var subs []string
	for i := 0; i < n; i++ {
		big := bigString()
		subs = append(subs, big[:5])
	}
	return subs
}

func keepClones(n int) []string {
	var subs []string
	for i := 0; i < n; i++ {
		big := bigString()
		subs = append(subs, strings.Clone(big[:5]))
	}
	return subs
}

func liveHeapMB() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc >> 20
}

func main() {
	const n = 100

	subslices := keepHeadersSlice(n)
	fmt.Println("subslices: ", liveHeapMB(), "MB") // ~100 MB: each kept []byte still points into a 1 MB array
	runtime.KeepAlive(subslices)

	copies := keepHeadersCopy(n)
	fmt.Println("copies:    ", liveHeapMB(), "MB") // ~0 MB: the 1 MB arrays were freed, only 11-byte clones remain
	runtime.KeepAlive(copies)

	substrings := keepSubstrings(n)
	fmt.Println("substrings:", liveHeapMB(), "MB") // ~100 MB: same problem, strings share a backing array too
	runtime.KeepAlive(substrings)

	clones := keepClones(n)
	fmt.Println("clones:    ", liveHeapMB(), "MB") // ~0 MB
	runtime.KeepAlive(clones)
}

// Question: a slice data[:11] has len 11 - why does it keep the whole 1 MB
// array alive?
//   A slice header is {pointer, len, cap} (Week 2). data[:11] gets a NEW
//   header with len=11, but its pointer still points at the SAME backing
//   array as data, and that array's size never shrinks just because the
//   slice's own len/cap did. The garbage collector doesn't know "this
//   slice only needs its first 11 bytes" - it only knows "this pointer
//   points somewhere inside this 1 MB array", so the whole array stays
//   reachable (and therefore alive) for as long as any slice points into
//   it, no matter how small that slice's len is. Exactly the same thing
//   happens with strings: big[:5] shares big's backing byte array, so
//   keeping big[:5] alive keeps all of big's memory alive too.
//
// Why does the measurement show 0 MB for the "leaking" version if you
// forget runtime.KeepAlive?
//   If a variable is never read again after the point where you measure,
//   the compiler/runtime can prove it's unreachable from that point
//   onward and the GC is free to collect it - even if the variable is
//   still "in scope" textually. runtime.GC() right before
//   ReadMemStats would then find nothing to keep around and report a
//   deceptively small heap. runtime.KeepAlive(x) after the measurement
//   tells the runtime "x must still be considered reachable up to this
//   point", forcing the GC to actually keep (and therefore count) the
//   memory we're trying to measure.
