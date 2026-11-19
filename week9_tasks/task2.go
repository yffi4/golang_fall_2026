package main

import (
	"fmt"
	"time"
)

// Task 2: Buffered vs unbuffered channels
//
// Run:
//   go run task2.go

func run(name string, ch chan int) {
	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }
	done := make(chan struct{})
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
			fmt.Printf("%s: sent %d at %v\n", name, i, since())
		}
		close(done)
	}()
	time.Sleep(300 * time.Millisecond) // the receiver is busy at first
	for i := 0; i < 3; i++ {
		<-ch
		time.Sleep(100 * time.Millisecond) // and slow
	}
	<-done
}

// demoCapacities runs the experiment above for cap 0, 1, 2, 3, 10.
func demoCapacities() {
	fmt.Println("=== 1. run(...) with capacity 0, 1, 2, 3, 10 ===")
	for _, bufSize := range []int{0, 1, 2, 3, 10} {
		run(fmt.Sprintf("cap=%-2d", bufSize), make(chan int, bufSize))
		fmt.Println()
	}
	fmt.Println(`Explanation per row:
cap=0  (unbuffered): every send blocks until the receiver is ready. The
       receiver is busy for the first 300ms (initial Sleep) then takes one
       value every 100ms. So sends land at ~300ms, ~400ms, ~500ms.
cap=1: the first send fills the 1 free slot immediately (~0ms). The 2nd
       send then blocks until the receiver frees a slot, which happens at
       300ms (receiver's first Sleep). The 3rd send waits for the next
       slot at ~400ms.
cap=2: first two sends fill both slots at ~0ms. The 3rd send blocks until
       a slot frees at 300ms.
cap=3: all three sends fit in the buffer immediately, ~0ms each — the
       sender never has to wait for the receiver at all.
cap=10: same as cap=3 — buffer is bigger than needed, so all three sends
       land at ~0ms too. Making the buffer bigger than the number of
       values sent has no further effect.`)
	fmt.Println()
}

// demoLenCapAfterClose: len, cap, and reading buffered values after close.
func demoLenCapAfterClose() {
	fmt.Println("=== 2. len, cap, and reading after close ===")
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	fmt.Println("len:", len(ch), "cap:", cap(ch))
	close(ch)
	v, ok := <-ch
	fmt.Println(v, ok)
	v, ok = <-ch
	fmt.Println(v, ok)
	v, ok = <-ch
	fmt.Println(v, ok)
	fmt.Println()
}

// demoSmallSends: the "a", "b" question — cap=2 vs cap=1.
func demoSmallSends() {
	fmt.Println("=== 3. small sends: cap=2 then cap=1 ===")

	fmt.Println("-- cap=2 --")
	ch2 := make(chan string, 2)
	ch2 <- "a"
	ch2 <- "b"
	fmt.Println(<-ch2, <-ch2)
	fmt.Println(`Why it doesn't deadlock: both sends fit in the buffer (cap=2), so
neither one blocks — "a" and "b" just sit in the queue. Only then do we
read them back out, after both sends already finished.`)

	fmt.Println("-- cap=1 --")
	ch1 := make(chan string, 1)
	ch1 <- "a"
	fmt.Println("sent a (fits in the 1 slot)")
	done := make(chan struct{})
	go func() {
		ch1 <- "b" // this would block forever here: the slot is full and
		// nobody will ever read "a" in this goroutine's lifetime in main's
		// current order — so we demonstrate it with a timeout instead of
		// actually deadlocking the whole program.
		close(done)
	}()
	select {
	case <-done:
		fmt.Println("second send completed")
	case <-time.After(200 * time.Millisecond):
		fmt.Println(`second send ("b") -> still blocked after 200ms:
with cap=1, the buffer already holds "a" and is full. Sending "b" blocks
until someone receives "a" to free the slot. If this were written exactly
like the original snippet (both sends directly in main, back to back,
with no concurrent receiver), the program would deadlock:
fatal error: all goroutines are asleep - deadlock!`)
		<-ch1 // drain so this demo's own goroutine can finish and not leak
		<-done
	}
	fmt.Println()
}

func main() {
	demoCapacities()
	demoLenCapAfterClose()
	demoSmallSends()

	fmt.Println(`Participation Task (Week 8 fetchLeaky fix):
make(chan int, 1) was enough there because fetchLeaky's goroutine sends
EXACTLY ONE value (the Report) and then returns — one buffer slot is all
that single send ever needs to not block, no matter whether the caller is
still listening. If that goroutine sent TWO values on the same channel, a
buffer of 1 would only cover the first send: the second send would still
block until a receiver takes the first value out, so the goroutine could
still leak waiting on the second send if the caller already gave up.`)
}
