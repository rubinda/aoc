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

	// SandSource is the point where sand starts flowing in. NEEDS OFFSET!
	SandSource = Point{500, 0}
	// MaterialSymbols contains string representations of materials that can occupy a space.
	MaterialSymbols = map[int]string{
		Void:       ".",
		SandSupply: "+",
		Rock:       "#",
		Sand:       "*",
	}
)

// Relative directions based on (0, 0) being the upper left corner.
var (
	relativeDownRight  = Point{x: 1, y: 1}
	relativeDown       = Point{x: 0, y: 1}
	relativeDownLeft   = Point{x: -1, y: 1}
	sandFlowDirections = []Point{relativeDown, relativeDownLeft, relativeDownRight}
)

// Represents possible type of Material occupying a point in space.
const (
	Void = iota
	SandSupply
	Rock
	Sand
)

// Point represents coordinates.
type Point struct {
	x int
	y int
}

// add sums coordinates of points.
func (p Point) add(d Point) Point {
	return Point{p.x + d.x, p.y + d.y}
}

// check panics on non nil error.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Sandbox holds space which can be filled with material.
type Sandbox struct {
	space      [][]int
	hasBottom  bool
	bottom     int
	width      int
	offsetX    int
	sandSource Point
}

// isOccupied returns true if given point is already occupied with non-Void material.
func (s *Sandbox) isOccupied(p Point) bool {
	return s.space[p.y][p.x] != Void
}

// moveGrainOfSand tries to move a sandcorn one space down
func (s *Sandbox) moveGrainOfSand(currentPos Point) (Point, bool) {
	if currentPos.y+1 == len(s.space) {
		return currentPos, false
	}

	for _, dir := range sandFlowDirections {
		newPos := currentPos.add(dir)
		if !s.isOccupied(newPos) {
			return newPos, true
		}
		if newPos.x-1 < 0 || newPos.y-1 < 0 || newPos.x+1 == len(s.space[0]) || newPos.y+1 == len(s.space) {
			return currentPos, false
		}
	}
	return currentPos, false
}

// SpawnGrainOfSand creates a new corn of sand and tries to move it as far as possible
func (s *Sandbox) SpawnGrainOfSand() bool {
	sandPos, canMove := s.moveGrainOfSand(s.sandSource)
	if !canMove {
		fmt.Println("STOP! Can't spawn more sand, source is blocked!")
		return false
	}

	for canMove {
		sandPos, canMove = s.moveGrainOfSand(sandPos)
	}
	if !s.hasBottom && sandPos.y == s.bottom {
		// Sand reached the bottom of the void!
		fmt.Println("STOP! Sand started flowing onto bottom")
		return false
	}
	if !canMove {
		// Sand has settled
		s.space[sandPos.y][sandPos.x] = Sand
		return true
	}
	return false
}

// DrawWall occupies space from given extremes in a straight line. Panics if non straight line given.
func (s *Sandbox) DrawWall(wallStart, wallEnd Point) {
	if wallStart.x == wallEnd.x {
		// vertical wall on y
		if wallStart.y < wallEnd.y {
			for y := wallStart.y; y <= wallEnd.y; y++ {
				s.space[y][wallStart.x] = Rock
			}
		} else {
			for wallStart.y >= wallEnd.y {
				s.space[wallStart.y][wallStart.x] = Rock
				wallStart.y--
			}
		}
	} else if wallStart.y == wallEnd.y {
		// horizontal wall on x
		if wallStart.x < wallEnd.x {
			for x := wallStart.x; x <= wallEnd.x; x++ {
				s.space[wallStart.y][x] = Rock
			}
		} else {
			for wallStart.x >= wallEnd.x {
				s.space[wallStart.y][wallStart.x] = Rock
				wallStart.x--
			}
		}
	} else {
		panic(fmt.Sprintf("I'm not supposed to draw diagonal walls (got %v -> %v)", wallStart, wallEnd))
	}
}

// Output formats the sandbox into ASCII art.
func (s *Sandbox) Output() string {
	output := ""
	for y := range s.space {
		for x := range s.space[y] {
			output += MaterialSymbols[s.space[y][x]]
		}
		output += "\n"
	}
	return output
}

// praseWallEdge returns coordinates from comma delimited value (e.g. "498,6" -> Point{498, 6, Rock}).
func parseWallEdge(wallEdgeDesc string) Point {
	coords := strings.Split(wallEdgeDesc, ",")
	x, err := strconv.Atoi(coords[0])
	check(err)
	y, err := strconv.Atoi(coords[1])
	check(err)

	return Point{x, y}
}

// InitSandbox creates a new sandbox from wall descriptions. Tries to create the optimal sandbox size.
func InitSandbox(sandBoxDesc string, hasBottom bool) *Sandbox {
	sandbox := &Sandbox{}
	sandbox.hasBottom = hasBottom

	// Parse wall instructions first so we can make the optimal sized sandbox
	wallPoints := make([]Point, 0)
	walls := strings.Split(sandBoxDesc, "\n")
	minX := -1
	for _, wall := range walls {
		edges := strings.Split(wall, " -> ")
		wallStart := parseWallEdge(edges[0])
		for i := 1; i < len(edges); i++ {
			// Find bottom -> helps calculate the optimal sandbox size
			if wallStart.y > sandbox.bottom {
				sandbox.bottom = wallStart.y
			}
			if wallStart.x > sandbox.width {
				sandbox.width = wallStart.x
			}
			if minX == -1 || minX > wallStart.x {
				minX = wallStart.x
			}
			wallEnd := parseWallEdge(edges[i])
			wallPoints = append(wallPoints, wallStart, wallEnd)
			wallStart = wallEnd
		}
	}
	// As per instrcutions, bottom is 2 spaces lower than lowest wall
	sandbox.bottom += 2
	// Minimum width to pad the leftmost and rightmost wall with 1 Void space
	sandbox.width -= (minX - 3)

	if sandbox.hasBottom && sandbox.width < (2*sandbox.bottom+1) {
		// Since sand flows into a triangle, optimal width = 2N+1 (see how to draw triangles with ASCII)
		sandbox.width = 2*sandbox.bottom + 1
		sandbox.offsetX = 500 - (sandbox.width / 2)
	} else {
		// Offset is such that leftmost wall is padded by 1
		sandbox.offsetX = minX - 1
	}
	// Create sandbox with size that fits all sand to up to sink (or walls)
	sandbox.space = make([][]int, sandbox.bottom+1)
	for i := range sandbox.space {
		sandbox.space[i] = make([]int, sandbox.width)
	}
	// Calculate the offset for X coordinates (so source is ~middle of sandbox) and mark the sand source
	offsetPoint := Point{-sandbox.offsetX, 0}
	sandbox.sandSource = SandSource.add(offsetPoint)
	sandbox.space[SandSource.y][SandSource.x-sandbox.offsetX] = SandSupply
	// Draw walls into sandbox

	for i := 1; i < len(wallPoints); i += 2 {
		sandbox.DrawWall(wallPoints[i-1].add(offsetPoint), wallPoints[i].add(offsetPoint))
	}
	// Draw bottom
	if sandbox.hasBottom {
		sandbox.DrawWall(Point{0, sandbox.bottom}, Point{len(sandbox.space[0]) - 1, sandbox.bottom})
	}
	return sandbox
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	hasBottom := false
	if (challengePart) == 2 {
		hasBottom = true
	}
	sandbox := InitSandbox(input, hasBottom)
	canSpawnMore := sandbox.SpawnGrainOfSand()
	cornsSpawned := 0
	for canSpawnMore {
		cornsSpawned++
		canSpawnMore = sandbox.SpawnGrainOfSand()
	}
	fmt.Println(sandbox.Output())
	if challengePart == 2 {
		// I like the idea that sand source stays where it originally was, but challenge wants it to change into sand.
		// So add 1 to get proper result.
		return cornsSpawned + 1
	}
	return cornsSpawned
}

func main() {
	fmt.Println("Grains of sand produced: ", runChallenge(2))
}
