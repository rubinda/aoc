package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

var (
	//go:embed challenge.in
	input string
)

const (
	// Used in part 2 to increase problem difficulty
	decryptionKeyMultiplier = 811589153
	// How many times to mix the list in challenge part 1
	Mixings_1 = 1
	// How many times to mix the list in challenge part 2
	Mixings_2 = 10
)

// ListItem is a node in a double linked list.
type ListItem struct {
	Value int
	Next  *ListItem
	Prev  *ListItem
}

// DoublyLinkedList is a wrapper for a circular linked list with connections in both directions.
type DoublyLinkedList struct {
	Len  int
	Head *ListItem
	Tail *ListItem
	// ZeroItem is the only item in challenge with value 0
	ZeroItem *ListItem
}

func (dll *DoublyLinkedList) NormalizeSteps(steps int) int {
	return steps % (dll.Len - 1)
}

// Move moves the given item for given steps.
// Negative steps mean move through "previous", positive steps mean move through "next".
func (dll *DoublyLinkedList) Move(item *ListItem, steps int) {
	// Example also works without the -1, took me some time to double check this :^)
	steps = dll.NormalizeSteps(steps)
	if steps == 0 {
		return
	}
	moveNext := true
	if steps < 0 {
		moveNext = false
	}
	current := item
	for steps != 0 {
		if moveNext {
			current = current.Next
			steps--
		} else {
			current = current.Prev
			steps++
		}
	}
	// Remove connections from old position
	item.Prev.Next = item.Next
	item.Next.Prev = item.Prev

	if !moveNext {
		// If negative direction, make 1 more negative step so we can replace the future component
		// ... combines the replace logic below for positive and negative stepping
		current = current.Prev
	}
	// Place moved item as "current" -> "item" -> "current.Next"
	item.Prev = current
	item.Next = current.Next
	current.Next.Prev = item
	current.Next = item
}

// String converts the linked list to textual format.
// The output starts with the original "head" (because it's circular indexing is ignored matter).
func (dll *DoublyLinkedList) String() string {
	out := ""
	current := dll.Head
	for i := 0; i < dll.Len; i++ {
		out += fmt.Sprintf("%d ", current.Value)
		current = current.Next
	}
	return out
}

// CreateLinkedList creates a circular double linked list.
// Returns the linked list and a slice containing the original indexing of list items.
func CreateLinkedList(numbers []int) (*DoublyLinkedList, []*ListItem) {
	originalPositions := make([]*ListItem, len(numbers))
	linkedList := &DoublyLinkedList{Len: len(numbers)}
	linkedList.Head = &ListItem{Value: numbers[0]}
	current := linkedList.Head
	for i := 0; i < len(numbers); i++ {
		if numbers[i] == 0 {
			// Store item with value 0, because challenge wants to operate on it later
			linkedList.ZeroItem = current
		}
		if i == len(numbers)-1 {
			// Tail is special case, because "Next" is overflow in a slice.
			linkedList.Tail = current
		} else {
			// Build list by constructing "Next" and settings it's "Previous" to "current"
			current.Next = &ListItem{Value: numbers[i+1]}
			current.Next.Prev = current
		}
		originalPositions[i] = current
		current = current.Next
	}
	// Set circular connections (last to first and vice versa)
	linkedList.Head.Prev = linkedList.Tail
	linkedList.Tail.Next = linkedList.Head
	return linkedList, originalPositions
}

// parseInput takes the challenge input and converts it to a list of integers.
func parseInput(encrypted string) []int {
	numbersStr := strings.Split(encrypted, "\n")
	numbers := make([]int, len(numbersStr))
	var err error
	for n := range numbersStr {
		numbers[n], err = strconv.Atoi(numbersStr[n])
		if err != nil {
			panic(err)
		}
	}
	return numbers
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	numbers := parseInput(input)
	mixings := Mixings_1
	if challengePart == 2 {
		// Part 2 has more runs and LARGER numbers (which we normalize later)
		mixings = Mixings_2
		for i := range numbers {
			numbers[i] *= decryptionKeyMultiplier
		}
	}
	linked, originalPositions := CreateLinkedList(numbers)
	for i := 0; i < mixings; i++ {
		for i := range originalPositions {
			item := originalPositions[i]
			linked.Move(item, item.Value)
		}
	}
	// The challenge output is the sum of 3 items which are [1000,2000,3000] steps away from ZeroItem
	groveCoordinatesSum := 0
	steps := 3000
	item := linked.ZeroItem
	for steps > 0 {
		item = item.Next
		steps--
		if steps == 2000 || steps == 1000 || steps == 0 {
			groveCoordinatesSum += item.Value
		}
	}
	return groveCoordinatesSum
}

func main() {
	fmt.Println(runChallenge(2))
}
