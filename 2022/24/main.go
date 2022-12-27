package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
)

var (
	//go:embed example.in
	input string

	relativeUp    = Point{0, -1}
	relativeRight = Point{1, 0}
	relativeDown  = Point{0, 1}
	relativeLeft  = Point{-1, 0}

	// blizzardDirections maps text direction to relative coordinates.
	blizzardDirections = map[string]Point{
		blizzardUp:    relativeUp,
		blizzardRight: relativeRight,
		blizzardDown:  relativeDown,
		blizzardLeft:  relativeLeft,
	}

	// exploreDirections are neighbouring locations expedition can move to.
	exploreDirections = []Point{
		relativeDown,
		relativeRight,
		relativeUp,
		relativeLeft,
		{0, 0},
	}
)

// Markings of challenge input.
const (
	valleyWall       = "#"
	safeGround       = "."
	blizzardLeft     = "<"
	blizzardRight    = ">"
	blizzardUp       = "^"
	blizzardDown     = "v"
	expeditionMarker = "E"
)

// Point contains 2D coordinates.
type Point struct {
	X, Y int
}

// add performs sum on each coordinates.
func (p Point) add(d Point) Point {
	return Point{p.X + d.X, p.Y + d.Y}
}

// Path represents a pathing
type Path struct {
	Previous *Path
	Steps    int
	Location Point
	Priority int
}

// Maze represents a board with moving blizzards we wish to navigate.
type Maze struct {
	Map             [][]string
	blizzards       []*Blizzard
	futureBlizzards map[Point]int

	ShortestPath *Path
}

// moveBlizzards makes all blizzards move to their next location on map.
func (m *Maze) moveBlizards() {
	m.futureBlizzards = make(map[Point]int)
	blizzardCount := make(map[Point]int)
	for _, blizzard := range m.blizzards {
		if _, hasBlizzard := blizzardCount[blizzard.Location]; !hasBlizzard {
			m.Map[blizzard.Location.Y][blizzard.Location.X] = safeGround
		}
		blizzard.Location = blizzard.NextLocation
		if blizzards, hasBlizzard := blizzardCount[blizzard.Location]; hasBlizzard {
			m.Map[blizzard.NextLocation.Y][blizzard.NextLocation.X] = fmt.Sprint(blizzards + 1)
		} else {
			m.Map[blizzard.Location.Y][blizzard.Location.X] = blizzard.Marker
		}
		blizzard.setNextAdvance(len(m.Map[0])-1, len(m.Map)-1)
		m.futureBlizzards[blizzard.NextLocation] += 1
		blizzardCount[blizzard.Location] += 1
	}
}

// isWall returns true if given location is a wall or out of bounds.
func (m *Maze) isWall(location Point) bool {
	return location.X < 0 || location.Y < 0 || location.X >= len(m.Map[0]) || location.Y >= len(m.Map) || m.Map[location.Y][location.X] == valleyWall
}

// MoveExpeditionTo finds shortest path among blizzards from start to goal.
func (m *Maze) MoveExpeditionTo(start, goal Point) *Path {
	explore := make(map[Point]*Path)
	explore[start] = &Path{Previous: nil, Steps: 0, Location: start}
	time := 0
	// Quick & dirty safeguard if a path doesn't exist
	for time < 10000 {
		time++
		next := make(map[Point]*Path)
		for _, path := range explore {
			if path.Location == goal {
				return path
			}
			for _, dir := range exploreDirections {
				// Check 5 directions - wait, up, right, down, left
				expeditionMove := path.Location.add(dir)
				_, blizzardLocation := m.futureBlizzards[expeditionMove]
				if blizzardLocation || m.isWall(expeditionMove) {
					// fmt.Printf("  %v is wall\n", expeditionMove)
					continue
				}
				// fmt.Println("Adding ", expeditionMove)
				next[expeditionMove] = &Path{Previous: path, Steps: time, Location: expeditionMove}

			}
		}
		explore = next
		m.moveBlizards()
	}
	return nil
}

// String returns ANSI colored text view of the map.
func (m *Maze) String() string {
	out := ""
	for y := range m.Map {
		for x := range m.Map[y] {
			switch m.Map[y][x] {
			case expeditionMarker:
				out += fmt.Sprintf("\u001b[31m%s\u001b[0m", expeditionMarker)
			case valleyWall:
				out += fmt.Sprintf("\u001b[33m%s\u001b[0m", valleyWall)
			case safeGround, "S", "G":
				out += m.Map[y][x]
			default:
				// Blizzards
				out += fmt.Sprintf("\u001b[34m%s\u001b[0m", m.Map[y][x])
			}
		}
		out += "\n"
	}
	return out
}

// Blizzard describes a moving snow blizzard we wish to avoid.
type Blizzard struct {
	Location     Point
	NextLocation Point
	Direction    Point
	Marker       string
}

// nextAdvance calculates the next spot blizzard will move to based on map boundaries.
func (b *Blizzard) setNextAdvance(wallX, wallY int) {
	b.NextLocation = b.Location.add(b.Direction)
	if b.NextLocation.Y == wallY {
		b.NextLocation.Y = 1
	} else if b.NextLocation.Y == 0 {
		b.NextLocation.Y = wallY - 1
	} else if b.NextLocation.X == wallX {
		b.NextLocation.X = 1
	} else if b.NextLocation.X == 0 {
		b.NextLocation.X = wallX - 1
	}
}

// parseMaze creates a maze structure from challenge input.
func parseMaze(mazeDesc string) (maze *Maze, start, goal Point) {
	maze = &Maze{}
	lines := strings.Split(mazeDesc, "\n")
	maze.Map = make([][]string, len(lines))
	maze.blizzards = make([]*Blizzard, 0)
	maze.futureBlizzards = make(map[Point]int)

	for y := range lines {
		spots := strings.Split(lines[y], "")
		maze.Map[y] = make([]string, len(spots))
		for x := range spots {
			if y == 0 && spots[x] == safeGround {
				start = Point{x, y}
			} else if y == len(lines)-1 && spots[x] == safeGround {
				goal = Point{x, y}
			}
			maze.Map[y][x] = spots[x]
			if dir, ok := blizzardDirections[spots[x]]; ok {
				blizzard := &Blizzard{Location: Point{x, y}, Direction: dir, Marker: spots[x]}
				blizzard.setNextAdvance(len(spots)-1, len(lines)-1)
				maze.blizzards = append(maze.blizzards, blizzard)
				maze.futureBlizzards[blizzard.NextLocation]++
			}
		}
	}
	return
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	maze, start, goal := parseMaze(input)
	initialPathing := maze.MoveExpeditionTo(start, goal)

	if challengePart == 1 {
		// Move start -> goal
		return initialPathing.Steps
	} else if challengePart == 2 {
		// Move start -> goal -> start -> goal and return steps
		backtrack := maze.MoveExpeditionTo(goal, start)
		thereAgain := maze.MoveExpeditionTo(start, goal)
		return initialPathing.Steps + backtrack.Steps + thereAgain.Steps
	} else if challengePart == 3 {
		// Visualize movement of challenge 1 path
		nodes := make([]Point, 0)
		for initialPathing != nil {
			nodes = append([]Point{initialPathing.Location}, nodes...)
			initialPathing = initialPathing.Previous
		}
		maze, start, goal := parseMaze(input)
		maze.Map[start.Y][start.X] = "S"
		maze.Map[goal.Y][goal.X] = "G"
		for _, node := range nodes {
			previous := maze.Map[node.Y][node.X]
			maze.Map[node.Y][node.X] = expeditionMarker
			fmt.Printf("\033[2J")
			fmt.Println(maze)
			maze.Map[node.Y][node.X] = previous
			maze.moveBlizards()
			time.Sleep(250 * time.Millisecond)
		}
	}
	return -1
}

func main() {
	runChallenge(3)
}
