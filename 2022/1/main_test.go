package main

import "testing"

const expected1 = 67027
const expected2 = 197291

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
