package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

// Task 4: Tour of Go — Exercise: Equivalent Binary Trees
//
// Run:
//   go run task4.go
//
// We reimplement the tiny tree.Tree/tree.New from the Tour of Go locally
// (no network access needed to fetch golang.org/x/tour/tree) — same shape,
// same behavior: New(k) returns a random BST holding k, 2k, ..., 10k.

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{Value: v}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

// New returns a new, random binary tree holding the values k, 2k, ..., 10k.
func New(k int) *Tree {
	vals := make([]int, 10)
	for i := range vals {
		vals[i] = (i + 1) * k
	}
	rand.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	var t *Tree
	for _, v := range vals {
		t = insert(t, v)
	}
	return t
}

// walk is the recursive worker. It does NOT close ch — close has to happen
// exactly once, after the ENTIRE tree has been sent, and walk calls itself
// once per subtree, so closing inside walk would either close multiple
// times (panic) or close too early (after just one subtree).
func walk(t *Tree, ch chan int) {
	if t == nil {
		return
	}
	walk(t.Left, ch)
	ch <- t.Value
	walk(t.Right, ch)
}

// Walk starts a non-recursive goroutine that drives the recursive walk and
// closes ch only once walk fully returns — i.e. after the whole tree has
// been sent.
func Walk(t *Tree, ch chan int) {
	go func() {
		walk(t, ch)
		close(ch)
	}()
}

// Same reports whether t1 and t2 hold the same values (in-order, so shape
// doesn't matter — only the sorted value sequence does).
func Same(t1, t2 *Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	Walk(t1, ch1)
	Walk(t2, ch2)
	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if ok1 != ok2 {
			return false
		}
		if !ok1 { // both exhausted: every value matched, none left over
			return true
		}
		if v1 != v2 {
			return false
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== Same(New(1), New(1)) — same value set {1..10}, different random shape ===")
	before := runtime.NumGoroutine()
	same := Same(New(1), New(1))
	time.Sleep(50 * time.Millisecond) // let any goroutines settle before measuring
	after := runtime.NumGoroutine()
	fmt.Println("Same:", same, "  NumGoroutine before:", before, "after:", after)

	fmt.Println("\n=== Same(New(1), New(2)) — different value sets ===")
	before2 := runtime.NumGoroutine()
	same2 := Same(New(1), New(2))
	time.Sleep(50 * time.Millisecond)
	after2 := runtime.NumGoroutine()
	fmt.Println("Same:", same2, "  NumGoroutine before:", before2, "after:", after2)

	fmt.Println(`
When the trees are the SAME, Same() drains both channels all the way to
close — both Walk goroutines finish normally, NumGoroutine returns to its
starting point.

When the trees are DIFFERENT, Same() returns early on the first mismatch
(or on the first ok1 != ok2) — but both Walk goroutines are still mid-walk,
blocked trying to send their NEXT value into an unbuffered channel that
Same() has already stopped reading from. Nobody will ever receive from
ch1/ch2 again, so those goroutines leak forever (classic Week 8 leak,
caused here by abandoning an unbuffered channel send). NumGoroutine stays
higher after the "different" case than it was before.`)
}
