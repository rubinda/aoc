package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Elf struct {
	id               int
	caloriesCarrying int
}

//go:embed challenge.in
var input string

func runChallenge(challengePart int) (int, []Elf) {
	calories := strings.Split(input, "\n")
	var elves []Elf
	elfId := 0
	elves = append(elves, Elf{id: elfId, caloriesCarrying: 0})

	for _, v := range calories {
		if v == "" {
			elfId++
			elves = append(elves, Elf{id: elfId, caloriesCarrying: 0})
			continue
		}
		calories, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		elves[elfId].caloriesCarrying += calories
	}
	sort.Slice(elves, func(i, j int) bool {
		return elves[i].caloriesCarrying > elves[j].caloriesCarrying
	})

	if challengePart == 1 {
		return elves[0].caloriesCarrying, elves[:1]
	} else {
		total := 0
		for _, v := range elves[:3] {
			total += v.caloriesCarrying
		}
		return total, elves[:3]
	}
}

func main() {
	calories, elves := runChallenge(1)
	fmt.Println(elves)
	fmt.Printf("They carry %d calories\n ", calories)
}
