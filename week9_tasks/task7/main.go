package main

import "fmt"

// Task 7: Benchmark — how expensive is a channel?
//
// Run:
//   go run .
//   go test -v
//   go test -bench . -benchmem -cpu 1,2

const N = 10_000

// passN sends N ints over a channel of the given buffer size and sums them
// on the receiving side.
func passN(size int) int {
	ch := make(chan int, size)
	go func() {
		for i := 0; i < N; i++ {
			ch <- i
		}
		close(ch)
	}()
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

const batchSize = 100

// passBatches sends the same N values, but grouped into []int slices of
// batchSize, so each channel operation carries 100 values instead of 1.
func passBatches() int {
	ch := make(chan []int, 1)
	go func() {
		buf := make([]int, 0, batchSize)
		for i := 0; i < N; i++ {
			buf = append(buf, i)
			if len(buf) == batchSize {
				ch <- buf
				buf = make([]int, 0, batchSize)
			}
		}
		if len(buf) > 0 {
			ch <- buf
		}
		close(ch)
	}()
	sum := 0
	for batch := range ch {
		for _, v := range batch {
			sum += v
		}
	}
	return sum
}

func main() {
	fmt.Println("passN(buf=100):   ", passN(100))
	fmt.Println("passBatches():    ", passBatches())
	fmt.Println("want:             ", N*(N-1)/2)
}
