package main

import (
	"fmt"
	"time"
)

// Task 5: defer in practice

// ---- 1. the four "predict" snippets ----

func lifo() {
	defer fmt.Println("1: deferred first")
	defer fmt.Println("2: deferred second")
	defer fmt.Println("3: deferred third")
	fmt.Println("main body")
	// main body, 3, 2, 1
}

func argsEvaluated() {
	x := 1
	defer fmt.Println("deferred x =", x)             // x is evaluated now: 1
	defer func() { fmt.Println("closure x =", x) }() // closure sees the variable: 2
	x = 2
	fmt.Println("current x =", x)
	// current x = 2, closure x = 2, deferred x = 1
}

func double() (result int) {
	defer func() { result *= 2 }()
	return 21 // sets result = 21, then the defer doubles it: 42
}

func notDouble() int {
	result := 21
	defer func() { result *= 2 }()
	return result // the value 21 is already copied to the return slot: 21
}

func trace(name string) func() {
	fmt.Println("enter", name)
	return func() { fmt.Println("exit", name) }
}

func work() {
	// trace("work") is evaluated NOW (prints "enter work"),
	// the returned func is deferred and runs at the end.
	defer trace("work")()
	fmt.Println("working...")
}

func helper() {
	if r := recover(); r != nil {
		fmt.Println("recovered:", r)
	}
}

func a() {
	defer func() {
		helper() // NOT called directly by the deferred func -> recover returns nil
	}()
	panic("boom A")
}

func b() {
	defer helper() // helper itself is the deferred function -> recovers
	panic("boom B")
}

func c() {
	fmt.Println("c: recover outside defer returned", recover()) // <nil>
}

// ---- 2. defer in a loop ----

type Resource struct{ Name string }

func Open(name string) (*Resource, error) {
	fmt.Println("open ", name)
	return &Resource{Name: name}, nil
}

func (r *Resource) Close() {
	fmt.Println("close", r.Name)
}

func processAllBad(names []string) error {
	for _, name := range names {
		r, err := Open(name)
		if err != nil {
			return err
		}
		defer r.Close() // runs only when processAllBad returns
		fmt.Println("use  ", r.Name)
	}
	return nil
}

func processOne(name string) error {
	r, err := Open(name)
	if err != nil {
		return err
	}
	defer r.Close() // runs when processOne returns, i.e. after each file
	fmt.Println("use  ", r.Name)
	return nil
}

func processAllGood(names []string) error {
	for _, name := range names {
		if err := processOne(name); err != nil {
			return err
		}
	}
	return nil
}

// ---- 3. timer ----

func timer(name string) func() {
	start := time.Now()
	return func() { fmt.Printf("%s took %v\n", name, time.Since(start)) }
}

func sortSomething() {
	defer timer("sort")()
	time.Sleep(20 * time.Millisecond)
}

func main() {
	fmt.Println("== LIFO ==")
	lifo()
	fmt.Println("== arguments ==")
	argsEvaluated()
	fmt.Println("== named result ==")
	fmt.Println(double(), notDouble())
	fmt.Println("== trace ==")
	work()
	fmt.Println("== recover ==")
	a2 := func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("a panicked through (not recovered by helper):", r)
			}
		}()
		a()
	}
	a2()
	b()
	c()

	names := []string{"a.txt", "b.txt", "c.txt"}
	fmt.Println("== processAllBad ==")
	_ = processAllBad(names)
	fmt.Println("== processAllGood ==")
	_ = processAllGood(names)

	fmt.Println("== timer ==")
	sortSomething()
}
