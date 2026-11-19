package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Task 3: Remove the allocations

// JoinIDsSlow builds "1,2,3" one string concatenation at a time.
func JoinIDsSlow(ids []int) string {
	s := ""
	for i, id := range ids {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprint(id)
	}
	return s
}

// JoinIDsFast builds the same result with a single growable buffer.
func JoinIDsFast(ids []int) string {
	var sb strings.Builder
	sb.Grow(len(ids) * 8)
	var tmp [20]byte
	for i, id := range ids {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.Write(strconv.AppendInt(tmp[:0], int64(id), 10))
	}
	return sb.String()
}

func main() {
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("JoinIDsSlow:", JoinIDsSlow(ids))
	fmt.Println("JoinIDsFast:", JoinIDsFast(ids))
	fmt.Println("run: go test -gcflags='-m -l' ./... ; go test -bench=. -benchmem")
}
