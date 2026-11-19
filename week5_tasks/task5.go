package main

import (
	"cmp"
	"fmt"
)

// Task 5: Generic toolbox

// Number uses `~` so that named types with an underlying int/int64/float64
// (like Score below) satisfy the constraint too. Without the `~`, Number
// would be the *exact* type set {int, int64, float64}, and Score (whose
// underlying type is int but whose named type is different) would fail
// to compile with:
//
//	Score does not satisfy Number (possibly missing ~ for int in Number)
//
// That's point 5 of the task: removing the `~` breaks Sum/Average for
// Score, so it stays in.
type Number interface {
	~int | ~int64 | ~float64
}

type Score int

func Sum[T Number](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}

func Average[T Number](nums []T) (float64, bool) {
	if len(nums) == 0 {
		return 0, false
	}
	return float64(Sum(nums)) / float64(len(nums)), true
}

func MaxOf[T cmp.Ordered](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	max := items[0]
	for _, v := range items[1:] {
		if v > max {
			max = v
		}
	}
	return max, true
}

func IndexOf[T comparable](items []T, target T) int {
	for i, v := range items {
		if v == target {
			return i
		}
	}
	return -1
}

func Filter[T any](items []T, keep func(T) bool) []T {
	result := make([]T, 0, len(items))
	for _, v := range items {
		if keep(v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Task 5: Generic toolbox ===")

	ints := []int{10, 20, 30}
	fmt.Println("Sum(ints):", Sum(ints))
	avg, ok := Average(ints)
	fmt.Println("Average(ints):", avg, ok)

	scores := []Score{60, 75, 90}
	fmt.Println("Sum(scores):", Sum(scores))
	avgScore, ok := Average(scores)
	fmt.Println("Average(scores):", avgScore, ok)

	avgEmpty, ok := Average([]int{})
	fmt.Println("Average(empty):", avgEmpty, ok)

	maxInt, ok := MaxOf([]int{4, 9, 2})
	fmt.Println("MaxOf(ints):", maxInt, ok)

	maxStr, ok := MaxOf([]string{"banana", "apple", "cherry"})
	fmt.Println("MaxOf(strings):", maxStr, ok)

	maxEmpty, ok := MaxOf([]int{})
	fmt.Println("MaxOf(empty):", maxEmpty, ok)

	fmt.Println("IndexOf(found):", IndexOf([]string{"go", "rust", "python"}, "rust"))
	fmt.Println("IndexOf(missing):", IndexOf([]string{"go", "rust"}, "java"))

	passing := Filter(scores, func(s Score) bool { return s >= 75 })
	fmt.Println("Filter(scores >= 75):", passing)

	// Point 5: uncommenting the line below (after removing `~` from
	// Number) fails to compile with:
	//   Score does not satisfy Number (possibly missing ~ for int in Number)
	// because without `~`, Number's type set is exactly {int, int64, float64}
	// and Score is a distinct named type, even though its underlying type
	// is int.
}
