package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed example.in
var input string

// check panics if error is not nil.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// add adds given integers.
func add(a, b int) int {
	return a + b
}

// multiply multiplies given integers.
func multiply(a, b int) int {
	return a * b
}

// pow raises a to the power of b.
func pow(a, b int) int {
	if b == 0 {
		return 1
	}
	if b == 1 {
		return a
	}
	y := pow(a, b/2)
	if b%2 == 0 {
		return y * y
	}
	return a * y * y
}

// gcd finds the greatest common divisor.
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// lcm finds the least common multiple.
func lcm(n ...int) int {
	if len(n) < 2 {
		return n[0]
	}
	t := n[0] * n[1] / gcd(n[0], n[1])
	for _, i := range n[2:] {
		t = lcm(t, i)
	}
	return t
}

// Monkey holds items and can throw them to others based on his decision function.
type Monkey struct {
	id int
	// items contains worry levels for each item held.
	items []int
	// worryIncrease describes how worry level increases after inspection.
	worryIncrease func(oldWorry int) int
	// throwRecipient returns the monkey index to throw the item to.
	throwRecipient func(worryLevel int) int
	// itemsInspected counts the number of inspections by monkey.
	itemsInspected int

	throwDivisor int
}

// SpotMonkey parses a monkey desc into a monkey struct.
func SpotMonkey(monkeyDesc string) *Monkey {
	// Use strings.NewReplacer and fmt.Sscanf
	monkey := &Monkey{}
	replacer := strings.NewReplacer(", ", ",", "* old", "^ 2")
	var itemArr, opSign string
	var opConstant, catcher1, catcher2 int
	_, err := fmt.Sscanf(replacer.Replace(monkeyDesc),
		`Monkey %d:
		Starting items: %s
		Operation: new = old %s %d
		Test: divisible by %d
		 If true: throw to monkey %d
		 If false: throw to monkey %d`,
		&monkey.id, &itemArr, &opSign, &opConstant, &monkey.throwDivisor, &catcher1, &catcher2,
	)
	check(err)

	// Parse integers from string items
	items := strings.Split(itemArr, ",")
	monkey.items = make([]int, len(items))
	for i := range items {
		monkey.items[i], err = strconv.Atoi(items[i])
		check(err)
	}

	// Worry level increase after each inspection
	var mathOp func(a, b int) int
	switch opSign {
	case "+":
		mathOp = add
	case "*":
		mathOp = multiply
	case "^":
		mathOp = pow
	default:
		panic("Unexpected sign [" + opSign + "] for worry increase function")
	}
	monkey.worryIncrease = func(worryLevel int) int {
		return mathOp(worryLevel, opConstant)
	}

	// Deciding function for which monkey to throw the item to (index)
	monkey.throwRecipient = func(worryLevel int) int {
		if worryLevel%monkey.throwDivisor == 0 {
			return catcher1
		}
		return catcher2
	}

	return monkey
}

// InspectItem causes monkey to inspect the first item in its possesion increasing worry level for that item.
// reducer keeps the worry level at normal levels.
func (m *Monkey) InspectItem(reducer func(worryLevel int) int) {
	if len(m.items) == 0 {
		return
	}
	m.itemsInspected++
	m.items[0] = reducer(m.worryIncrease(m.items[0]))
}

// ThrowFirstItem throws monkey's first item possesed to given recipient.
func (m *Monkey) ThrowFirstItem(recipient *Monkey) {
	if len(m.items) == 0 {
		return
	}
	throwedItem := m.items[0]
	m.items = m.items[1:]
	recipient.items = append(recipient.items, throwedItem)
}

// ParseInput takes an input string (monkey description) and converts it to objects.
func ParseInput(inputDesc string) ([]*Monkey, []int) {
	descriptions := strings.Split(input, "\n\n")
	monkeys := make([]*Monkey, len(descriptions))
	divisors := make([]int, len(descriptions))
	for _, desc := range descriptions {
		monkey := SpotMonkey(desc)
		monkeys[monkey.id] = monkey
		divisors[monkey.id] = monkey.throwDivisor
	}
	return monkeys, divisors
}

// runChallenge returns the desired output for the days challenge.
// May print additional information to stdout.
func runChallenge(challengePart int) int {
	monkeys, divisors := ParseInput(input)
	// Part 1 related
	reduceFunc := func(worryLevel int) int {
		return worryLevel / 3
	}
	rounds := 20
	if challengePart == 2 {
		monkeyLCM := lcm(divisors...)
		reduceFunc = func(worryLevel int) int {
			return worryLevel % monkeyLCM
		}
		rounds = 10000
	}

	for i := 0; i < rounds; i++ {
		for _, monkey := range monkeys {
			for range monkey.items {
				monkey.InspectItem(reduceFunc)
				recipient := monkeys[monkey.throwRecipient(monkey.items[0])]
				monkey.ThrowFirstItem(recipient)
			}
		}
	}
	sort.Slice(monkeys, func(i, j int) bool { return monkeys[i].itemsInspected > monkeys[j].itemsInspected })
	return monkeys[0].itemsInspected * monkeys[1].itemsInspected
}

func main() {
	fmt.Printf("==== Part 1 ==== \n%d\n", runChallenge(1))
	fmt.Printf("==== Part 2 ==== \n%d\n", runChallenge(2))
}
