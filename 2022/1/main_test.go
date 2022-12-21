package main

import "testing"

const expected1 = 24000
const expected2 = 45000

func TestChallenge1(t *testing.T) {
	actual, _ := runChallenge(1)

	if actual != expected1 {
		t.Errorf("Wrong result! Expected: %d, actual: %d", expected1, actual)
	}
}

func TestChallenge2(t *testing.T) {
	actual, _ := runChallenge(2)

	if actual != expected2 {
		t.Errorf("Wrong result! Expected: %d, actual: %d", expected2, actual)
	}
}

func BenchmarkChallenge2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runChallenge(2)
	}
}
