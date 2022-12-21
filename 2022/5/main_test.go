package main

import "testing"

const expected1 = "CMZ"
const expected2 = "MCD"

func TestChallenge1(t *testing.T) {
	actual := runChallenge(1)

	if actual != expected1 {
		t.Errorf("Wrong result! Expected: %s, actual: %s", expected1, actual)
	}
}

func TestChallenge2(t *testing.T) {
	actual := runChallenge(2)

	if actual != expected2 {
		t.Errorf("Wrong result! Expected: %s, actual: %s", expected2, actual)
	}
}

func Benchmark2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runChallenge(2)
	}
}
