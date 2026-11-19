package main

import (
	"fmt"
	"sync"
)

// Task 7: Benchmark — how many goroutines are enough?
//
// Run:
//   go run .
//   go test -bench . -benchmem -cpu 1,2,4,8

const N = 200_000

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// CountSeq counts primes below n, sequentially, no goroutines at all.
func CountSeq(n int) int {
	count := 0
	for i := 2; i < n; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

// CountPar splits [0, n) into `parts` chunks and counts each chunk in its
// own goroutine.
func CountPar(n, parts int) int {
	if parts < 1 {
		parts = 1
	}
	results := make([]int, parts)
	var wg sync.WaitGroup
	chunk := (n + parts - 1) / parts
	for p := 0; p < parts; p++ {
		from := p * chunk
		to := from + chunk
		if to > n {
			to = n
		}
		if from >= to {
			continue
		}
		wg.Add(1)
		go func(p, from, to int) {
			defer wg.Done()
			count := 0
			for i := from; i < to; i++ {
				if isPrime(i) {
					count++
				}
			}
			results[p] = count
		}(p, from, to)
	}
	wg.Wait()

	total := 0
	for _, r := range results {
		total += r
	}
	return total
}

func main() {
	seq := CountSeq(N)
	fmt.Println("CountSeq:", seq)
	for _, parts := range []int{1, 2, 8, 1000, N} {
		par := CountPar(N, parts)
		status := "OK"
		if par != seq {
			status = "MISMATCH"
		}
		fmt.Printf("CountPar(parts=%-7d): %d (%s)\n", parts, par, status)
	}
}
