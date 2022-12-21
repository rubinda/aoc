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

// Operations between two Monkey Number fields
const (
	operationAdd      = "+"
	operationSubract  = "-"
	operationDivide   = "/"
	operationMultiply = "*"
)

const (
	// The monkey whose number is wanted by the challenge.
	wantedMonkeyName = "root"
	// In challenge part 2 this monkey is a special case (it's the user yelling a calculated number)
	myMonkeyName = "humn"
)

// doOperation executes a numeric operation between two integers.
func doOperation(a, b int, opSign string) int {
	switch opSign {
	case operationAdd:
		return a + b
	case operationSubract:
		return a - b
	case operationDivide:
		return a / b
	case operationMultiply:
		return a * b
	}
	panic("Unrecognized operation {" + opSign + "}")
}

// inverseOpSign returns the inverse of the given operation
func inverseOpSign(opSign string) string {
	switch opSign {
	case operationAdd:
		return operationSubract
	case operationSubract:
		return operationAdd
	case operationDivide:
		return operationMultiply
	case operationMultiply:
		return operationDivide
	}
	panic("Unrecognized operation {" + opSign + "}")
}

// Monkey yells a number which can be known from start or constructed from other monkey's numbers.
type Monkey struct {
	Name             string
	NumberYelled     int
	HasNumber        bool
	DependsOn        []string
	DependsOperation string
}

// GetNumberYelled returns monkey's number or calculates it based on dependencies.
func (m *Monkey) GetNumberYelled(allMonkeys map[string]*Monkey) int {
	if m.HasNumber {
		return m.NumberYelled
	}
	monkey1 := allMonkeys[m.DependsOn[0]]
	monkey2 := allMonkeys[m.DependsOn[1]]
	m.NumberYelled = doOperation(monkey1.GetNumberYelled(allMonkeys), monkey2.GetNumberYelled(allMonkeys), m.DependsOperation)
	m.HasNumber = true
	return m.NumberYelled
}

// RewriteDependency makes current monkey become dependant on one of dependants based on equation reordering.
func (m *Monkey) RewriteDependency(extractedDependant string, allMonkeys map[string]*Monkey) {
	dependant := allMonkeys[extractedDependant]
	if len(dependant.DependsOn) == 0 {
		dependant.DependsOn = make([]string, 2)
		dependant.HasNumber = false
	}

	if m.DependsOn[0] == extractedDependant {
		// m = X + b => X = m-b
		// m = X - b => X = m+b
		// m = X * b => X = m/b
		// m = X / b => X = m*b
		dependant.DependsOn[0] = m.Name
		dependant.DependsOn[1] = m.DependsOn[1]
		dependant.DependsOperation = inverseOpSign(m.DependsOperation)
	} else if m.DependsOn[1] == extractedDependant {
		// m = b + X => X = m-b
		// m = b - X => X = b-m
		// m = b * X => X = m/b
		// m = b / X => X = b/m
		dependant.DependsOperation = m.DependsOperation
		dependant.DependsOn[0] = m.DependsOn[0]
		dependant.DependsOn[1] = m.Name
		switch m.DependsOperation {
		case operationAdd:
			dependant.DependsOperation = operationSubract
			dependant.DependsOn[0] = m.Name
			dependant.DependsOn[1] = m.DependsOn[0]
		case operationMultiply:
			dependant.DependsOperation = operationDivide
			dependant.DependsOn[0] = m.Name
			dependant.DependsOn[1] = m.DependsOn[0]
		}
	}
}

// NewMonkey creates a new Monkey struct from challenge input line.
func NewMonkey(monkeyDesc string) *Monkey {
	monkey := &Monkey{}
	parts := strings.Split(monkeyDesc, ": ")
	monkey.Name = parts[0]

	if n, err := strconv.Atoi(parts[1]); err == nil {
		monkey.NumberYelled = n
		monkey.HasNumber = true
	} else {
		dependsOn := strings.Fields(parts[1])

		if len(dependsOn) != 3 {
			panic(fmt.Sprintf("Expected operation between 2 monkeys if monkey number isn't already given (given %v)", parts[1]))
		}
		monkey.DependsOn = []string{dependsOn[0], dependsOn[2]}
		monkey.DependsOperation = dependsOn[1]
	}
	return monkey
}

// ParseMonkeys returns a map of monkeyName pointing to Monkey pointer
func ParseMonkeys(challengeInput string) map[string]*Monkey {
	monkeys := make(map[string]*Monkey)
	monkeyDescs := strings.Split(challengeInput, "\n")
	for _, monkeyDesc := range monkeyDescs {
		m := NewMonkey(monkeyDesc)
		monkeys[m.Name] = m
	}
	return monkeys
}

// findDependentMonkey returns the monkey who is dependent on monkeyName.
func findDependentMonkey(monkeyName string, allMonkeys map[string]*Monkey) *Monkey {
	for _, m := range allMonkeys {
		if m.HasNumber {
			continue
		}
		if m.DependsOn[0] == monkeyName || m.DependsOn[1] == monkeyName {
			return m
		}
	}
	return nil
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	monkeys := ParseMonkeys(input)
	if challengePart == 1 {
		return monkeys[wantedMonkeyName].GetNumberYelled(monkeys)
	}
	if challengePart == 2 {
		// The process on example.in
		// ptdq = humn - dvpt -> humn = ptdq + dvpt
		// lgvd = ljgn * ptdq -> ptdq = lgvd / ljgn
		// cczh = sllz + lgvd -> lgvd = cczh - sllz
		// pppw = cczh / lfgf -> cczh = pppw * lfgf
		// root = pppw + sjmn -> root = pppw = sjmn + 0 (zer0)
		zer0 := &Monkey{Name: "zer0", HasNumber: true, NumberYelled: 0}
		monkeys[zer0.Name] = zer0
		rewrittenName := myMonkeyName
		dependant := findDependentMonkey(rewrittenName, monkeys)
		for dependant.Name != wantedMonkeyName {
			// Since "rewrittenName" will now also depend on dependant, we look for the original dependant beforehand
			nextDependant := findDependentMonkey(dependant.Name, monkeys)
			// fmt.Printf("%s = %s %s %s \n", dependant.Name, dependant.DependsOn[0], dependant.DependsOperation, dependant.DependsOn[1])
			dependant.RewriteDependency(rewrittenName, monkeys)
			// fmt.Printf("  => %s = %s %s %s \n", rewrittenName, monkeys[rewrittenName].DependsOn[0], monkeys[rewrittenName].DependsOperation, monkeys[rewrittenName].DependsOn[1])
			rewrittenName = dependant.Name
			dependant = nextDependant
		}
		// Rewrite one of the dependants of root, so that
		// root = rewrittenName = dependantB => rewrittenName = dependantB + zer0
		rootMonkey := monkeys[wantedMonkeyName]
		lastRewrite := monkeys[rewrittenName]
		lastRewrite.DependsOn[0] = zer0.Name
		lastRewrite.DependsOperation = operationAdd
		if rootMonkey.DependsOn[0] == rewrittenName {
			lastRewrite.DependsOn[1] = rootMonkey.DependsOn[1]
		} else if rootMonkey.DependsOn[1] == rewrittenName {
			lastRewrite.DependsOn[1] = rootMonkey.DependsOn[0]
		}
		// fmt.Printf("  => %s = %s %s %s \n", lastRewrite.Name, lastRewrite.DependsOn[0], lastRewrite.DependsOperation, lastRewrite.DependsOn[1])
		return monkeys[myMonkeyName].GetNumberYelled(monkeys)
	}
	return -1
}

func main() {
	fmt.Println(runChallenge(2))
}
