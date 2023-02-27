package main

import (
	"testing"
)

func BenchmarkFindMaxSizeColorGroup(b *testing.B) {
	f := generateField(5,5,2)
	for i := 0; i < b.N; i++ {
		f.findMaxSizeColorGroup()
	}
}