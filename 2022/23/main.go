package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
)

var (
	//go:embed example.in
	input string

	north     = Point{0, -1}
	northEast = Point{1, -1}
	east      = Point{1, 0}
	southEast = Point{1, 1}
	south     = Point{0, 1}
	southWest = Point{-1, 1}
	west      = Point{-1, 0}
	northWest = Point{-1, -1}

	directions8 = []Point{
		northEast, north, northWest, southEast, south, southWest, southWest, west, northWest, northEast, east, southEast,
	}
)

const (
	// elfMarker represents space occupied by an elf.
	elfMarker = "#"
	// groundMarker represents empty space.
	groundMarker = "."

	offsetX = 150
	offsetY = 150
)

type Point struct {
	X, Y int
}

func (p Point) add(d Point) Point {
	return Point{p.X + d.X, p.Y + d.Y}
}

type Elf struct {
	// Location holds grove coordinates of elf's current location.
	Location Point
	// ProposedMove holds a relative direction for wanted move.
	ProposedMove Point
}

type Grove struct {
	Ground [][]string
	Elves  []*Elf
	// DesiredDirection points to which direction elves want to travel in current round.
	DesiredDirection int
}

func (g *Grove) CountEmptySpots() int {
	smallestX := math.MaxInt
	smallestY := math.MaxInt
	biggestX := 0
	biggestY := 0

	for _, elf := range g.Elves {
		if elf.Location.X < smallestX {
			smallestX = elf.Location.X
		}
		if elf.Location.X > biggestX {
			biggestX = elf.Location.X
		}
		if elf.Location.Y > biggestY {
			biggestY = elf.Location.Y
		}
		if elf.Location.Y < smallestY {
			smallestY = elf.Location.Y
		}
	}

	elfSeparationWidth := biggestX - smallestX + 1
	elfSeparationDepth := biggestY - smallestY + 1

	return elfSeparationWidth*elfSeparationDepth - len(g.Elves)

}

func (g *Grove) SpaceAt(p Point) string {
	return g.Ground[p.Y][p.X]
}

func (g *Grove) acceptProposedMove(elf *Elf) {
	g.Ground[elf.Location.Y][elf.Location.X] = groundMarker
	// fmt.Printf("Elf move from %v to %v \n", elf.Location, elf.ProposedMove)
	elf.Location = Point{elf.ProposedMove.X, elf.ProposedMove.Y}
	g.Ground[elf.Location.Y][elf.Location.X] = elfMarker

}

// MoveElves spaces out elves. Returns true if atleast one elf has moved.
func (g *Grove) MoveElves() bool {
	movesOnto := make(map[Point]int)
	for _, elf := range g.Elves {
		elf.ProposedMove = Point{}
		foundElf := false
		hasMove := false
		for dir := 1; dir < 11; dir += 3 {
			d1 := (g.DesiredDirection + dir - 1) % 12
			d2 := (g.DesiredDirection + dir) % 12
			d3 := (g.DesiredDirection + dir + 1) % 12
			stepPos := elf.Location.add(directions8[d2])
			if g.SpaceAt(stepPos) == groundMarker && g.SpaceAt(elf.Location.add(directions8[d1])) == groundMarker && g.SpaceAt(elf.Location.add(directions8[d3])) == groundMarker {
				if !hasMove {
					elf.ProposedMove = stepPos
					// fmt.Printf("  > wants spot %v\n", stepPos)
					movesOnto[stepPos] += 1
					hasMove = true
				}
			} else {
				foundElf = true
				if hasMove {
					break
				}
			}
		}
		if !foundElf {
			movesOnto[elf.ProposedMove] -= 1
			elf.ProposedMove = Point{}
		}
	}
	movement := false
	for _, elf := range g.Elves {
		if elf.ProposedMove.X == 0 && elf.ProposedMove.Y == 0 {
			continue
		}
		if spotWantedLevel := movesOnto[elf.ProposedMove]; spotWantedLevel == 1 {
			g.acceptProposedMove(elf)
			movement = true
		}
	}
	g.DesiredDirection = (g.DesiredDirection + 3) % 12

	return movement
}

// String outputs the string representation of grove.
func (g *Grove) String() string {
	out := ""
	for y := range g.Ground {
		for x := range g.Ground[y] {
			out += g.Ground[y][x]
			if g.Ground[y][x] == "" {
				out += groundMarker
			}
		}
		out += "\n"
	}
	return out
}

// PraseGrove structure the challenge input data.
func ParseGrove(groveDesc string) *Grove {
	grove := &Grove{}
	lines := strings.Split(groveDesc, "\n")
	grove.Ground = make([][]string, offsetY*2+1)
	grove.Elves = make([]*Elf, 0)

	for y := range grove.Ground {
		grove.Ground[y] = make([]string, offsetX*2+1)
		for x := range grove.Ground[y] {
			grove.Ground[y][x] = groundMarker
		}
	}

	for y := range lines {
		spaces := strings.Split(lines[y], "")
		for x := range spaces {
			grove.Ground[offsetY+y][offsetX+x] = spaces[x]
			if spaces[x] == elfMarker {
				grove.Elves = append(grove.Elves, &Elf{Location: Point{offsetX + x, offsetY + y}})
			}
		}
	}
	return grove
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	grove := ParseGrove(input)
	movement := true
	round := 0
	for movement {
		movement = grove.MoveElves()
		round++
		if challengePart == 1 && round == 10 {
			return grove.CountEmptySpots()
		}
	}
	if challengePart == 2 {
		return round
	}
	return -1
}

func main() {
	fmt.Println(runChallenge(2))
}
