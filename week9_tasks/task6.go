package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Task 6: A worker pool with backpressure and cancellation
//
// Run:
//   go run task6.go

type Result struct {
	Word  string
	Upper string
}

// producer owns and closes jobs. It stops producing (and lets jobs close)
// as soon as done is closed, even mid-stream.
func producer(words []string, done <-chan struct{}) <-chan string {
	jobs := make(chan string)
	go func() {
		defer close(jobs)
		for _, w := range words {
			select {
			case jobs <- w:
			case <-done: // the consumer gave up: stop producing
				return
			}
		}
	}()
	return jobs
}

// startWorkers runs n workers reading jobs and sending Result; a closer
// goroutine closes results once all workers are done.
func startWorkers(n int, jobs <-chan string, done <-chan struct{}) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range jobs {
				res := Result{Word: w, Upper: strings.ToUpper(w)}
				select {
				case results <- res:
				case <-done:
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

// producerLeaky and startWorkersLeaky: the same pipeline with the <-done
// cases removed, to show what leaks without them.
func producerLeaky(words []string) <-chan string {
	jobs := make(chan string)
	go func() {
		defer close(jobs)
		for _, w := range words {
			jobs <- w // no done case: blocks forever if nobody ever reads this
		}
	}()
	return jobs
}

func startWorkersLeaky(n int, jobs <-chan string) <-chan Result {
	results := make(chan Result)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range jobs {
				results <- Result{Word: w, Upper: strings.ToUpper(w)} // no done case
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

var words = []string{"go", "channel", "select", "worker", "pool", "backpressure", "context", "scheduler", "goroutine", "pipeline"}

func demoRunToEnd() {
	fmt.Println("=== run to the end with range ===")
	done := make(chan struct{})
	results := startWorkers(3, producer(words, done), done)
	for r := range results {
		fmt.Printf("  %s -> %s\n", r.Word, r.Upper)
	}
	close(done) // nothing left depending on it, but tidy up like a real caller would
	time.Sleep(20 * time.Millisecond)
	fmt.Println("NumGoroutine after full run:", runtime.NumGoroutine())
	fmt.Println()
}

func demoEarlyStop() {
	fmt.Println("=== take only 3 results, then close(done) ===")
	done := make(chan struct{})
	results := startWorkers(3, producer(words, done), done)
	count := 0
	for r := range results {
		fmt.Printf("  %s -> %s\n", r.Word, r.Upper)
		count++
		if count == 3 {
			break
		}
	}
	close(done)
	time.Sleep(50 * time.Millisecond) // give producer/workers/closer time to notice done and exit
	fmt.Println("NumGoroutine after early stop + close(done):", runtime.NumGoroutine(), "(must be 1)")
	fmt.Println()
}

func demoEarlyStopLeaky() {
	fmt.Println("=== same, but with the <-done cases removed ===")
	before := runtime.NumGoroutine()
	jobs := producerLeaky(words)
	results := startWorkersLeaky(3, jobs)
	count := 0
	for r := range results {
		fmt.Printf("  %s -> %s\n", r.Word, r.Upper)
		count++
		if count == 3 {
			break
		}
	}
	// There is no done channel to close anymore — nothing can tell the
	// producer or the remaining workers to stop.
	time.Sleep(50 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Println("NumGoroutine before:", before, " after taking only 3 (no done to cancel with):", after)
	fmt.Println(`Without the <-done cases, the producer goroutine is stuck on "jobs <- w"
for whichever word came after the 3rd result (nobody is left reading
jobs/results fast enough, or at all, once main stops ranging), and the
workers are similarly stuck trying to send further results into a channel
main no longer drains. Both leak forever — this is exactly the Week 8
goroutine-leak pattern, just with more goroutines involved at once.`)
	fmt.Println()
}

// --- Submit: bounded wait, then give up ---

var ErrQueueFull = errors.New("queue is full")

func Submit(queue chan<- string, job string, timeout time.Duration) error {
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case queue <- job:
		return nil
	case <-t.C:
		return ErrQueueFull // an HTTP server would answer 503 "try again later"
	}
}

func demoSubmit() {
	fmt.Println("=== Submit(queue, job, timeout) ===")
	queue := make(chan string, 1)

	err1 := Submit(queue, "first", 50*time.Millisecond)
	fmt.Println(`Submit("first", 50ms) with an empty queue (cap 1) ->`, err1, "(success: fits in the free slot)")

	// queue is now full (cap 1, holds "first"); nobody is draining it, so
	// the next Submit has to wait out its timeout and give up.
	err2 := Submit(queue, "second", 50*time.Millisecond)
	fmt.Println(`Submit("second", 50ms) with a full queue       ->`, err2, "(gives up: ==ErrQueueFull is", errors.Is(err2, ErrQueueFull), ")")
	fmt.Println()
}

func main() {
	demoRunToEnd()
	demoEarlyStop()
	demoEarlyStopLeaky()
	demoSubmit()
}
