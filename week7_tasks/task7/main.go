package main

import (
	"bytes"
	"fmt"
	"sync"
)

// Task 7: Reuse buffers with sync.Pool

var bufPool = sync.Pool{
	New: func() any { b := new(bytes.Buffer); b.Grow(4096); return b },
}

// RenderNew allocates a fresh 4 KB buffer every call.
func RenderNew(ids []int) int {
	buf := new(bytes.Buffer)
	buf.Grow(4096)
	for _, id := range ids {
		fmt.Fprintf(buf, "id=%d\n", id)
	}
	return buf.Len()
}

// RenderPool reuses a buffer from the pool instead of allocating a new one.
func RenderPool(ids []int) int {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	for _, id := range ids {
		fmt.Fprintf(buf, "id=%d\n", id)
	}
	n := buf.Len()
	bufPool.Put(buf)
	return n
}

func main() {
	ids := make([]int, 200)
	for i := range ids {
		ids[i] = i
	}
	fmt.Println("RenderNew: ", RenderNew(ids), "bytes")
	fmt.Println("RenderPool:", RenderPool(ids), "bytes")
	fmt.Println("run: go test -bench=. -benchmem")
}
