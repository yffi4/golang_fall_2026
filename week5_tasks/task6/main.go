package main

// Task 6: Benchmark interfaces vs generics
//
// Run with:
//   go test ./task6/... -bench=. -benchmem

type Number interface {
	~int | ~int64 | ~float64
}

// SumInts is the concrete, hand-written version.
func SumInts(nums []int) int {
	s := 0
	for _, n := range nums {
		s += n
	}
	return s
}

// Sum is the generic version - the compiler generates a specialized
// version of this function for each instantiation (e.g. for []int), so
// at runtime it does exactly the same work as SumInts.
func Sum[T Number](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}

// SumAny takes a slice of `any` and must type-assert every element back
// to int before it can add it.
func SumAny(nums []any) int {
	total := 0
	for _, n := range nums {
		total += n.(int)
	}
	return total
}

// ToAny copies every int into a []any. Each int above a small cached
// range gets boxed into its own heap allocation to be stored as `any`.
func ToAny(nums []int) []any {
	out := make([]any, len(nums))
	for i, n := range nums {
		out[i] = n
	}
	return out
}

func main() {}
