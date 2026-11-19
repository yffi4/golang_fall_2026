package main

import (
	"fmt"
	"os"
)

// Task 1: Explore the stack

// depthSmall is the depth() from the lecture: pad is [64]byte.
func depthSmall(n int) int {
	var pad [64]byte
	pad[n%len(pad)] = byte(n)
	if n == 0 {
		return 0
	}
	return depthSmall(n-1) + 1 + int(pad[(n+1)%len(pad)]) // reads a 0, but the compiler can't know that
}

// depthBig is the same function with a 16x bigger pad.
func depthBig(n int) int {
	var pad [1024]byte
	pad[n%len(pad)] = byte(n)
	if n == 0 {
		return 0
	}
	return depthBig(n-1) + 1 + int(pad[(n+1)%len(pad)])
}

// depthNoRead is depthSmall with the line that *reads* pad removed - pad is
// now write-only.
func depthNoRead(n int) int {
	var pad [64]byte
	pad[n%len(pad)] = byte(n)
	if n == 0 {
		return 0
	}
	return depthNoRead(n-1) + 1
}

func forever(n int) int {
	return forever(n+1) + 1
}

// tryForever calls forever wrapped in recover, to check whether a stack
// overflow is a recoverable panic.
func tryForever() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	fmt.Println(forever(0))
	fmt.Println("unreachable")
}

func main() {
	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "big":
		// depthBig alone, isolated: its frame (1056 B) times 1,000,000
		// calls is just over the 1 GB stack limit - see note below.
		fmt.Println("depth, pad [1024]byte: ", depthBig(1_000_000))
	case "forever":
		fmt.Println("--- calling forever(0), wrapped in recover ---")
		tryForever()
	default:
		fmt.Println("depth, pad [64]byte:   ", depthSmall(1_000_000))
		fmt.Println("depth, pad never read: ", depthNoRead(1_000_000))
		fmt.Println(`(run "go run task1.go big" to see depthBig overflow)`)
		fmt.Println(`(run "go run task1.go forever" to see the stack overflow)`)
	}
}

// 1. Frame sizes, measured with:
//      go build -gcflags=-S -o /dev/null task1.go  →  grep 'TEXT main\.depth.*SB'
//    on Apple M5 Pro, go1.27.1:
//      main.depthSmall  frame size:   96 bytes (pad [64]byte  + saved regs/args)
//      main.depthBig    frame size: 1056 bytes (pad [1024]byte + saved regs/args)
//    depthSmall(1_000_000) works fine: ~96 bytes x 1,000,000 ≈ 96 MB, well
//    under the 1 GB goroutine stack limit; the runtime grows and copies the
//    stack a few times along the way (see the note below) and it just
//    works.
//    depthBig(1_000_000) does NOT work ("go run task1.go big"):
//      runtime: goroutine stack exceeds 1000000000-byte limit
//      fatal error: stack overflow
//    1056 bytes x 1,000,000 ≈ 1.056 GB, which is just *over* the 1 GB
//    limit - so increasing pad to [1024]byte is enough to turn a program
//    that worked into a stack overflow, with the same n and the same
//    recursion shape. This is exactly the point of the exercise: stack
//    usage is "frame size x recursion depth", and both factors matter.
//
// 2. depthNoRead: removing the line that *reads* pad shrinks the frame a
//    lot, from 96 to 32 bytes. pad[n%len(pad)] = byte(n) (the write) is
//    still a visible side effect, so the compiler can't drop the array,
//    but once nothing ever reads pad back, the compiler no longer needs
//    to keep it addressable/escaping in the way the read required, and
//    the stack slot layout shrinks. Only removing *both* the read and the
//    write would let the compiler delete pad entirely.
//
// 3. forever(0) crashes with:
//      runtime: goroutine stack exceeds 1000000000-byte limit
//      fatal error: stack overflow
//    This is a *fatal error*, not a panic - recover() never runs (you can
//    see this by running `go run task1.go forever`: the "recovered: ..."
//    line never prints, the process just exits with status 2). The Go
//    runtime treats a handful of conditions (stack overflow, out of
//    memory, concurrent map misuse, etc.) as unrecoverable: they abort the
//    whole program before any deferred function gets a chance to run,
//    because by that point there may be no stack left to run them on.
//
// Participation Task: why the goroutine stack is copied when it grows, and
// why that's safe even though pointers may point into the stack.
//   A goroutine starts with a tiny (2 KB) stack. When a call would overflow
//   it, the runtime allocates a new, bigger stack (typically 2x) and copies
//   every byte of the old stack into it - frames, saved registers, local
//   variables, everything. This is only possible because the Go compiler
//   knows the type of every stack frame (it's not opaque memory the way C's
//   stack is): the runtime's stack-copying code walks the frames using that
//   type information and rewrites every pointer that pointed *into the old
//   stack* so it points at the matching address in the new one. Pointers
//   that point into the heap are untouched (the heap doesn't move). Because
//   the rewrite is exhaustive and precise, no pointer is ever left dangling,
//   which is exactly why Go can safely copy/move stacks while C (which has
//   no such metadata) cannot.
