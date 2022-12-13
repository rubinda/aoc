package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Represents compare function results.
const (
	LessThan    int = -1
	Equal       int = 0
	GreaterThan int = 1
)

const (
	specialPacket1 = "[[2]]"
	specialPacket2 = "[[6]]"
)

// comparisonString converts numeric comparison result to text.
var comparisonString = map[int]string{
	Equal:       "CONT",
	LessThan:    "PASS",
	GreaterThan: "FAIL",
}

// check panics if error is not nil
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// messageItem represents an item in a message (hack to get array of strings from JSON unmarshal regardless of item type).
type messageItem string

// UnmarshalJSON provides a method to convert JSON data to messgeItem type.
func (m *messageItem) UnmarshalJSON(data []byte) error {
	if n := len(data); n > 1 && data[0] == '"' && data[n-1] == '"' {
		return json.Unmarshal(data, (*string)(m))
	}
	*m = messageItem(data)
	return nil
}

func (m messageItem) containsArray() bool {
	return string(m[0]) == "["
}

// compareAsInt tries to convert message items to integers before comparison.
func compareAsInt(left, right messageItem) int {
	aI, err := strconv.Atoi(string(left))
	check(err)
	bI, err := strconv.Atoi(string(right))
	check(err)
	if aI > bI {
		return GreaterThan
	} else if aI < bI {
		return LessThan
	}
	return Equal
}

// arrayFrom parses one level of an array described in a valid JSON string
// e.g. "[[1, 2, 3], 0, [4, 5]]" => ["[1, 2, 3]", "0", "[4, 5]"]
// e.g. "[[[4], 5, 6], 9]" => ["[[4], 5, 6]", "9"]
func arrayFrom(arrayDesc string) []messageItem {
	if _, err := strconv.Atoi(arrayDesc); err == nil {
		mi := messageItem(arrayDesc)
		return []messageItem{mi}
	}
	var a []messageItem
	err := json.Unmarshal([]byte(arrayDesc), &a)
	if err != nil {
		panic(err)
	}
	return a
}

// Compare returns the order of 2 packets from the challenge input.
func Compare(packet1, packet2 string) int {
	left := arrayFrom(packet1)
	right := arrayFrom(packet2)
	// Current item in left and right array
	i := 0
	for i < len(left) && i < len(right) {
		// fmt.Printf(">Comparing {%v} vs {%v} \n", left[i], right[i])
		if left[i].containsArray() || right[i].containsArray() {
			// fmt.Printf(">>Going recursive on %v and %v \n", left[i], right[i])
			result := Compare(string(left[i]), string(right[i]))
			// fmt.Printf(">>  %v \n", comparisonString[result])
			if result != Equal {
				// fmt.Println(">>>Returning from recursive ", result, EQUAL)
				return result
			}
			// fmt.Println(">>>Recursion says continue ", left[i], right[i], len(left), len(right))
		} else {
			if r := compareAsInt(left[i], right[i]); r != Equal {
				// fmt.Printf("[%s] because %s", comparisonString[r], string(left[i]))
				return r
			}
		}
		i++
	}
	// No comparison between characters was able to resolve order and both reached end of list
	if i == len(left) && i == len(right) {
		return Equal
	}
	// Left reached end of list (but right didn't beacuse it would have already returned)
	if i == len(left) {
		return LessThan
	}
	// Right reached end of list and left didn't
	return GreaterThan
}

//go:embed challenge.in
var input string

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	packetPairs := strings.Split(input, "\n\n")
	packets := make([]string, len(packetPairs)*2)

	alreadySortedPairs := 0
	for i, doubleLine := range packetPairs {
		pair := strings.Split(doubleLine, "\n")
		// fmt.Printf("==== Pair %d ==== \n", i+1)
		// fmt.Printf("0: %v \n", pair[0])
		// fmt.Printf("1: %v \n", pair[1])
		if challengePart == 1 {
			result := Compare(pair[0], pair[1])
			if result == LessThan {
				alreadySortedPairs += i + 1
			}
		}
		packets[i*2] = pair[0]
		packets[i*2+1] = pair[1]
	}
	if challengePart == 1 {
		return alreadySortedPairs
	}
	if challengePart == 2 {
		packets = append(packets, specialPacket1, specialPacket2)
		sort.Slice(packets, func(i, j int) bool {
			return Compare(packets[i], packets[j]) == LessThan
		})
		// Find positions of special packets (starting count with 1)
		packet1 := -1
		packet2 := -1
		// fmt.Printf("==== SORTED PACKETS (ascending) ====\n")
		for i, s := range packets {
			if s == specialPacket1 {
				packet1 = i + 1
			} else if s == specialPacket2 {
				packet2 = i + 1
				break
			}
			// fmt.Println(s)
		}
		return packet1 * packet2
	}
	return -1
}

func main() {
	fmt.Println(runChallenge(1))
}
