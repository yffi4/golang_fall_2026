package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/pprof"
	"time"
)

// Task 5: Find and fix the leaks
//
// Run:
//   go run .

func slowQuery() int {
	time.Sleep(50 * time.Millisecond) // pretend to ask a database
	return 42
}

// fetchLeaky leaks a goroutine on every timeout: the unbuffered channel send
// blocks forever once nobody is left to receive from it.
func fetchLeaky(timeout time.Duration) (int, error) {
	ch := make(chan int) // unbuffered
	go func() {
		ch <- slowQuery() // blocks forever if fetchLeaky already returned
	}()
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

// fetchFixed: give the channel a buffer of one, so the send always
// succeeds immediately and the goroutine can finish and exit on its own,
// whether or not anyone ever reads the value.
func fetchFixed(timeout time.Duration) (int, error) {
	ch := make(chan int, 1)
	go func() {
		ch <- slowQuery()
	}()
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(timeout):
		return 0, errors.New("timeout")
	}
}

// startTickerLeaky starts a worker with no way to ever stop it.
func startTickerLeaky() {
	go func() {
		for {
			time.Sleep(10 * time.Millisecond) // "check something"
		}
	}()
}

// startTicker: give the caller a way to say "stop" via context.
func startTicker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
		}
	}()
}

func goroutineProfile(debug int) string {
	var buf bytes.Buffer
	_ = pprof.Lookup("goroutine").WriteTo(&buf, debug)
	return buf.String()
}

func main() {
	fmt.Println("start:", runtime.NumGoroutine())

	fmt.Println("\n--- leaking versions ---")
	for i := 0; i < 100; i++ {
		_, _ = fetchLeaky(10 * time.Millisecond) // slowQuery takes 50ms > 10ms timeout: always leaks
	}
	fmt.Println("after 100 fetchLeaky:", runtime.NumGoroutine())

	for i := 0; i < 10; i++ {
		startTickerLeaky()
	}
	fmt.Println("after 10 leaky tickers:", runtime.NumGoroutine())

	fmt.Println("\ngoroutine profile (debug=1), trimmed:")
	fmt.Println(goroutineProfile(1))

	fmt.Println("\n--- fixed versions ---")
	for i := 0; i < 100; i++ {
		_, _ = fetchFixed(10 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond) // give the 50ms queries time to finish and the goroutines to exit
	fmt.Println("after 100 fetchFixed:", runtime.NumGoroutine())

	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < 10; i++ {
		startTicker(ctx)
	}
	fmt.Println("after 10 stoppable tickers:", runtime.NumGoroutine())
	cancel()
	time.Sleep(50 * time.Millisecond) // give them time to notice ctx.Done()
	fmt.Println("after cancel():", runtime.NumGoroutine())

	fmt.Println("\nnote: the 10 leaky tickers from above are still alive and always will be —")
	fmt.Println("only process exit stops them. Count:", runtime.NumGoroutine())
}
