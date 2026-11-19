package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Task 1: Concurrency vs parallelism
//
// Run:
//   go run task1.go
//   GOMAXPROCS=1 go run task1.go
//
// cpuWork is CPU-bound: it spends all its time calculating (counting primes).
// It only gets faster if there are free CPU cores to actually run the extra
// goroutines in parallel.
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

func countPrimes(from, to int) int {
	count := 0
	for n := from; n < to; n++ {
		if isPrime(n) {
			count++
		}
	}
	return count
}

func cpuWork(workers int) (int, time.Duration) {
	const N = 6_000_000
	start := time.Now()
	results := make([]int, workers)
	var wg sync.WaitGroup
	chunk := N / workers
	for w := 0; w < workers; w++ {
		w := w
		wg.Add(1)
		go func() {
			defer wg.Done()
			from := w * chunk
			to := from + chunk
			if w == workers-1 {
				to = N // last chunk absorbs the remainder
			}
			results[w] = countPrimes(from, to)
		}()
	}
	wg.Wait()

	total := 0
	for _, r := range results {
		total += r
	}
	return total, time.Since(start)
}

// ioWork is I/O-bound: each "request" just waits, like a database call.
// Waiting goroutines don't use the CPU, so this scales almost for free.
func ioWork(workers int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func main() {
	fmt.Println("runtime.NumCPU():", runtime.NumCPU())
	fmt.Println("runtime.GOMAXPROCS(0):", runtime.GOMAXPROCS(0))
	fmt.Println()

	for _, w := range []int{1, 2, 4, runtime.NumCPU(), runtime.NumCPU() * 2} {
		primes, d := cpuWork(w)
		fmt.Printf("CPU work, %d goroutine(s): %d primes in %v\n", w, primes, d)
	}
	fmt.Println()

	for _, w := range []int{1, 10, 1000} {
		d := ioWork(w)
		fmt.Printf("I/O work, %4d goroutine(s): %v\n", w, d)
	}
}
