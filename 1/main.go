package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Elf struct {
	id               int
	caloriesCarrying int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data, err := os.ReadFile("challenge.in")
	check(err)
	calories := strings.Split(string(data), "\n")
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
		check(err)
		elves[elfId].caloriesCarrying += calories
	}
	sort.Slice(elves, func(i, j int) bool {
		return elves[i].caloriesCarrying > elves[j].caloriesCarrying
	})

	fmt.Println("Tope 3 elves carrying food are: ")
	fmt.Println(elves[:3])
	caloricSum := 0
	for _, elf := range elves[:3] {
		caloricSum += elf.caloriesCarrying
	}
	fmt.Printf("Together they carry: %d\n", caloricSum)
}
