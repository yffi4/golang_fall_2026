package main

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"
)

// Task 4: Explore the scheduler
//
// Run:
//   go run task4.go
//   GODEBUG=schedtrace=500 go run task4.go
//   GODEBUG=schedtrace=500 GOMAXPROCS=1 go run task4.go
//   GODEBUG=asyncpreemptoff=1 go run task4.go   (the busy loop will hang: Ctrl+C)

// busy burns CPU for d without making any function call that the old,
// cooperative scheduler could use as a preemption point.
func busy(d time.Duration) {
	end := time.Now().Add(d)
	for time.Now().Before(end) {
	}
}

// demoBusyGoroutines starts 8 goroutines that are each busy for 2 seconds —
// run this binary with GODEBUG=schedtrace=500 to watch the scheduler react.
func demoBusyGoroutines() {
	fmt.Println("--- 4.1: 8 busy goroutines (watch with GODEBUG=schedtrace=500) ---")
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			busy(500 * time.Millisecond)
		}()
	}
	wg.Wait()
	fmt.Println("done")
}

// demoPreemption shows that since Go 1.14 a goroutine that never makes a
// function call can still be preempted: main keeps running even though one
// CPU is "stuck" in an infinite loop next to it.
func demoPreemption() {
	fmt.Println("\n--- 4.2: preemption (try GODEBUG=asyncpreemptoff=1 to see it hang) ---")
	runtime.GOMAXPROCS(1) // only one goroutine can run at a time
	go func() {
		for { // a busy loop with no function calls
		}
	}()
	time.Sleep(10 * time.Millisecond)
	fmt.Println("main is still alive")
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// demoThreadCreate compares how many OS threads get created for 100
// goroutines that wait with time.Sleep (parked by the runtime, no thread
// needed) versus 100 goroutines that block in a real blocking syscall (each
// one locks its own OS thread in the kernel).
func demoThreadCreate() {
	fmt.Println("\n--- 4.3: threads created, time.Sleep vs blocking syscall ---")

	threads := func() int { return pprof.Lookup("threadcreate").Count() }
	run := func(label string, n int, f func()) {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				f()
			}()
		}
		wg.Wait()
		fmt.Printf("%-24s threads created so far: %d\n", label, threads())
	}

	fmt.Printf("%-24s threads created so far: %d\n", "start", threads())
	run("100 x time.Sleep", 100, func() {
		time.Sleep(200 * time.Millisecond)
	})
	run("100 x blocking syscall", 100, func() {
		// syscall.Nanosleep only exists on Linux; syscall.Select with no fds
		// is the portable equivalent — a real blocking syscall on any OS.
		tv := syscall.NsecToTimeval(int64(200 * time.Millisecond))
		syscall.Select(0, nil, nil, nil, &tv) // the whole OS thread is blocked
	})
}

func main() {
	demoBusyGoroutines()
	demoPreemption()
	demoThreadCreate()

	fmt.Println(`
--- 4.4: why is "worker 5" printed first? (see task2.go's WaitGroup demo) ---
Each P keeps a special one-slot fast lane called runnext for the goroutine
it most recently made runnable. When the loop does go func(){...}() five
times in a row on the SAME P, goroutine 5 is placed into that P's runnext
slot, bumping out whatever was there. When the scheduler needs a new G to
run, it checks runnext FIRST — so goroutine 5 tends to run before 1..4,
which are waiting in the regular local run queue behind it. This is an
implementation detail, not a language guarantee: with GOMAXPROCS=1 it is
very consistent (one P, no other core to interfere), but it is still not
part of the Go specification, so real code must never depend on it.`)
}
