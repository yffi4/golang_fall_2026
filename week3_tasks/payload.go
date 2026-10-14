package main

import (
	"bytes"
	"strings"
)

// processPayloadOriginal is the version from the task description.
// It converts data ([]byte) to a string just to run TrimSpace, then runs
// ReplaceAll on that string, and returns a string. Both conversions
// ([]byte -> string for the trim, and the implicit string result) force
// Go to allocate a brand-new copy of the underlying bytes, because
// strings are immutable: the runtime can never hand out a string that
// aliases a mutable []byte, so every conversion copies.
func processPayloadOriginal(data []byte) string {
	// Unnecessary type conversions
	trimmed := string(bytes.TrimSpace(data)) // allocation #1: []byte -> string copy
	clean := strings.ReplaceAll(trimmed, "\r", "")
	return clean
}

// processPayloadOptimized works on the []byte the whole way through and
// does the trim + "\r" removal in a single pass instead of two separate
// allocating steps (string conversion, then ReplaceAll).
//
// It uses strings.Builder rather than a plain []byte/bytes.Buffer for the
// output because Builder.String() is a zero-copy unsafe cast of its
// internal buffer - unlike bytes.Buffer.String() or a []byte->string
// conversion, it does not allocate a second time. Combined with
// sb.Grow(len(trimmed)) up front, the whole function performs exactly one
// allocation (for the builder's backing array) versus two in the
// original.
func processPayloadOptimized(data []byte) string {
	trimmed := bytes.TrimSpace(data) // no allocation: TrimSpace returns a subslice

	var sb strings.Builder
	sb.Grow(len(trimmed)) // one allocation, sized exactly, no reallocation while writing
	for _, c := range trimmed {
		if c != '\r' {
			sb.WriteByte(c)
		}
	}
	return sb.String() // zero-copy: no second allocation
}
