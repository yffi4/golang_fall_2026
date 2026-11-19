package main

import (
	"errors"
	"fmt"
)

// Task 6: Safe task runner

type Task struct {
	Name string
	Run  func() error
}

func runOne(t Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("task %s: panic: %v", t.Name, r)
		}
	}()
	if e := t.Run(); e != nil {
		return fmt.Errorf("task %s: %w", t.Name, e)
	}
	return nil
}

func SafeRun(tasks []Task) []error {
	var errs []error
	for _, t := range tasks {
		if err := runOne(t); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

var ErrSentinel = errors.New("sentinel failure")

func main() {
	tasks := []Task{
		{"ok", func() error { return nil }},
		{"index", func() error {
			s := []int{1, 2, 3}
			i := 5
			_ = s[i]
			return nil
		}},
		{"nilmap", func() error {
			var m map[string]int
			m["x"] = 1
			return nil
		}},
		{"config", func() error { panic("config is missing") }},
		{"sentinel", func() error { return ErrSentinel }},
		{"last", func() error { fmt.Println("last task still ran"); return nil }},
	}

	errs := SafeRun(tasks)
	for _, e := range errs {
		fmt.Println(e)
	}
	joined := errors.Join(errs...)
	fmt.Println("sentinel found in joined:", errors.Is(joined, ErrSentinel)) // true
}
