package main

import "testing"

func BenchmarkReport_FirstValue(b *testing.B) {
	r := Report{}
	b.ReportAllocs()
	for b.Loop() {
		r.FirstValue()
	}
}

func BenchmarkReport_FirstPointer(b *testing.B) {
	r := &Report{}
	b.ReportAllocs()
	for b.Loop() {
		r.FirstPointer()
	}
}

func BenchmarkPoint_FirstValue(b *testing.B) {
	p := Point{}
	b.ReportAllocs()
	for b.Loop() {
		p.FirstValue()
	}
}

func BenchmarkPoint_FirstPointer(b *testing.B) {
	p := &Point{}
	b.ReportAllocs()
	for b.Loop() {
		p.FirstPointer()
	}
}
