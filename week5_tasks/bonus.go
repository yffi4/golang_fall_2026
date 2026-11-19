package main

import "fmt"

// Bonus question
//
// Why does this compile at all?
//   `any` (interface{}) is always a comparable *type* as far as the
//   compiler is concerned, so every `== ` below type-checks fine,
//   regardless of what ends up stored inside at runtime.

func bonusDemo() {
	var a, b any = 1, 1
	// a == b -> true: both hold dynamic type int and the same value 1.
	fmt.Println("a == b:", a == b)

	var c, d any = 1, int64(1)
	// c == d -> false: comparing two interface values first compares
	// their dynamic *types*. int and int64 are different types, so the
	// values are unequal even though they "look like" the same number.
	fmt.Println("c == d:", c == d)

	var e, f any = []int{1}, []int{1}
	// e == f -> this is the line that does something "much worse than
	// print false": both hold dynamic type []int, which is NOT
	// comparable (slices can only be compared to nil). Go only catches
	// this at runtime, so instead of a boolean we get a panic:
	//   panic: runtime error: comparing uncomparable type []int
	// recover() here just lets the program keep running so we can
	// print what happened instead of crashing the whole demo.
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("e == f panicked:", r)
			}
		}()
		fmt.Println("e == f:", e == f)
	}()
}

func main() {
	fmt.Println("=== Bonus question ===")
	bonusDemo()
}
