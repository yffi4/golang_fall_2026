package main

import "fmt"

// Task 4: Benchmark - value vs pointer, with vs without capacity

type Point struct{ X, Y int }

//go:noinline
func NewPointValue(x, y int) Point { return Point{x, y} }

//go:noinline
func NewPointPointer(x, y int) *Point { return &Point{x, y} }

// CollectNoCap grows the slice from a zero-value nil slice: append has to
// guess the capacity and reallocate+copy every time it runs out of room.
func CollectNoCap(n int) []int {
	var out []int
	for i := 0; i < n; i++ {
		out = append(out, i)
	}
	return out
}

// CollectWithCap pre-sizes the backing array once, so append never
// reallocates.
func CollectWithCap(n int) []int {
	out := make([]int, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, i)
	}
	return out
}

func main() {
	fmt.Println(NewPointValue(1, 2))
	fmt.Println(*NewPointPointer(1, 2))
	fmt.Println(len(CollectNoCap(1000)), len(CollectWithCap(1000)))
	fmt.Println("run: go test -bench=. -benchmem")
}
