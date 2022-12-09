package main

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
)

//go:embed challenge.in
var input string

// sign returns the integer sign of a number.
func sign(a int) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return 1
	}
	return 0
}

// position contains coordinates for an imaginary grid.
type position struct {
	x int
	y int
}

// add joins the two positions coordinatewise.
func (p position) add(d position) position {
	p.x += d.x
	p.y += d.y
	return p
}

// distanceTo returns the Euclidean distance between 2 points.
func (p position) distanceTo(b position) float64 {
	y2 := math.Pow(float64(b.y-p.y), 2)
	x2 := math.Pow(float64(b.x-p.x), 2)
	return math.Sqrt(y2 + x2)
}

// Based on bottom left as (0,0).
var (
	up    = position{0, 1}
	right = position{1, 0}
	down  = position{0, -1}
	left  = position{-1, 0}
)

// directionMap contains translations from input to relative coordinates.
var directionMap = map[string]position{
	"U": up,
	"D": down,
	"L": left,
	"R": right,
}

// moveInstruction represents a line of the challenge input.
type moveInstruction struct {
	direction position
	steps     int
}

// rope consists of knots. First one moves according to the instruction, the others follow.
// Almost like a snake with a head.
type rope struct {
	knots []position
	size  int
	// visitedPositions contains number of visits for each location __ONLY__ the __TAIL__ has been to.
	visitedPositions map[position]int
}

// newRope returns an new rope with 'size' knots and initial knot positions at (0,0).
func newRope(size int) *rope {
	r := &rope{
		size:  size,
		knots: make([]position, size),
		visitedPositions: map[position]int{
			{0, 0}: 1,
		},
	}
	for i := 0; i < size; i++ {
		r.knots[i] = position{}
	}
	return r
}

// followPrevious moves the current knot closer to the previous.
// Distance between knots should always be >= 0 and < 2. Returns if current know has moved.
func (r *rope) followPrevious(previous, current int) bool {
	if current >= r.size || r.knots[current].distanceTo(r.knots[previous]) < 2 {
		return false
	}
	// Only move 1 step in any given direction, not gradient!
	dx := sign(r.knots[previous].x - r.knots[current].x)
	dy := sign(r.knots[previous].y - r.knots[current].y)
	moveDirection := position{dx, dy}
	r.knots[current] = r.knots[current].add(moveDirection)
	// Only log the position of the TAIL!!
	if current == r.size-1 {
		r.visitedPositions[r.knots[current]]++
	}
	return true
}

// followHead moves all the knots after their predecessor (while distance between adjacent knots is >= 2).
func (r *rope) followHead() {
	hasMoved := true
	i := 0
	for hasMoved {
		hasMoved = r.followPrevious(i, i+1)
		i++
	}
}

// moveHead moves the head in a given direction (4 DOF) for a given amount of steps.
func (r *rope) moveHead(move moveInstruction) {
	steps := move.steps
	for steps > 0 {
		steps--
		r.knots[0] = r.knots[0].add(move.direction)
		// Diagonal distance is ~1.4142 and still ok
		if r.knots[1].distanceTo(r.knots[0]) >= 2 {
			r.followHead()
		}
	}
}

// parseMoves reads the challenge input - move instructions for our snakey rope. See example.in or challenge.in.
func parseMoves(desc string) []moveInstruction {
	lines := strings.Split(desc, "\n")
	moves := make([]moveInstruction, len(lines))
	for i, line := range lines {
		s := strings.Fields(line)
		steps, _ := strconv.Atoi(s[1])
		moves[i] = moveInstruction{direction: directionMap[s[0]], steps: steps}
	}
	return moves
}

// runChallenge returns the desired output for the days challenge.
// May print additional information to stdout.
func runChallenge(challengePart int) int {
	moves := parseMoves(input)
	knots := 2
	if challengePart == 2 {
		knots = 10
	}
	// Head & body all start at 0,0 (any part can overlap even in future steps)
	snakeyRope := newRope(knots)
	for _, move := range moves {
		snakeyRope.moveHead(move)
	}
	return len(snakeyRope.visitedPositions)

}

func main() {
	fmt.Println(runChallenge(2))
}
