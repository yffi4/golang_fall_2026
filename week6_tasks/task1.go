package main

import "fmt"

// Task 1: Closures

// makeAccumulator returns a function that adds its argument to a running
// total (kept in a variable captured by the closure) and returns the total.
func makeAccumulator() func(int) int {
	total := 0
	return func(n int) int {
		total += n
		return total
	}
}

// memoize remembers results of f in a map, so f runs once per distinct input.
func memoize(f func(int) int) func(int) int {
	cache := make(map[int]int)
	return func(n int) int {
		if v, ok := cache[n]; ok {
			return v
		}
		v := f(n)
		cache[n] = v
		return v
	}
}

// Participation Task: counter() from the lecture.
//
// `count` is captured by the returned closure, so the compiler moves it to
// the heap (escape analysis) and it lives as long as the closure does.
// Every call of counter() creates a NEW count variable, so c1 and c2 don't
// share it.
func counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	acc := makeAccumulator()
	fmt.Println(acc(10), acc(5), acc(-3)) // 10 15 12

	calls := 0
	slow := func(n int) int {
		calls++
		return n * n
	}
	fast := memoize(slow)
	for _, n := range []int{3, 3, 4, 3, 4, 5} {
		fmt.Printf("fast(%d) = %d\n", n, fast(n))
	}
	fmt.Println("slow was called", calls, "times for 3 distinct inputs") // 3

	c1 := counter()
	fmt.Println(c1(), c1(), c1()) // 1 2 3
	c2 := counter()
	fmt.Println(c2(), c1()) // 1 4

	x := 10
	show := func() { fmt.Println("x =", x) }
	x = 20
	show() // x = 20: closures capture variables, not values

	triple := func(k int) func(int) int {
		return func(n int) int { return n * k }
	}(3)
	fmt.Println(triple(5)) // 15

	// Part 3 of the task: see task1_loopvar/go121 and task1_loopvar/go122.
}
