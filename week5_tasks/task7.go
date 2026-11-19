package main

import (
	"fmt"
	"strings"
)

// Task 7: Tour of Go -> Generic types
//
// Starting point from the tour:
//
//	type List[T any] struct {
//		next *List[T]
//		val  T
//	}
//
// List[T] represents a node of a singly-linked list; a `nil *List[T]`
// is the empty list. Push can't reassign the caller's pointer variable
// (a Go method only ever receives a copy of the pointer), so it returns
// the new head instead - the caller writes `l = l.Push(v)`. That's
// exactly what makes Push work starting from an empty (nil) list: the
// very first call has a nil receiver, and `next: l` simply stores that
// nil, producing a valid one-node list.

type List[T any] struct {
	next *List[T]
	val  T
}

func (l *List[T]) Push(v T) *List[T] {
	return &List[T]{next: l, val: v}
}

func (l *List[T]) Len() int {
	n := 0
	for cur := l; cur != nil; cur = cur.next {
		n++
	}
	return n
}

func (l *List[T]) String() string {
	if l == nil {
		return "[]"
	}
	parts := make([]string, 0, l.Len())
	for cur := l; cur != nil; cur = cur.next {
		parts = append(parts, fmt.Sprint(cur.val))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func main() {
	fmt.Println("=== Task 7: generic List[T] ===")

	var l *List[int] // empty (nil) list
	fmt.Println("empty list:", l, "len:", l.Len())

	l = l.Push(1)
	l = l.Push(2)
	l = l.Push(3)
	fmt.Println("after three Push calls:", l, "len:", l.Len())

	var names *List[string]
	names = names.Push("go")
	names = names.Push("rust")
	fmt.Println("string list:", names, "len:", names.Len())
}
