package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed challenge.in
var input string

func add(a, b int) int {
	return a + b
}
func multiply(a, b int) int {
	return a * b
}

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
//
//	Monkey 1:
//	Starting items: 54, 65, 75, 74
//	Operation: new = old + 6
//	Test: divisible by 19
//	  If true: throw to monkey 2
//	  If false: throw to monkey 0
func SpotMonkey(monkeyDesc string) *Monkey {
	monkey := &Monkey{}
	lines := strings.Split(monkeyDesc, "\n")

	// id parsing
	monkey.id, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(lines[0], "Monkey ", ""), ":", ""))
	// Items parsing (csv into list of int)
	items := strings.Split(strings.Split(lines[1], ":")[1], ",")
	monkey.items = make([]int, len(items))
	var err error
	for i := 0; i < len(items); i++ {
		a, err := strconv.Atoi(strings.TrimSpace(items[i]))
		monkey.items[i] = int(a)
		if err != nil {
			panic(err)
		}
	}

	// throw recipient calculation
	b, _ := strconv.Atoi(strings.Fields(lines[3])[3])
	monkey.throwDivisor = int(b)
	line4 := strings.Fields(lines[4])
	line5 := strings.Fields(lines[5])
	monkeyA, _ := strconv.Atoi(line4[len(line4)-1])
	monkeyB, _ := strconv.Atoi(line5[len(line5)-1])
	monkey.throwRecipient = func(worryLevel int) int {
		if worryLevel%monkey.throwDivisor == 0 {
			return monkeyA
		}
		return monkeyB
	}

	// worry increase function
	opDesc := strings.Fields(strings.Split(lines[2], "=")[1])
	if len(opDesc) < 3 || opDesc[0] != "old" {
		panic(fmt.Errorf("[Error] Error with worry increase operation (got %v) ", opDesc))
	}
	usesOldWorry := false
	a, err := strconv.Atoi(opDesc[2])
	increaseConstant := int(a)
	if err != nil {
		usesOldWorry = true
	}
	switch opDesc[1] {
	case "+":
		monkey.worryIncrease = func(oldWorry int) int {
			if usesOldWorry {
				return add(oldWorry, oldWorry)
			}
			return add(oldWorry, increaseConstant)
		}
	case "*":
		monkey.worryIncrease = func(oldWorry int) int {
			if usesOldWorry {
				return multiply(oldWorry, oldWorry)
			}
			return multiply(oldWorry, increaseConstant)
		}
	default:
		panic(fmt.Errorf("[Error] Worry increase function has different sign: [%s]", opDesc[1]))
	}

	return monkey
}

// InspectItem causes monkey to inspect the first item in its possesion increasing worry level for that item.
func (m *Monkey) InspectItem() {
	if len(m.items) == 0 {
		return
	}
	m.itemsInspected++
	m.items[0] = m.worryIncrease(m.items[0]) % 9699690
}

// ThrowFirstItem throws monkeys first item possesed to given recipient
func (m *Monkey) ThrowFirstItem(recipient *Monkey) {
	if len(m.items) == 0 {
		return
	}
	throwedItem := m.items[0]
	m.items = m.items[1:]
	recipient.items = append(recipient.items, throwedItem)
}

func ParseInput(inputDesc string) []*Monkey {
	descriptions := strings.Split(input, "\n\n")
	monkeys := make([]*Monkey, len(descriptions))

	for _, desc := range descriptions {
		monkey := SpotMonkey(desc)
		monkeys[monkey.id] = monkey
	}
	return monkeys
}

func runChallenge(challengePart int) int {
	result := -1
	monkeys := ParseInput(input)
	rounds := 10000
	for i := 0; i < rounds; i++ {
		fmt.Printf(">>>>>>> ROUND %d <<<<<<<\n", i)
		for _, monkey := range monkeys {
			fmt.Printf("====== Monkey %d ======\n", monkey.id)
			for range monkey.items {
				fmt.Printf("  worry change: %d ->", monkey.items[0])
				monkey.InspectItem()
				fmt.Printf("%d\n", monkey.items[0])
				fmt.Printf("  Throwing [%d] to [Monkey %d]\n", monkey.items[0], monkey.throwRecipient(monkey.items[0]))
				recipient := monkeys[monkey.throwRecipient(monkey.items[0])]
				monkey.ThrowFirstItem(recipient)
			}

		}
	}

	sort.Slice(monkeys, func(i, j int) bool { return monkeys[i].itemsInspected > monkeys[j].itemsInspected })
	for _, m := range monkeys {
		fmt.Printf("[Monkey %d] \n  Inspections: %d \n", m.id, m.itemsInspected)
	}
	result = monkeys[0].itemsInspected * monkeys[1].itemsInspected
	return result
}

func main() {
	fmt.Println(runChallenge(1))
}
