package main

import (
	"fmt"
	"testing"
)

func TestCountParMatchesCountSeq(t *testing.T) {
	want := CountSeq(N)
	for _, parts := range []int{1, 2, 3, 8, 1000, 200_000} {
		got := CountPar(N, parts)
		if got != want {
			t.Errorf("CountPar(N, %d) = %d, want %d", parts, got, want)
		}
	}
}

var sink int

func BenchmarkCountSeq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = CountSeq(N)
	}
}

func BenchmarkCountPar(b *testing.B) {
	for _, parts := range []int{2, 8, 1000, 200_000} {
		b.Run(fmt.Sprintf("parts=%d", parts), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sink = CountPar(N, parts)
			}
		})
	}
}
