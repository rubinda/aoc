package main

import (
	_ "embed"
	"fmt"
	"strings"
)

var (
	//go:embed example.in
	input string
	// pieceMaterial helps visualizing output
	pieceMaterial = []string{"-", "+", "J", "I", "B"}
	// RockPieces are tetris blocks that repeat.
	RockPieces = []RockPiece{
		{
			{MaterialRock, MaterialRock, MaterialRock, MaterialRock},
		},
		{
			{MaterialVoid, MaterialRock, MaterialVoid},
			{MaterialRock, MaterialRock, MaterialRock},
			{MaterialVoid, MaterialRock, MaterialVoid},
		},
		{
			{MaterialVoid, MaterialVoid, MaterialRock},
			{MaterialVoid, MaterialVoid, MaterialRock},
			{MaterialRock, MaterialRock, MaterialRock},
		},
		{
			{MaterialRock},
			{MaterialRock},
			{MaterialRock},
			{MaterialRock},
		},
		{
			{MaterialRock, MaterialRock},
			{MaterialRock, MaterialRock},
		},
	}
)

// Relative directions for tetris blocks to move.
var (
	relativeLeft  = Point{-1, 0}
	relativeRight = Point{1, 0}
	relativeDown  = Point{0, 1}
)

const (
	// chamberWidth represents the playing field width.
	chamberWidth = 7
	// offsetX represents initial rock position from left chamber edge
	offsetX = 2
	// offsetY represents initial rock position from the floor
	offsetY      = 3
	MaterialVoid = "."
	MaterialRock = "#"

	// When visualizing the chamber, incrementally use sections of depth 100
	chamberSectionDepth = 100
	// Section increases keep these lines from previous section top
	chamberSectionIncrease = 50
	// Input instruction for jet of wind to left side.
	moveLeft = "<"
	// Input instruction for jet of wind to right side.
	moveRight = ">"

	// Blocks spawned in challenge 1
	challenge1Runs = 2022
	// Blocks spawned in challenge 2
	challenge2Runs = 1000000000000
)

// RockPiece is a shortcut name for 2d string arrays.
type RockPiece [][]string

// Height represents the Y axis of a rock piece.
func (r RockPiece) Height() int {
	return len(r)
}

// Width represensts the X axis of a rock piece.
func (r RockPiece) Width() int {
	return len(r[0])
}

// Point holds coordinates.
type Point struct {
	x int
	y int
}

// add returns a new point with coordinate-wise sums.
func (p Point) add(d Point) Point {
	return Point{p.x + d.x, p.y + d.y}
}

// Chamber is the playing field for tetris.
type Chamber struct {
	Section            [][]string
	TowerHeight        int
	sectionTowerHeight int
	WindJets           []string
	piecesSpawned      int
	jetPointer         int
}

// NewChamber returns a new chamber with periodic rock pieces and jets of wind.
func NewChamber(windJets string) *Chamber {
	chamber := &Chamber{}
	chamber.Section = make([][]string, chamberSectionDepth)
	for i := range chamber.Section {
		chamber.Section[i] = make([]string, chamberWidth)
		for j := 0; j < chamberWidth; j++ {
			chamber.Section[i][j] = MaterialVoid
		}
	}
	chamber.WindJets = strings.Split(windJets, "")
	return chamber
}

// Deepen increases space in the chamber for new pieces.
func (c *Chamber) Deepen() {
	moreSpace := make([][]string, chamberSectionIncrease)
	for i := range moreSpace {
		moreSpace[i] = make([]string, chamberWidth)
		for j := 0; j < chamberWidth; j++ {
			moreSpace[i][j] = MaterialVoid
		}
	}
	// Modified for part two from:
	// c.State = append(moreSpace, c.State...)
	c.Section = append(moreSpace, c.Section[:chamberSectionIncrease]...)
	c.sectionTowerHeight = c.sectionTowerHeight % chamberSectionIncrease
}

// NextWindJet returns the next direction a wind jet will blow.
func (c *Chamber) NextWindJet() Point {
	jet := c.WindJets[c.jetPointer%len(c.WindJets)]
	c.jetPointer++
	if jet == moveLeft {
		return relativeLeft
	} else if jet == moveRight {
		return relativeRight
	}
	panic(fmt.Sprintf("Given Jet [%s] is unrecognized", jet))
}

// NextRockPiece returns the next RockPiece to be spawned.
func (c *Chamber) NextRockPiece() RockPiece {
	rockPiece := RockPieces[c.piecesSpawned%len(RockPieces)]
	// fmt.Println("Spawning piece ", c.piecesSpawned%len(RockPieces))
	c.piecesSpawned++
	return rockPiece
}

// SpawnPiece places a new piece on top of chamber and moves it down untill the piece is blocked.
func (c *Chamber) SpawnPiece() {
	rockPiece := c.NextRockPiece()
	// Check if we can place piece on starting point
	if len(c.Section) < ((c.sectionTowerHeight) + offsetY + rockPiece.Height()) {
		// fmt.Println("  >deepen")
		c.Deepen()
	}
	spawnY := len(c.Section) - (c.sectionTowerHeight + offsetY + rockPiece.Height())
	corner := Point{offsetX, spawnY}
	// fmt.Printf("Spawning at (%d, %d) \n", offsetX, spawnY)

	corner, hasMovedDown := c.movePiece(corner, rockPiece)
	// fmt.Printf(" >Move: (%d, %d), next: %v \n", corner.x, corner.y, hasMovedDown)
	for hasMovedDown {
		corner, hasMovedDown = c.movePiece(corner, rockPiece)
		// fmt.Printf(" >Move: (%d, %d), next: %v \n", corner.x, corner.y, hasMovedDown)
	}
	c.drawPiece(corner, rockPiece)
	// Recalculate tower height
	c.increaseTowerHeight()
}

// increaseTowerHeight recalculates the total tetris tower height in the chamber.
func (c *Chamber) increaseTowerHeight() {
	for y := len(c.Section) - 1 - c.sectionTowerHeight; y >= 0; y-- {
		noRock := true
		for x := 0; x < chamberWidth; x++ {
			if c.Section[y][x] != MaterialVoid {
				noRock = false
				c.TowerHeight++
				c.sectionTowerHeight++
				break
			}
		}
		if noRock {
			return
		}
	}
}

// movePiece applies a jet of wind and downward movement onto a RockPiece with given upper left corner.
// Returns false if piece turned solid rock and can't be moved anymore.
func (c *Chamber) movePiece(corner Point, piece RockPiece) (Point, bool) {
	jetDir := c.NextWindJet()
	if c.canPlace(corner.add(jetDir), piece) {
		corner = corner.add(jetDir)
	}
	moved := false
	if moved = c.canPlace(corner.add(relativeDown), piece); moved {
		corner = corner.add(relativeDown)
	}

	return corner, moved
}

// canPlace checks if a RockPiece can fit in chamber with given upper left corner.
func (c *Chamber) canPlace(corner Point, piece RockPiece) bool {
	if corner.x < 0 || corner.x+piece.Width() > chamberWidth || corner.y+piece.Height() > len(c.Section) {
		// fmt.Println("  !place overflow")
		return false
	}

	for y := 0; y < piece.Height(); y++ {
		for x := 0; x < piece.Width(); x++ {
			if piece[y][x] != MaterialVoid && c.Section[corner.y+y][corner.x+x] != MaterialVoid {
				// fmt.Println("  !place hit rock")
				return false
			}
		}
	}
	return true
}

// drawPiece draws a rock onto the canvas.
func (c *Chamber) drawPiece(corner Point, piece RockPiece) {
	material := pieceMaterial[(c.piecesSpawned-1)%len(RockPieces)]
	for y := 0; y < len(piece); y++ {
		for x := 0; x < len(piece[0]); x++ {
			if piece[y][x] == MaterialRock {
				c.Section[corner.y+y][corner.x+x] = material
			}
		}
	}
}

// Output draws the (sectioned) chamber into a string.
func (c *Chamber) Output() string {
	out := ""
	for i := 0; i < len(c.Section); i++ {
		out += "|"
		for j := 0; j < len(c.Section[i]); j++ {
			out += c.Section[i][j]
		}
		out += "|\n"
	}
	out += "+-------+\n"
	return out
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	chamber := NewChamber(input)
	if challengePart == 1 {
		for i := 0; i < challenge1Runs; i++ {
			chamber.SpawnPiece()
		}
		return chamber.TowerHeight
	}
	if challengePart == 2 {
		initialRuns := 10000
		heightGains := make([]int, initialRuns)
		previousHeight := 0

		// isPeriodic := false
		for i := 0; i < initialRuns; i++ {
			chamber.SpawnPiece()
			heightGains[i] = chamber.TowerHeight - previousHeight
			previousHeight = chamber.TowerHeight
		}
		// For some reason, heigh gain becomes periodic after first N runs. BUT if we search for the periodic gains in reverse, there is no offset!
		// So we find the periodic length L first and then offset
		periodic := heightGains[len(heightGains)-5:]
		i := initialRuns - 6
		foundPeriod := false
		for i >= 0 && !foundPeriod {
			allMatch := true
			for j := 0; j < len(periodic); j++ {
				if heightGains[i-len(periodic)+j] != periodic[j] {
					allMatch = false
					break
				}
			}
			if allMatch {
				foundPeriod = true
			}
			periodic = append([]int{heightGains[i]}, periodic...)
			i--
		}
		// fmt.Println("Found periodic growth: ", len(periodic))

		// Find offset
		periodOffset := -1
		foundOffset := false
		for !foundOffset {
			periodOffset++
			foundOffset = true
			for i := range periodic {
				if heightGains[periodOffset+i] != periodic[i] {
					foundOffset = false
					break
				}
			}
		}
		// fmt.Println("Found offset for periodic growth: ", periodOffset)
		periodicRuns := (challenge2Runs - periodOffset) / len(periodic)
		lastPeriodRemainder := (challenge2Runs - periodOffset) - (periodicRuns * len(periodic))

		// fmt.Printf("Need to do %d offset runs, then %d to get period (times %d) and finally add %d of last period \n", periodOffset, len(periodic), periodicRuns, lastPeriodRemainder)
		chamber := NewChamber(input)
		offsetHeight := 0
		periodicGrowth := 0
		remainderGrowth := 0
		for i := 0; i < (periodOffset + len(periodic)); i++ {
			chamber.SpawnPiece()
			if i == (periodOffset - 1) {
				offsetHeight = chamber.TowerHeight
			} else if i == (periodOffset-1)+lastPeriodRemainder {
				remainderGrowth = chamber.TowerHeight - offsetHeight
			} else if i == (periodOffset + len(periodic) - 1) {
				periodicGrowth = chamber.TowerHeight - offsetHeight
			}
		}
		return offsetHeight + (periodicGrowth * periodicRuns) + remainderGrowth
	}

	return -1
}

func main() {
	fmt.Println("Total height: ", runChallenge(2))
}
