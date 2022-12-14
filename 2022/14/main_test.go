package main

import (
	"testing"
)

const (
	expected1 = 728
	expected2 = 27623
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
