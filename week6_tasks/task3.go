package main

import (
	"errors"
	"fmt"
)

// Task 3: Grade book

var ErrNotFound = errors.New("student not found")

type GradeError struct {
	Student string
	Grade   int
}

func (e *GradeError) Error() string {
	return fmt.Sprintf("invalid grade %d for %s", e.Grade, e.Student)
}

type GradeBook struct {
	grades map[string]int
}

func NewGradeBook() *GradeBook {
	return &GradeBook{grades: make(map[string]int)}
}

func (gb *GradeBook) Get(name string) (int, error) {
	g, ok := gb.grades[name]
	if !ok {
		return 0, ErrNotFound
	}
	return g, nil
}

func (gb *GradeBook) Set(name string, grade int) error {
	if grade < 0 || grade > 100 {
		return &GradeError{Student: name, Grade: grade}
	}
	gb.grades[name] = grade
	return nil
}

func AddBonus(gb *GradeBook, name string, bonus int) error {
	g, err := gb.Get(name)
	if err != nil {
		return fmt.Errorf("add bonus to %q: get grade: %w", name, err)
	}
	if err := gb.Set(name, g+bonus); err != nil {
		return fmt.Errorf("add bonus to %q: set grade: %w", name, err)
	}
	return nil
}

// AddBonusBroken is the same as AddBonus but uses %v for the Get error.
func AddBonusBroken(gb *GradeBook, name string, bonus int) error {
	g, err := gb.Get(name)
	if err != nil {
		return fmt.Errorf("add bonus to %q: get grade: %v", name, err)
	}
	if err := gb.Set(name, g+bonus); err != nil {
		return fmt.Errorf("add bonus to %q: set grade: %w", name, err)
	}
	return nil
}

func report(err error) {
	var ge *GradeError
	switch {
	case err == nil:
		fmt.Println("OK")
	case errors.Is(err, ErrNotFound):
		fmt.Println("404:", err)
	case errors.As(err, &ge):
		fmt.Printf("400: %v (student=%s grade=%d)\n", err, ge.Student, ge.Grade)
	default:
		fmt.Println("500:", err)
	}
}

func main() {
	gb := NewGradeBook()
	_ = gb.Set("Aruzhan", 95)
	_ = gb.Set("Dias", 60)

	report(AddBonus(gb, "Dias", 5))     // OK
	report(AddBonus(gb, "Timur", 5))    // 404
	report(AddBonus(gb, "Aruzhan", 10)) // 400 (105)

	fmt.Println("--- with %v instead of %w for Get ---")
	report(AddBonusBroken(gb, "Timur", 5)) // 500
	// Explanation: %v only copies the message text, the ErrNotFound value is
	// no longer in the chain, so errors.Is returns false and neither the
	// 404 nor the 400 branch matches: report falls into default (500).
}
