package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	crateCaptureGroup = `\[([A-Z])\]`
	crateMoveNumbers  = `[0-9]+`
)

//go:embed example.in
var input string

type Stack []string

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

func (s *Stack) PushN(strs []string) {
	*s = append(*s, strs...)
}

func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	}
	last := len(*s) - 1
	lastElement := (*s)[last]
	*s = (*s)[:last]
	return lastElement, true
}

func (s *Stack) PopN(n int) []string {
	if n > len(*s) {
		n = len(*s)
	}
	newLen := len(*s) - n
	popped := (*s)[newLen:]
	*s = (*s)[:newLen]
	return popped
}

func (s *Stack) PeekLast() string {
	if (*s).IsEmpty() {
		return ""
	}
	return (*s)[len(*s)-1]
}

type CargoShip struct {
	crateStacks []Stack
}

type MoveInstruction struct {
	nCrates     int
	sourceStack int
	destStack   int
}

func RecklessParseInt(s string) int {
	n, _ := strconv.ParseInt(s, 10, 0)
	return int(n)
}

func parseInput(crateDesc string) (*CargoShip, []MoveInstruction) {
	lines := strings.Split(crateDesc, "\n")
	crateSectionEnd := 0
	inCratesSection := true
	numStacks := 0
	moves := make([]MoveInstruction, 0)
	moveSplitter := regexp.MustCompile(crateMoveNumbers)
	for i, line := range lines {
		//fmt.Println(line)
		if line == "" {
			crateSectionEnd = i - 2 // Naming is trivial so ignore the number line
			numStacks = len(strings.Fields(lines[i-1]))
			inCratesSection = false
			continue
		}
		if !inCratesSection {
			// will return [N S D] (N - number of crates, S - source, D - destionation)
			// Source and destination crates start counting from 1
			inst := moveSplitter.FindAllString(line, -1)
			moves = append(moves,
				MoveInstruction{
					nCrates:     RecklessParseInt(inst[0]),
					sourceStack: RecklessParseInt(inst[1]) - 1,
					destStack:   RecklessParseInt(inst[2]) - 1,
				},
			)
		}
	}
	// Parse the crates in reverse to properly populate stack
	crateMatcher := regexp.MustCompile(crateCaptureGroup)
	cargoShip := &CargoShip{}
	cargoShip.crateStacks = make([]Stack, numStacks)
	for i := crateSectionEnd; i >= 0; i-- {
		// Returns [ [ groupStart groupEnd insideStart insideEnd ] ... ]
		for _, m := range crateMatcher.FindAllStringSubmatchIndex(lines[i], -1) {
			crateName := lines[i][m[2]:m[3]]
			crateStack := m[2] / 4
			//fmt.Printf("Crate named [%s] belongs to stack %d \n", crateName, crateStack)
			cargoShip.crateStacks[crateStack].Push(crateName)
		}
	}

	return cargoShip, moves
}

func runChallenge(challengePart int) string {
	ship, moves := parseInput(input)
	result := ""
	for _, move := range moves {
		if challengePart == 1 {
			for i := 0; i < move.nCrates; i++ {
				crate, ok := ship.crateStacks[move.sourceStack].Pop()
				if ok {
					ship.crateStacks[move.destStack].Push(crate)
				}
			}
		} else if challengePart == 2 {
			crates := ship.crateStacks[move.sourceStack].PopN(move.nCrates)
			ship.crateStacks[move.destStack].PushN(crates)
		}
	}
	for _, stack := range ship.crateStacks {
		result += stack.PeekLast()
	}

	return result
}

func main() {
	result := runChallenge(2)
	fmt.Println(result)
}
