package main

import "testing"

func makeInts(n int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	return nums
}

func BenchmarkSumInts(b *testing.B) {
	nums := makeInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumInts(nums)
	}
}

func BenchmarkSumGeneric(b *testing.B) {
	nums := makeInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sum(nums)
	}
}

func BenchmarkSumAny(b *testing.B) {
	nums := makeInts(1000)
	anyNums := ToAny(nums) // convert once, before the timer starts
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumAny(anyNums)
	}
}

func BenchmarkSumAnyWithBoxing(b *testing.B) {
	nums := makeInts(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SumAny(ToAny(nums))
	}
}

// Explanations (measured with `go test ./task6/... -bench=. -benchmem`
// on Apple M5 Pro, 1000 ints per run):
//
// - BenchmarkSumInts          227.0 ns/op      0 B/op      0 allocs/op
// - BenchmarkSumGeneric       227.4 ns/op      0 B/op      0 allocs/op
// - BenchmarkSumAny           227.5 ns/op      0 B/op      0 allocs/op
// - BenchmarkSumAnyWithBoxing 4779  ns/op  22336 B/op    745 allocs/op
//
// Why is the generic version as fast as the concrete one?
//   The compiler generates a fully specialized version of Sum for T=int
//   at compile time (monomorphization), so the generated machine code is
//   essentially identical to SumInts - there is no interface, no
//   dynamic dispatch and no boxing involved.
//
// Why isn't SumAny slower here even though it does 0 allocations?
//   In principle every loop iteration still pays for a type assertion
//   (n.(int)): a runtime check that the interface's dynamic type is
//   exactly int, instead of a plain inlined add. On this benchmark that
//   extra check is cheap and predictable enough (same type every single
//   time, branch predictor loves it) that it gets hidden in the noise
//   next to the add itself - the real cost of `any` here isn't the
//   assertion, it's the boxing, which is exactly what
//   BenchmarkSumAnyWithBoxing exposes once it happens on every call
//   instead of once up front.
//
// Where do the allocs/op in BenchmarkSumAnyWithBoxing come from, and why
// isn't it exactly 1001 (1000 ints + 1 backing array)?
//   Go caches the boxed representation of small integers (0-255, see the
//   runtime's `staticuint64s` table), so converting those small values to
//   `any` reuses a pre-existing allocation instead of creating a new one
//   on the heap. Since our slice is 0..999, the values 0-255 are "free"
//   (no allocation), so we get one allocation for the []any backing
//   array plus one allocation per int outside the cached 0-255 range
//   (i.e. roughly 1 + (1000-256) = 745, not 1001).
