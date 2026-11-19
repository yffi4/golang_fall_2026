package main

import (
	"fmt"
	"os"
)

// Bonus: os.Exit
//
// Only "hello" is printed. os.Exit terminates the process immediately,
// without unwinding the stack, so deferred calls (the "bye!" line) never run.
// log.Fatal calls os.Exit(1) after printing, so defers (closing files,
// unlocking, flushing buffers) are skipped too. Use log.Fatal only in main
// (or not at all) and prefer returning errors to main.

func main() {
	defer fmt.Println("bye!")
	fmt.Println("hello")
	os.Exit(3)
}
