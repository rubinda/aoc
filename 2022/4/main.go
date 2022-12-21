package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed example.in
var input string

type ElfCleaner struct {
	sectionMin int
	sectionMax int
}

func (a ElfCleaner) overlaps(b ElfCleaner) bool {
	return !(a.sectionMax < b.sectionMin || b.sectionMax < a.sectionMin)
}

func (a ElfCleaner) fullyContains(b ElfCleaner) bool {
	return a.sectionMin <= b.sectionMin && a.sectionMax >= b.sectionMax
}

func parseAssignments(assignmentDesc string) []ElfCleaner {
	sections := strings.Split(assignmentDesc, ",")
	cleaners := make([]ElfCleaner, len(sections))
	for i, section := range sections {
		bounds := strings.Split(section, "-")
		cleaners[i].sectionMin, _ = strconv.Atoi(bounds[0])
		cleaners[i].sectionMax, _ = strconv.Atoi(bounds[1])
	}
	return cleaners
}

func isFullyContained(cleaners []ElfCleaner, n int) bool {
	for i, elf := range cleaners {
		if i == n {
			continue
		}
		if cleaners[n].fullyContains(elf) {
			return true
		}
	}
	return false
}

func hasOverlap(cleaners []ElfCleaner, n int) bool {
	for i, elf := range cleaners {
		if i == n {
			continue
		}

		if cleaners[n].overlaps(elf) {
			return true
		}
	}
	return false
}

func runChallenge(challengePart int) int {
	contained := 0
	for _, assignment := range strings.Split(input, "\n") {
		cleaners := parseAssignments(assignment)

		for i := range cleaners {
			if challengePart == 1 && isFullyContained(cleaners, i) {
				contained++
				// If the cleaners have the same ranges it counts as 1 contain, not 2 (15-20,15-20 => 1 overlap not 2)
				break
			} else if challengePart == 2 && hasOverlap(cleaners, i) {
				contained++
				break
			}
		}
	}
	return contained
}

func main() {
	answer := runChallenge(2)
	fmt.Println(answer)
}
