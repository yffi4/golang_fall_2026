package main

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Task 5: select
//
// Run:
//   go run task5.go
//   go run -race task5.go

// --- 1. random choice, and making select prefer a ---

func demoRandomChoice() {
	fmt.Println("=== 1. random choice: both a and b have 1000 values ready ===")
	for run := 1; run <= 3; run++ {
		a := make(chan int, 1000)
		b := make(chan int, 1000)
		for i := 0; i < 1000; i++ {
			a <- i
			b <- i
		}
		count := map[string]int{}
		for i := 0; i < 1000; i++ {
			select {
			case v := <-a:
				_ = v
				count["a"]++
			case v := <-b:
				_ = v
				count["b"]++
			}
		}
		fmt.Printf("run %d: a=%d b=%d\n", run, count["a"], count["b"])
	}
	fmt.Println("The split is close to 50/50 each run, but never exactly — select picks uniformly at random among ready cases, it does not guarantee an even count.")

	fmt.Println("\n=== prefer a: two selects, first one with a default ===")
	a := make(chan int, 1000)
	b := make(chan int, 1000)
	for i := 0; i < 1000; i++ {
		a <- i
	}
	for i := 0; i < 200; i++ {
		b <- i
	}
	count := map[string]int{}
	for len(a) > 0 || len(b) > 0 {
		select {
		case v := <-a: // checked first, non-blocking: if a has something, always take it
			_ = v
			count["a"]++
			continue
		default:
		}
		select {
		case v, ok := <-a:
			_ = v
			if ok {
				count["a"]++
			}
		case v, ok := <-b:
			_ = v
			if ok {
				count["b"]++
			}
		}
	}
	fmt.Println("a:", count["a"], "b:", count["b"], "(a dominates: it's drained by the first, priority select before the random one ever runs)")
	fmt.Println()
}

// --- 2. fetchWithTimeout ---

func slowFetch(d time.Duration) int {
	time.Sleep(d)
	return 42
}

var errTimeout = errors.New("timeout")

// fetchWithTimeout waits up to timeout for slowFetch. The channel is
// buffered with capacity 1 so that if the timeout wins, the goroutine can
// still deliver its result into the buffer and exit on its own — instead
// of blocking forever on an unbuffered send that nobody will ever receive
// (Week 8 leak).
func fetchWithTimeout(timeout time.Duration) (int, error) {
	ch := make(chan int, 1)
	go func() {
		ch <- slowFetch(200 * time.Millisecond)
	}()
	select {
	case v := <-ch:
		return v, nil
	case <-time.After(timeout):
		return 0, errTimeout
	}
}

func demoFetchWithTimeout() {
	fmt.Println("=== 2. fetchWithTimeout ===")
	before := runtime.NumGoroutine()
	v, err := fetchWithTimeout(100 * time.Millisecond) // slower than the fetch itself
	fmt.Println("timeout=100ms, fetch=200ms ->", v, err)
	v2, err2 := fetchWithTimeout(500 * time.Millisecond) // slower than the timeout
	fmt.Println("timeout=500ms, fetch=200ms ->", v2, err2)
	time.Sleep(300 * time.Millisecond) // let the first fetch's goroutine actually finish and exit
	after := runtime.NumGoroutine()
	fmt.Println("NumGoroutine before:", before, "after (goroutine from the timed-out fetch had time to exit):", after)
	fmt.Println()
}

// --- 3. merge ---

// merge reads from a and b until BOTH are closed. Setting a (or b) to nil
// once it's closed switches that case off (a nil channel is never ready),
// so the loop doesn't keep hammering a channel that has nothing left to
// give.
func merge(a, b <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

// mergeBuggy is the same, but WITHOUT the "a = nil" (and "b = nil") line —
// once a is closed, its case keeps firing (closed channels are always
// ready) and the loop busy-spins receiving useless (0, false) values from
// it while waiting for b.
func mergeBuggy(a, b <-chan int) (<-chan int, *int64) {
	out := make(chan int)
	var uselessFromA int64
	aOpen, bOpen := true, true
	go func() {
		defer close(out)
		for aOpen || bOpen {
			select {
			case v, ok := <-a:
				if !ok {
					atomic.AddInt64(&uselessFromA, 1)
					continue // no "a = nil": this case stays switched on forever
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					bOpen = false
					continue
				}
				out <- v
			}
		}
	}()
	return out, &uselessFromA
}

func feed(vals []int, delay time.Duration) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range vals {
			time.Sleep(delay)
			ch <- v
		}
	}()
	return ch
}

func demoMerge() {
	fmt.Println("=== 3. merge(a, b) ===")
	a := feed([]int{1, 2}, 5*time.Millisecond)           // closes quickly
	b := feed([]int{100, 101, 102}, 40*time.Millisecond) // closes much later
	var got []int
	for v := range merge(a, b) {
		got = append(got, v)
	}
	fmt.Println("merged (with a = nil fix):", got)

	fmt.Println("\n-- without the a = nil line --")
	a2 := feed([]int{1, 2}, 5*time.Millisecond)
	b2 := feed([]int{100, 101, 102}, 40*time.Millisecond)
	out2, uselessFromA := mergeBuggy(a2, b2)
	// We know exactly 5 real values are coming (2 from a, 3 from b). We
	// read exactly that many and stop — mergeBuggy's goroutine itself
	// never closes out and never exits (see explanation below), so a
	// plain "range out2" would hang forever.
	var got2 []int
	for i := 0; i < 5; i++ {
		got2 = append(got2, <-out2)
	}
	fmt.Println("merged (buggy, no nil fix):", got2)
	fmt.Println("useless zero values received from already-closed a while we were reading:", atomic.LoadInt64(uselessFromA))
	fmt.Println(`
Why: per the Part 1 table, receiving from a CLOSED channel never blocks —
it immediately returns the zero value and ok=false. Without "a = nil",
a's case in the select stays switched on forever, and since it's always
immediately ready, select keeps firing it in a tight busy-loop (burning
CPU) every time b isn't ready yet, instead of actually waiting for b.
Worse than just wasted values: since "a" is never reassigned, the loop
guard "a != nil || b != nil" is permanently true (a, the channel
variable, is still non-nil — only closed) even after b also closes, so
this buggy goroutine never calls close(out) and never exits at all — a
genuine infinite busy-loop, not just extra zero values. That's why we
stopped reading after exactly 5 values instead of using range: the real
consequence of dropping that one line isn't "some wasted CPU", it's "this
merge never finishes". With "a = nil" restored, a's case is switched off
for good (a nil channel blocks forever = never ready) the moment a
closes, so the loop guard correctly goes false once both inputs are
exhausted, matching the PDF's measured 3039-4397 useless values before
the fix in a comparable run.`)
	fmt.Println()
}

// --- Participation Task: ball game ---

func demoBallGame() {
	fmt.Println("=== Participation Task: ball game ===")
	before := runtime.NumGoroutine()

	ball := make(chan int)
	done := make(chan struct{})
	var hits atomic.Int64

	player := func() {
		for {
			select {
			case v, ok := <-ball:
				if !ok {
					return
				}
				hits.Add(1)
				select {
				case ball <- v + 1:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); player() }()
	go func() { defer wg.Done(); player() }()

	go func() { ball <- 0 }() // kick off, in a goroutine so main doesn't block if both players are briefly busy

	time.Sleep(100 * time.Millisecond)
	close(done)
	wg.Wait() // both players guaranteed to exit: every blocking op in player() also selects on <-done

	time.Sleep(10 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Println("hits in 100ms:", hits.Load())
	fmt.Println("NumGoroutine before:", before, "after (players stopped via done):", after)
	fmt.Println()
}

func main() {
	demoRandomChoice()
	demoFetchWithTimeout()
	demoMerge()
	demoBallGame()
}
