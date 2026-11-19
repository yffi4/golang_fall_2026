package main

// ⭐ Bonus: does nums escape?

func sum() int {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() { _ = sum() }

// go build -gcflags='-m -l' -o /dev/null bonus.go, go1.27.1:
//   1000:  "make([]int, 1000) does not escape"  - nums never leaves sum
//          (only `total`, an int, is returned), so with a small, constant
//          size the compiler keeps it on the stack.
//   10000: "make([]int, 10000) escapes to heap" - now it moves to the
//          heap, even though nums STILL never leaves sum. The only thing
//          that changed is the size.
//
// Binary search between 1000 and 10000 (same code, only the constant in
// make([]int, N) changes) finds the exact cutoff, measured with go1.27.1:
//   N=8192 ("does not escape")  ->  8192 * 8 bytes = 65536 bytes = 64 KB
//   N=8193 ("escapes to heap")  ->  8193 * 8 bytes = 65544 bytes > 64 KB
// So the cutoff is exactly 64 KB of backing-array size, matching the
// handout's "make() > 64 KB" rule (Part 2) to the byte: int is 8 bytes on
// a 64-bit build, and 8192 is exactly the largest N that keeps the array
// at or under 65536 bytes.
//
// Why does the compiler have such a limit at all?
//   Stack frames are meant to be small and cheap: a function's whole
//   frame has to fit in the current stack segment, and every goroutine
//   stack grows (and gets copied, Part 1) in increments the runtime wants
//   to keep bounded and predictable. If the compiler let frames of
//   arbitrary size live on the stack, one large local array could blow up
//   a frame, make stack growth/copying far more expensive, and in the
//   worst case make a single function call responsible for most of a
//   goroutine's whole 1 GB stack budget. Capping any single stack
//   allocation at 64 KB (and the whole frame at ~128 KB, per the handout)
//   keeps frame sizes - and therefore the cost of growing or copying a
//   stack - bounded, at the cost of pushing genuinely large values to the
//   heap where size isn't a concern.
