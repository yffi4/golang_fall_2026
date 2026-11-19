package main

import "fmt"

// Task 2: Predict, then check
//
// Predictions (written BEFORE running `go build -gcflags='-m -l' task2.go`):
//
//   f1  (x := 42; return x)          -> stays on stack: a plain value copy
//                                       is returned, nothing is shared up.
//   f2  (x := 42; return &x)         -> escapes: &x is returned, so x must
//                                       outlive f2's frame.
//   f3  (p *Point; return p.X+p.Y)   -> stays on stack: p is only read, the
//                                       pointer itself is never stored or
//                                       returned ("p does not escape").
//   f4  (p := Point{1,2}; Println(p))-> escapes: fmt.Println takes `any`,
//                                       so p gets boxed into an interface.
//   f5  (s := make([]int,10); ret s) -> escapes: the backing array is
//                                       returned to the caller.
//   f6  (s := make([]int,10); len(s))-> stays on stack: only len(s) (a
//                                       plain int) leaves the function, the
//                                       slice itself never does.
//   f7  (s := make([]int,n); len(s))-> escapes: size `n` is not a compile-
//                                       time constant, so the compiler
//                                       can't size a stack slot for it.
//   f8  (closure capturing x)        -> escapes: the returned closure
//                                       outlives f8, so x (and the closure
//                                       itself) must live on the heap.
//   f9  (global = p)                 -> leaking param: p itself doesn't
//                                       escape *inside* f9, but p is stored
//                                       into a global, so whatever the
//                                       caller passes in must already be on
//                                       the heap (or become so).
//   f10 (var big [200_000]byte; ...) -> escapes: 200,000 bytes is far past
//                                       the ~128 KB local-variable limit,
//                                       too big for a stack frame.
//   main (p := Point{3,4}; f9(&p))   -> because f9 leaks its parameter,
//                                       &p escapes too: p itself must move
//                                       to the heap even though it's a
//                                       local variable in main.
//
// Actual compiler report (measured with
// `go build -gcflags='-m -l' -o /dev/null task2.go`, go1.27.1), only the
// lines for f1-f10 and main:
//
//   task2.go:81:25: moved to heap: x                   (f2)
//   task2.go:82:9:  p does not escape                  (f3)
//   task2.go:83:55: p escapes to heap                   (f4, boxed by Println)
//   task2.go:84:34: make([]int, 10) escapes to heap     (f5)
//   task2.go:85:34: make([]int, 10) does not escape     (f6)
//   task2.go:86:34: make([]int, n) does not escape      (f7)  <- see below
//   task2.go:87:25: moved to heap: x                     (f8, captured var)
//   task2.go:87:40: func literal escapes to heap         (f8, the closure)
//   task2.go:88:9:  leaking param: p                     (f9)
//   task2.go:89:29: moved to heap: big                   (f10)
//   task2.go:102:2: moved to heap: p                     (main's p, because of f9(&p))
//   (f1 gets no line at all - nothing to report, exactly like newUserValue
//   in the lecture: it returns a plain copy, x stays on the stack.)
//
// Mismatches between my predictions and the report:
//   f7 was predicted to escape ("size n is not a compile-time constant"),
//   but go1.27.1 reports "make([]int, n) does not escape". This is exactly
//   the compiler-version warning on page 1 of the handout: Go 1.25/1.26
//   learned to keep more non-constant-size slices on the stack when the
//   compiler can still prove an upper bound/usage pattern that fits. The
//   lecture's own dynamicSize example (`make([]byte, n)`) was measured on
//   Go 1.24 and DID escape there - same source pattern, different answer
//   on a newer compiler. That's the whole point of "never rely on a value
//   being on the stack - measure instead" (Warning box, Part 2).
//   Every other prediction matched the report once the exact wording is
//   accounted for: f2/f8/f10 show up as "moved to heap: <var>" (the
//   variable itself moves) rather than "<expr> escapes to heap" (a value
//   escapes) - same outcome, different message shape depending on whether
//   the compiler is talking about a named local or an anonymous value.

type Point struct{ X, Y int }

var global *Point

func f1() int         { x := 42; return x }
func f2() *int        { x := 42; return &x }
func f3(p *Point) int { return p.X + p.Y }
func f4()             { p := Point{1, 2}; fmt.Println(p) }
func f5() []int       { s := make([]int, 10); return s }
func f6() int         { s := make([]int, 10); return len(s) }
func f7(n int) int    { s := make([]int, n); return len(s) }
func f8() func() int  { x := 0; return func() int { x++; return x } }
func f9(p *Point)     { global = p }
func f10() byte       { var big [200_000]byte; return big[0] }

func main() {
	fmt.Println(f1())
	fmt.Println(*f2())
	fmt.Println(f3(&Point{1, 2}))
	f4()
	fmt.Println(f5())
	fmt.Println(f6())
	fmt.Println(f7(10))
	fmt.Println(f8()())
	fmt.Println(f10())

	p := Point{3, 4}
	f9(&p) // p escapes: f9 leaks its parameter into the global `global`
	fmt.Println(*global)
}
