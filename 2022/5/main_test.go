package main

import "testing"

const expected1 = "VJSFHWGFT"
const expected2 = "LCTQFBVZV"

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
