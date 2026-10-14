package main

import (
	"fmt"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------
// Basic map operations (read/update/delete) + two-value lookup
// ---------------------------------------------------------------------

func demoMapOperations() {
	fmt.Println("\n=== Basic map operations ===")

	m := map[string]int{"one": 1, "two": 2}
	fmt.Println("read m[\"one\"]:", m["one"])

	m["two"] = 20
	fmt.Println("after update m[\"two\"] = 20:", m["two"])

	delete(m, "one")
	_, stillThere := m["one"]
	fmt.Println("after delete, \"one\" present?", stillThere)
}

// Task: two-value assignment to distinguish "missing" from "present with
// zero value" - answers the "Bob scored 0 vs Charlie never took the exam"
// problem.
func demoBobCharlieCheck() {
	fmt.Println("\n=== Two-value map lookup (Bob vs Charlie) ===")

	scores := map[string]int{
		"Alice": 90,
		"Bob":   0, // Bob took the exam but scored 0
	}

	if _, ok := scores["Bob"]; ok {
		fmt.Println("Bob is in the map")
	} else {
		fmt.Println("Bob is missing")
	}

	if _, ok := scores["Charlie"]; ok {
		fmt.Println("Charlie is in the map")
	} else {
		fmt.Println("Charlie is missing")
	}
}

// ---------------------------------------------------------------------
// Inventory management task
// ---------------------------------------------------------------------

func demoInventory() {
	fmt.Println("\n=== Inventory management ===")

	// 1. Initialize
	inventory := map[string]int{
		"apples":  10,
		"bananas": 5,
	}

	// 2. Read
	fmt.Println("current stock of apples:", inventory["apples"])

	// 3. Update
	inventory["bananas"] = 12

	// 4. Insert
	inventory["oranges"] = 8

	// 5. Delete
	delete(inventory, "apples")

	// 6. Print (sorted keys for deterministic output, since map iteration
	// order in Go is randomized)
	printSortedMap("updated inventory", inventory)
}

func printSortedMap(label string, m map[string]int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Printf("%s: {", label)
	for i, k := range keys {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%q: %d", k, m[k])
	}
	fmt.Println("}")
}

// ---------------------------------------------------------------------
// Tour of Go -> Exercise: Maps (WordCount)
// https://go.dev/tour/moretypes/23
// ---------------------------------------------------------------------

// WordCount splits s on whitespace and returns a map of word -> count of
// occurrences.
func WordCount(s string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(s) {
		counts[word]++
	}
	return counts
}

func demoWordCount() {
	fmt.Println("\n=== Tour of Go: Exercise Maps (WordCount) ===")

	inputs := []string{
		"I ate a donut. Then I ate another donut.",
		"the quick brown fox jumped over the lazy dog the fox ran",
	}

	for _, s := range inputs {
		fmt.Printf("WordCount(%q) =\n", s)
		counts := WordCount(s)

		words := make([]string, 0, len(counts))
		for w := range counts {
			words = append(words, w)
		}
		sort.Strings(words)

		for _, w := range words {
			fmt.Printf("  %q: %d\n", w, counts[w])
		}
	}
}
func some(s string) map[string]int {
	m := make(map[string]int)
	for _, word := range strings.Fields(s) {
		m[word]++
	}
}
