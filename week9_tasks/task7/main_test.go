package main

import (
	"fmt"
	"testing"
)

func TestPassNSum(t *testing.T) {
	want := N * (N - 1) / 2
	for _, size := range []int{0, 1, 100, 10000} {
		if got := passN(size); got != want {
			t.Errorf("passN(%d) = %d, want %d", size, got, want)
		}
	}
}

func TestPassBatchesSum(t *testing.T) {
	want := N * (N - 1) / 2
	if got := passBatches(); got != want {
		t.Errorf("passBatches() = %d, want %d", got, want)
	}
}

var sink int

func BenchmarkPassN(b *testing.B) {
	for _, size := range []int{0, 1, 100, 10000} {
		b.Run(fmt.Sprintf("buf=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sink = passN(size)
			}
		})
	}
}

func BenchmarkPassBatches(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = passBatches()
	}
}
