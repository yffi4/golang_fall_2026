package main

import (
	"fmt"
	"time"
)

// Task 1: Channel semantics
//
// Run:
//   go run task1.go

// --- 1. ping example ---

// demoPingOriginal: Sleep is in the SENDER. The receiver's <-ch does all the
// waiting (500ms), even though the receiver never called Sleep itself.
func demoPingOriginal() {
	fmt.Println("=== 1a. ping example: Sleep in sender (original) ===")
	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }
	ch := make(chan string) // unbuffered
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("sender:   about to send", since())
		ch <- "ping"
		fmt.Println("sender:   send finished", since())
	}()
	fmt.Println("receiver: waiting...  ", since())
	msg := <-ch
	fmt.Println("receiver: got", msg, "     ", since())
	fmt.Println()
}

// demoPingSleepInReceiver: Sleep is in the RECEIVER, before <-ch. Now the
// sender has nothing to wait for and reaches "ch <- ping" almost
// immediately — it is the SENDER that blocks on the channel operation this
// time, waiting for main to wake up and receive.
func demoPingSleepInReceiver() {
	fmt.Println("=== 1b. ping example: Sleep in receiver (modified) ===")
	start := time.Now()
	since := func() time.Duration { return time.Since(start).Round(100 * time.Millisecond) }
	ch := make(chan string) // unbuffered
	go func() {
		fmt.Println("sender:   about to send", since())
		ch <- "ping"
		fmt.Println("sender:   send finished", since())
	}()
	fmt.Println("receiver: waiting...  ", since())
	time.Sleep(500 * time.Millisecond)
	msg := <-ch
	fmt.Println("receiver: got", msg, "     ", since())
	fmt.Println(`
What changed: total wall-clock time is still ~500ms, but WHO is blocked on
the channel operation flips. In 1a, the sender sleeps and the receiver's
<-ch is what waits. In 1b, the sender has no Sleep, so it reaches "ch <-
ping" almost at 0s and blocks there until main finishes its own Sleep and
reaches <-ch at 500ms — now the SEND is what waits. Either way, the
rendezvous happens at ~500ms: an unbuffered channel always makes one side
wait for the other, no matter which side got there first.`)
	fmt.Println()
}

// --- 2. the nil / open / closed table: 12 tiny programs ---

// runWithTimeout runs a blocking channel op in its own goroutine and reports
// whether it completed within the timeout. We can't actually wait forever
// (that would hang this demo), so "did not finish in time" stands in for
// "blocks forever" — same idea as a goroutine leak in Week 8: the goroutine
// below is abandoned on purpose, to show the blocking behavior.
func runWithTimeout(label string, op func()) {
	done := make(chan struct{})
	go func() {
		op()
		close(done)
	}()
	select {
	case <-done:
		fmt.Println(label, "-> completed")
	case <-time.After(200 * time.Millisecond):
		fmt.Println(label, "-> still blocked after 200ms (would block forever)")
	}
}

// runRecovering runs op and reports a panic if one happens, without
// crashing the whole program.
func runRecovering(label string, op func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(label, "-> panic:", r)
		}
	}()
	op()
	fmt.Println(label, "-> no panic")
}

func demoNilOpenClosedTable() {
	fmt.Println("=== 2. nil / open / closed channel table (12 cases) ===")

	// --- nil channel ---
	var nilCh chan int // zero value: nil
	runWithTimeout("nil:    send             ", func() { nilCh <- 1 })
	runWithTimeout("nil:    receive          ", func() { <-nilCh })
	runRecovering("nil:    close             ", func() { close(nilCh) })
	fmt.Println("nil:    len/cap           ->", len(nilCh), cap(nilCh))

	// --- open channel (unbuffered: send/receive block until the other side
	// shows up; we provide the other side so these actually complete) ---
	openUnbuf := make(chan int)
	runWithTimeout("open:   send (unbuf)     ", func() {
		go func() { <-openUnbuf }() // give it a receiver so it can complete
		openUnbuf <- 1
	})
	openUnbuf2 := make(chan int)
	runWithTimeout("open:   receive (unbuf)  ", func() {
		go func() { openUnbuf2 <- 1 }() // give it a sender so it can complete
		<-openUnbuf2
	})
	openBuf := make(chan int, 2)
	runRecovering("open:   close             ", func() { close(openBuf) })
	openBuf2 := make(chan int, 3)
	openBuf2 <- 1
	openBuf2 <- 2
	fmt.Println("open:   len/cap           ->", len(openBuf2), cap(openBuf2))

	// --- closed channel ---
	closedCh := make(chan int)
	close(closedCh)
	runRecovering("closed: send              ", func() { closedCh <- 1 })
	closedBuf := make(chan int, 2)
	closedBuf <- 10
	close(closedBuf)
	v1, ok1 := <-closedBuf
	v2, ok2 := <-closedBuf
	fmt.Println("closed: receive           -> first:", v1, ok1, " second (empty):", v2, ok2)
	runRecovering("closed: close (again)     ", func() { close(closedBuf) })
	fmt.Println("closed: len/cap           ->", len(closedBuf), cap(closedBuf))

	fmt.Println()
}

// --- 3. say("goroutine") from Week 8, rewritten with a channel instead of
// sync.WaitGroup ---

func sayWithChannel(msg string) {
	done := make(chan struct{})
	go func() {
		defer close(done) // signal "I'm finished" by closing, not by sending a value
		for i := 0; i < 3; i++ {
			fmt.Println(msg)
		}
	}()
	<-done // blocks until the goroutine above closes done
}

func demoSayWithChannel() {
	fmt.Println("=== 3. say(\"goroutine\") using a channel instead of WaitGroup ===")
	sayWithChannel("goroutine")
	fmt.Println()
}

func main() {
	demoPingOriginal()
	demoPingSleepInReceiver()
	demoNilOpenClosedTable()
	demoSayWithChannel()
}
