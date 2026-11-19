package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Task 3: How much does a goroutine cost?
//
// Run:
//   go run .
//   go run . 1000000   (try a bigger N if your laptop has enough memory)

func startAndStop(n int) (stackPerGoroutineKB float64, elapsed time.Duration) {
	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)

	stop := make(chan struct{})
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-stop
		}()
	}
	elapsed = time.Since(start)

	fmt.Println("goroutines alive:", runtime.NumGoroutine())
	runtime.ReadMemStats(&after)
	fmt.Printf("stack memory: %.1f MB\n", float64(after.StackSys)/1024/1024)

	stackPerGoroutineKB = float64(after.StackSys-before.StackSys) / float64(n) / 1024

	close(stop)
	wg.Wait()
	fmt.Println("after close:", runtime.NumGoroutine(), "goroutine(s)")
	return
}

func main() {
	fmt.Println("=== 100 000 goroutines ===")
	kb, d := startAndStop(100_000)
	fmt.Printf("%.2f KB per goroutine, started in %v\n\n", kb, d)

	fmt.Println("=== 1 000 000 goroutines ===")
	kb, d = startAndStop(1_000_000)
	fmt.Printf("%.2f KB per goroutine, started in %v\n", kb, d)
}
