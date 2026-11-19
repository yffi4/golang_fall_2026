package main

import (
	"errors"
	"fmt"
)

// Task 3: Find the bug

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string { return "invalid field: " + e.Field }

// buggyValidate reproduces the classic "non-nil nil error" bug: it
// declares a typed *ValidationError pointer and returns it through the
// error interface.
func buggyValidate(age int) error {
	var err *ValidationError
	if age < 0 {
		err = &ValidationError{Field: "age"}
	}
	return err
}

// fixedValidate follows the rules from the slides:
//  1. returns the error interface (already the case here),
//  2. returns a literal nil when there is no error,
//  3. never smuggles a typed nil pointer through the interface.
func fixedValidate(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age"}
	}
	return nil
}

func main() {
	fmt.Println("--- 1 & 2: run the buggy version with age = 20 ---")
	err := buggyValidate(20)
	// The if is NOT entered for age=20, so `err` stays a nil
	// *ValidationError. That nil pointer is still stored in the error
	// interface, so the interface itself is NOT nil - it holds the pair
	// (*ValidationError, nil). That's why `err != nil` is true and the
	// "error:" branch runs.
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Println("error:", err) prints "<nil>" here (not a panic and not
	// the real message) because fmt detects that the dynamic value is a
	// nil pointer of a type with pointer methods and special-cases it
	// instead of calling Error() on a nil receiver.
	fmt.Println("err == nil:", err == nil) // false: interface holds a type
	fmt.Printf("%%T of err: %T\n", err)    // *main.ValidationError

	fmt.Println("\n--- 3: fixed version ---")
	fixedErr := fixedValidate(20)
	fmt.Println("age=20, err == nil:", fixedErr == nil) // true now

	fixedErr = fixedValidate(-1)
	fmt.Println("age=-1, err:", fixedErr)

	var ve *ValidationError
	if errors.As(fixedErr, &ve) {
		fmt.Println("extracted field via errors.As:", ve.Field)
	}
}
