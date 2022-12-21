package main

import (
	_ "embed"
	"fmt"
)

const (
	startMarkerLen   = 4
	messageMarkerLen = 14
)

//go:embed example.in
var input string

func isUniqueChars(str string) bool {
	charMap := make(map[rune]bool, len(str))

	for _, r := range str {
		if _, ok := charMap[r]; ok {
			return false
		}
		charMap[r] = true
	}
	return true
}

func runChallenge(challengePart int) int {
	inputLen := len(input)
	bufferLen := startMarkerLen
	if challengePart == 2 {
		bufferLen = messageMarkerLen
	}
	for i := range input {
		if (i + bufferLen) >= inputLen {
			break
		}
		buffer := input[i : i+bufferLen]
		if isUniqueChars(buffer) {
			// fmt.Printf("'%s' @ {%d:%d} is unique \n", buffer, i, i+bufferLen)
			// Return the number of characters processed
			return i + bufferLen
		}
	}
	return -1
}

func main() {
	fmt.Println(runChallenge(1))
}
