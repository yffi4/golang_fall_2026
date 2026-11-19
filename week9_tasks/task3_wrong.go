package main

import (
	"fmt"
	"sync"
)

// Task 3, part 1a: WRONG — every writer closes the shared channel.
//
// This is MEANT to panic: a panic in a non-main goroutine cannot be caught
// by a recover() living in a different goroutine, so the whole process
// crashes. Run it three times as separate processes to see the sum and
// the exact panic point vary with scheduling:
//
//   go run task3_wrong.go
//   go run task3_wrong.go
//   go run task3_wrong.go

func main() {
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				results <- w*10 + i
			}
			close(results) // "I'm done" — but the others are not!
		}(w)
	}
	sum := 0
	for v := range results {
		sum += v
	}
	fmt.Println("sum:", sum)
	wg.Wait()
}
