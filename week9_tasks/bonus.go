package main

import "fmt"

// Bonus question
//
// Run:
//   go run bonus.go

func main() {
	fmt.Println("=== closed channel ===")
	ch := make(chan int)
	close(ch)
	select {
	case v, ok := <-ch:
		fmt.Println("closed:", v, ok)
	default:
		fmt.Println("default")
	}
	fmt.Println(`Always prints "closed: 0 false". Receiving from a CLOSED channel never
blocks — it is immediately ready and returns the zero value with ok==false
every single time. Since that case is always ready, select always picks
it over default (default only runs when NO case is ready).`)

	fmt.Println("\n=== nil channel, with default ===")
	var nilCh chan int
	select {
	case v, ok := <-nilCh:
		fmt.Println("nil:", v, ok)
	default:
		fmt.Println("default")
	}
	fmt.Println(`Always prints "default". A nil channel blocks forever on both send and
receive — it is never ready. So the <-nilCh case can never fire, and
select falls through to default every time.`)

	fmt.Println("\n=== nil channel, default removed ===")
	fmt.Println(`select {
case v, ok := <-nilCh:
    fmt.Println(v, ok)
}
With no default and no other case, select has exactly one case and it can
never become ready — so this blocks forever. With only one goroutine
(main) and nothing else running, that's an immediate:
fatal error: all goroutines are asleep - deadlock!
(We don't run this branch here — it would just hang/crash the demo on
purpose, same as the deadlock examples in Part 1.)`)
}
