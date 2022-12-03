package main

import (
	_ "embed"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/exp/maps"
)

//go:embed challenge.in
var input string

const (
	upperCaseShift = 38
	lowerCaseShift = 96
	elvesInGroup   = 3
)

func findSmallestRucksack(groups map[int]map[rune]int) int {
	smallest := 0
	for key, _ := range groups {
		if len(groups[key]) < len(groups[smallest]) {
			smallest = key
		}
	}
	return smallest
}

func itemToPriority(item rune) int {
	shift := lowerCaseShift
	if unicode.IsUpper(item) {
		shift = upperCaseShift
	}
	return int(item) - shift
}

func challenge1() int {
	compartment1 := make(map[rune]int, 0)
	compartment2 := make(map[rune]int, 0)
	var commonItems []rune
	for _, rucksack := range strings.Split(input, "\n") {
		half := len(rucksack) / 2

		for i, item := range rucksack {
			if i < half {
				compartment1[item]++
			} else {
				compartment2[item]++
			}
		}
		fmt.Println(compartment1)
		for item, _ := range compartment1 {
			if _, ok := compartment2[item]; ok {
				commonItems = append(commonItems, item)
			}
		}
		maps.Clear(compartment1)
		maps.Clear(compartment2)
	}
	priorityScore := 0
	for _, item := range commonItems {
		priorityScore += itemToPriority(item)
	}
	return priorityScore
}

func challenge2() int {
	elfGroups := make(map[int]map[rune]int, 0)
	var authBadges []rune

	for i, rucksack := range strings.Split(input, "\n") {
		groupNum := i % elvesInGroup
		if _, ok := elfGroups[i]; !ok {
			elfGroups[groupNum] = make(map[rune]int, 0)
		}
		for _, item := range rucksack {
			elfGroups[groupNum][item]++
		}
		smallestRucksack := findSmallestRucksack(elfGroups)
		for item, _ := range elfGroups[smallestRucksack] {
			allContain := true
			for g := 0; g < elvesInGroup; g++ {
				if g == smallestRucksack {
					continue
				}
				if _, ok := elfGroups[g][item]; !ok {
					allContain = false
					break
				}
			}
			if allContain {
				authBadges = append(authBadges, item)
			}
		}
		if (groupNum + 1) == elvesInGroup {
			maps.Clear(elfGroups)
		}
	}

	priorityScore := 0
	for _, item := range authBadges {
		priorityScore += itemToPriority(item)
	}

	return priorityScore
}

func runChallenge(challengePart int) int {
	priorityScore := 0
	if challengePart == 1 {
		priorityScore = challenge1()
	} else {
		priorityScore = challenge2()
	}
	return priorityScore
}

func main() {
	priorityScore := runChallenge(2)
	fmt.Printf("Priority is %d \n", priorityScore)
}
