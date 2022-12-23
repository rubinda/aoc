package main

import (
	"testing"
)

const (
	expected1 = 6032
	expected2 = 10006
)

func TestChallenge1(t *testing.T) {
	actual := runChallenge(1)

	if actual != expected1 {
		t.Errorf("Wrong result! Expected: %v, actual: %v", expected1, actual)
	}
}

func TestChallenge2(t *testing.T) {
	actual := runChallenge(2)

	if actual != expected2 {
		t.Errorf("Wrong result! Expected: %v, actual: %v", expected2, actual)
	}
}

func Benchmark1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runChallenge(1)
	}
}
func Benchmark2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runChallenge(2)
	}
}
