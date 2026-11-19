package main

import (
	"fmt"
)

// Task 4: Tour of Go -> Exercise: Errors

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z, nil
}

// BuggyError is the version of Error() that would cause infinite
// recursion. It's a separate type (ErrNegativeSqrt already has its own
// correct Error() above) defined only for discussion - run `go vet` on
// this file to see it caught at compile time, before it ever runs:
//
//	task4.go:43:51: fmt.Sprintf format %v with arg e causes recursive (BuggyError).Error method call
//
// Why: fmt.Sprintf("%v", e) formats e, and because e implements the
// error interface, fmt recognizes that and calls e.Error() to produce
// the string - which calls fmt.Sprintf("%v", e) again, forever (a stack
// overflow in practice, not a clean panic). The fix is the same one used
// in ErrNegativeSqrt.Error() above: convert to the underlying type
// first, e.g. %v with float64(e), so fmt formats a plain float64
// instead of calling Error() again.
type BuggyError float64

func (e BuggyError) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", e) // %v with e, not float64(e)
}

func main() {
	v, err := Sqrt(2)
	fmt.Println(v, err)

	v, err = Sqrt(-2)
	fmt.Println(v, err)
}
