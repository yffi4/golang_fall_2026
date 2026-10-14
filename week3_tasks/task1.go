package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ---------------------------------------------------------------------
// Strings / runes section
// ---------------------------------------------------------------------

// demoStringLen shows the difference between len() (bytes) and
// utf8.RuneCountInString() (actual characters).
func demoStringLen() {
	fmt.Println("=== len() vs utf8.RuneCountInString() ===")

	s := "hello"
	fmt.Printf("%q -> len=%d (ASCII, 1 byte per rune, so bytes == runes)\n", s, len(s))

	s2 := "世界"
	// Each Chinese character here is 3 bytes in UTF-8, so len() returns 6,
	// not 2. len() always counts bytes, never characters/runes.
	fmt.Printf("%q -> len=%d, RuneCountInString=%d\n", s2, len(s2), utf8.RuneCountInString(s2))
}

// demoCorrectIteration shows the wrong (byte-indexed) way to iterate a
// string containing multi-byte characters and the correct (range) way.
func demoCorrectIteration() {
	fmt.Println("\n=== Incorrect vs correct string iteration ===")

	s := "hêllo"

	fmt.Println("Incorrect (s[i] indexing, prints garbage for multi-byte runes):")
	for i := 0; i < len(s); i++ {
		fmt.Printf("  position %d: %c\n", i, s[i])
	}

	fmt.Println("Correct (for i, r := range s):")
	for i, r := range s {
		fmt.Printf("  position %d: %c\n", i, r)
	}
}

// secondRune returns the second rune of s in a readable way.
// Converting the string to a []rune slice is the simplest, most readable
// way to do random access by rune index (as opposed to byte index).
// It works by decoding the whole UTF-8 string into a slice of int32 code
// points; after that, ordinary slice indexing gives you rune N directly.
func secondRune(s string) (rune, bool) {
	runes := []rune(s)
	if len(runes) < 2 {
		return 0, false
	}
	return runes[1], true
}

func demoSecondRune() {
	fmt.Println("\n=== Accessing the second rune ===")
	s := "hêllo"
	if r, ok := secondRune(s); ok {
		fmt.Printf("second rune of %q is %q\n", s, r)
	}

	s = "世界杯"
	if r, ok := secondRune(s); ok {
		fmt.Printf("second rune of %q is %q\n", s, r)
	}
}

// ---------------------------------------------------------------------
// Trim functions section
// ---------------------------------------------------------------------

func demoTrimFunctions() {
	fmt.Println("\n=== Trim functions ===")

	s := "123oxo"
	fmt.Printf("TrimRight(%q, \"xo\")  = %q  (strips ANY trailing chars in the cutset \"xo\")\n",
		s, strings.TrimRight(s, "xo"))
	fmt.Printf("TrimSuffix(%q, \"xo\") = %q  (strips the exact literal suffix \"xo\" once)\n",
		s, strings.TrimSuffix(s, "xo"))

	p := "oxo123"
	fmt.Printf("TrimLeft(%q, \"xo\")   = %q  (strips ANY leading chars in the cutset \"xo\")\n",
		p, strings.TrimLeft(p, "xo"))
	fmt.Printf("TrimPrefix(%q, \"oxo\") = %q  (strips the exact literal prefix \"oxo\" once)\n",
		p, strings.TrimPrefix(p, "oxo"))
}

// ---------------------------------------------------------------------
// strings.Builder / Grow section
// ---------------------------------------------------------------------

// concatNaive is the "incorrect"/slow approach: every += allocates a new
// string, copying everything accumulated so far. O(n^2) in the worst case.
func concatNaive(values []string) string {
	s := ""
	for _, v := range values {
		s += v
	}
	return s
}

// concatBuilder uses strings.Builder without pre-sizing.
func concatBuilder(values []string) string {
	var sb strings.Builder
	for _, v := range values {
		sb.WriteString(v)
	}
	return sb.String()
}

// concatBuilderWithGrow uses strings.Builder and pre-allocates the
// underlying buffer with Grow, since the total length is known in
// advance. This avoids repeated reallocation/copying as the builder
// grows.
func concatBuilderWithGrow(values []string) string {
	total := 0
	for _, v := range values {
		total += len(v)
	}

	var sb strings.Builder
	sb.Grow(total) // pre-allocate exactly the bytes we need
	for _, v := range values {
		sb.WriteString(v)
	}
	return sb.String()
}

func demoBuilderGrow() {
	fmt.Println("\n=== strings.Builder with Grow ===")
	values := []string{"Go ", "is ", "great ", "for ", "backend ", "development!"}

	fmt.Println("naive (+=):        ", concatNaive(values))
	fmt.Println("Builder (no Grow): ", concatBuilder(values))
	fmt.Println("Builder + Grow:    ", concatBuilderWithGrow(values))
}

// ---------------------------------------------------------------------
// entry point
// ---------------------------------------------------------------------

func main() {
	demoStringLen()
	demoCorrectIteration()
	demoSecondRune()
	demoTrimFunctions()
	demoBuilderGrow()

	demoMapOperations()
	demoBobCharlieCheck()
	demoInventory()
	demoWordCount()
}
