package main

import (
	"errors"
	"fmt"
)

// Task 2: Describe anything

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1f°C", float64(c)) }

func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int, int64:
		return fmt.Sprintf("integer %v", x)
	case string:
		return fmt.Sprintf("string of %d bytes", len(x))
	case []int:
		return fmt.Sprintf("slice with %d ints", len(x))
	case bool:
		return fmt.Sprintf("bool %v", x)
	case map[string]int:
		return fmt.Sprintf("map with %d keys", len(x))
	case error:
		return "error: " + x.Error()
	case fmt.Stringer:
		return "Stringer: " + x.String()
	default:
		return fmt.Sprintf("unknown type %T", x)
	}
}

// Participation Task (Part 2): rewrite describe without a type switch,
// using only "comma ok" assertions.
func describeNoSwitch(v any) string {
	if v == nil {
		return "nil"
	}
	if n, ok := v.(int); ok {
		return fmt.Sprintf("integer %v", n)
	}
	if n, ok := v.(int64); ok {
		return fmt.Sprintf("integer %v", n)
	}
	if s, ok := v.(string); ok {
		return fmt.Sprintf("string of %d bytes", len(s))
	}
	if sl, ok := v.([]int); ok {
		return fmt.Sprintf("slice with %d ints", len(sl))
	}
	if b, ok := v.(bool); ok {
		return fmt.Sprintf("bool %v", b)
	}
	if m, ok := v.(map[string]int); ok {
		return fmt.Sprintf("map with %d keys", len(m))
	}
	if err, ok := v.(error); ok {
		return "error: " + err.Error()
	}
	if s, ok := v.(fmt.Stringer); ok {
		return "Stringer: " + s.String()
	}
	return fmt.Sprintf("unknown type %T", v)
}

func main() {
	values := []struct {
		v        any
		expected string
	}{
		// 1. nil -> hits the `case nil` branch directly.
		{nil, "nil"},
		// 2. plain int literal -> matches `case int, int64`.
		{7, "integer 7"},
		// 3. explicit int64 -> same branch as above, x keeps type any
		//    (multi-type case), so %v just prints the number.
		{int64(7), "integer 7"},
		// 4. "привет" is 6 Cyrillic runes, but each one is 2 bytes in
		//    UTF-8, so len() (byte length) is 12, not 6 (see Week 3).
		{"привет", "string of 12 bytes"},
		{[]int{1, 2, 3}, "slice with 3 ints"},
		// errors.New returns a value whose type implements `error`.
		{errors.New("disk full"), "error: disk full"},
		// Celsius has a String() method -> matches `case fmt.Stringer`.
		{Celsius(21.5), "Stringer: 21.5°C"},
		// a bare float64 has no matching case above `default`.
		{3.14, "unknown type float64"},
		{struct{}{}, "unknown type struct {}"},
		{true, "bool true"},
		{map[string]int{"a": 1, "b": 2}, "map with 2 keys"},
	}

	fmt.Println("=== Task 2: describe (type switch) ===")
	for _, tc := range values {
		got := describe(tc.v)
		fmt.Printf("describe(%#v) = %-28s // expected: %s\n", tc.v, got, tc.expected)
	}

	fmt.Println("\n=== Participation Task: describeNoSwitch (comma-ok only) ===")
	for _, tc := range values {
		fmt.Printf("describeNoSwitch(%#v) = %s\n", tc.v, describeNoSwitch(tc.v))
	}

	// Which version is easier to read?
	// The type switch (Task 2) is easier: one `switch` block makes the
	// set of handled types and their order obvious at a glance, and the
	// compiler gives each case a correctly narrowed variable for free.
	// describeNoSwitch needs a repeated `if _, ok := ...; ok { ... }`
	// block per type, duplicates the "ok" check every time, and it's
	// easy to forget an early `return` and fall through into the wrong
	// branch.
}
