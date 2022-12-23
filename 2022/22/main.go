package main

import (
	_ "embed"
	"fmt"
	"strings"
)

var (
	//go:embed example.in
	input string
	// sortedFacing provides a (index) score for each turn direction and helps turn direction determination.
	sortedFacing = []string{faceRight, faceDown, faceLeft, faceUp}
	// facingDirection provides relative directional coordinates
	relativeFacing = map[string]Position{
		faceRight: {1, 0},
		faceDown:  {0, 1},
		faceLeft:  {-1, 0},
		faceUp:    {0, -1},
	}

	cubeMinBoundaries = []Position{
		{0, 0},
		{cubeSize, 0},
		{0, cubeSize},
		{0, cubeSize * 2},
		{cubeSize, 2 * cubeSize},
		{0, 3 * cubeSize},
	}
	cubeMaxBoundaries = []Position{
		{cubeSize - 1, cubeSize - 1},
		{2*cubeSize - 1, cubeSize - 1},
		{cubeSize - 1, 2*cubeSize - 1},
		{cubeSize - 1, 3*cubeSize - 1},
		{2*cubeSize - 1, 3*cubeSize - 1},
		{cubeSize - 1, 4*cubeSize - 1},
	}
)

const (
	// openTile allows movement.
	openTile = "."
	// solidTile blocks movement.
	solidTile = "#"
	// voidTile wraps map (also present in state overflow).
	voidTile = " "

	// turnClockwise denotes in-place turning for 90 degrees.
	turnClockwise = "R"
	// turnCounterClocwise denotes in-place turning for 90 degrees.
	turnCounterClockwise = "L"

	// faceRight equals facing east.
	faceRight = ">"
	// faceDown equals facing south.
	faceDown = "v"
	// faceLeft equals facing west.
	faceLeft = "<"
	// faceUp equals facing north.
	faceUp = "^"
	// flatMap means the map description is a 2D flat surface
	flatMap = "MAP_FLAT"
	// cubeMap means the map description represents a cube surface. 2x2x2 Cube would looke like:
	//     11
	//     11
	// 223344
	// 223344
	//     5566
	//     5566
	cubeMap = "MAP_CUBE"

	cubeSize = 4
)

// Movement represents a move instruction
type Movement struct {
	Steps         int
	TurnDirection string
}

// Position holds 2D coordinates.
type Position struct {
	X int
	Y int
}

// add performs coordinate-wise sum on 2 positions.
func (p Position) add(d Position) Position {
	return Position{p.X + d.X, p.Y + d.Y}
}

// MonkeysDescription is a representation of Monkey's map.
type MonkeysDescription struct {
	MyPosition Position
	MyFacing   string

	Map     [][]string
	MapType string
	// StepperFunc controls moving in whichever direction we are facing
	StepperFunc func(currentPos Position, facing string) (Position, string)
}

// TurnPlayer rotates current facing clock- or counterclockwise.
func (bm *MonkeysDescription) TurnPlayer(turnDirection string) {
	myDir := 0
	for d := range sortedFacing {
		if sortedFacing[d] == bm.MyFacing {
			myDir = d
		}
	}
	switch turnDirection {
	case turnClockwise:
		myDir = (myDir + 1) % 4
	case turnCounterClockwise:
		myDir--
		if myDir < 0 {
			myDir = 3
		}
	}
	bm.MyFacing = sortedFacing[myDir]
}

// tileAt returns the tile type at given position.
func (md *MonkeysDescription) tileAt(p Position) string {
	if p.X < 0 || p.Y < 0 || p.Y >= len(md.Map) || p.X >= len(md.Map[p.Y]) {
		return voidTile
	}
	return md.Map[p.Y][p.X]
}

// canMove returns if we can move from
func (md *MonkeysDescription) canMove() (Position, bool) {
	// Diagonal movements are not possible atm
	newPos, newFacing := md.StepperFunc(md.MyPosition, md.MyFacing)

	// PART 1
	// Depending on which direction we moved in different map 'overflow' can happen
	// switch md.MyFacing {
	// case faceRight:
	// 	// array overflow (because each depth level col of right-most tile wide)
	// 	if newPos.X >= len(md.Map[newPos.Y]) {
	// 		// Wrap onto left side
	// 		newPos.X = 0
	// 		for md.tileAt(newPos) == voidTile {
	// 			newPos.X++
	// 		}
	// 	}
	// case faceDown:
	// 	// Overflow OR step onto voidTile
	// 	if md.tileAt(newPos) == voidTile {
	// 		newPos.Y = 0
	// 		for md.tileAt(newPos) == voidTile {
	// 			newPos.Y += 1
	// 		}
	// 	}
	// case faceLeft:
	// 	// Overflow OR step onto void tile
	// 	if md.tileAt(newPos) == voidTile {
	// 		// Rightmost tile should always be non void!
	// 		newPos.X = len(md.Map[newPos.Y]) - 1
	// 		if md.tileAt(newPos) == voidTile {
	// 			panic(fmt.Sprintf("Overflow stepped onto void tile! %v", newPos))
	// 		}
	// 	}
	// case faceUp:
	// 	// Overflow OR step onto void tile
	// 	if md.tileAt(newPos) == voidTile {
	// 		newPos.Y = len(md.Map) - 1
	// 		for md.tileAt(newPos) == voidTile {
	// 			newPos.Y--
	// 		}
	// 	}
	// }
	// if md.tileAt(newPos) == solidTile {
	// 	return md.MyPosition, false
	// }

	// Draw my steps onto map so it's easier to debug!

	if md.tileAt(newPos) == solidTile {
		// fmt.Printf("  (!) no move\n")
		return md.MyPosition, false
	}
	md.MyFacing = newFacing
	md.Map[newPos.Y][newPos.X] = md.MyFacing
	return newPos, true
}

// MovePlayer moves the player's position into given direction (while can) and turns (if turn instruction)
func (md *MonkeysDescription) MovePlayer(move Movement) {
	haveMoved := true
	for move.Steps > 0 && haveMoved {
		md.MyPosition, haveMoved = md.canMove()
		move.Steps--
	}
	md.TurnPlayer(move.TurnDirection)
}

func calculateSide(currentPos Position) int {
	currentSide := -1
	if currentPos.Y >= 3*cubeSize {
		currentSide = 5
	} else if currentPos.Y >= cubeSize && currentPos.Y < cubeSize*2 {
		currentSide = 2
	} else if currentPos.Y < cubeSize {
		currentSide = 0
		if currentPos.X >= cubeSize {
			currentSide = 1
		}
	} else {
		currentSide = 3
		if currentPos.X >= cubeSize {
			currentSide = 4
		}
	}
	return currentSide
}

// String returns the textual map representation with player's last position marked with a facing marker.
func (md *MonkeysDescription) String() string {
	out := ""
	for y := range md.Map {
		for x := range md.Map[y] {
			if md.MyPosition.Y == y && md.MyPosition.X == x {
				out += md.MyFacing
				continue
			}
			out += md.Map[y][x]
		}
		out += "\n"
	}
	return out
}

// NewMonkeyDescription constructs a new map from Monkey's description.
func NewMonkeyDescription(desc string, mapType string) *MonkeysDescription {
	monkeyDesc := &MonkeysDescription{
		MyFacing: faceRight,
	}

	switch mapType {
	case flatMap:
		monkeyDesc.StepperFunc = func(currentPos Position, facing string) (Position, string) {
			return currentPos.add(relativeFacing[facing]), facing
		}
	case cubeMap:
		monkeyDesc.StepperFunc = func(currentPos Position, facing string) (Position, string) {
			currentSide := calculateSide(currentPos)
			if currentSide < 0 {
				panic(fmt.Sprintf("Current side is <0 @ %v", currentPos))
			}

			// fmt.Printf(" on side (stepping '%s'): %d (%v) \n", facing, currentSide+1, currentPos)
			naiveStep := currentPos.add(relativeFacing[facing])
			if naiveStep.X < cubeMinBoundaries[currentSide].X || naiveStep.Y < cubeMinBoundaries[currentSide].Y || naiveStep.X > cubeMaxBoundaries[currentSide].X || naiveStep.Y > cubeMaxBoundaries[currentSide].Y {
				// Move onto different surface -> can change facing there
				// fmt.Printf("new side (from going '%s'): %v \n", facing, currentPos)
				switch currentSide {
				case 0:
					switch facing {
					case faceRight:
						return Position{cubeMinBoundaries[1].X, currentPos.Y}, faceRight
					case faceDown:
						return Position{currentPos.X, cubeMinBoundaries[2].Y}, faceDown
					case faceLeft:
						return Position{cubeMinBoundaries[3].X, cubeMaxBoundaries[0].Y - currentPos.Y + cubeMinBoundaries[3].Y}, faceRight
					case faceUp:
						return Position{0, cubeMinBoundaries[5].Y + currentPos.X}, faceRight
					}
				case 1:
					switch facing {
					case faceRight:
						return Position{cubeMaxBoundaries[4].X, cubeMaxBoundaries[1].Y - currentPos.Y + cubeMinBoundaries[4].Y}, faceLeft
					case faceDown:
						return Position{cubeMaxBoundaries[2].X, currentPos.X - cubeMinBoundaries[1].X + cubeMinBoundaries[2].Y}, faceLeft
					case faceLeft:
						return Position{cubeMaxBoundaries[0].X, currentPos.Y}, faceLeft
					case faceUp:
						return Position{currentPos.X - cubeMinBoundaries[1].X, cubeMaxBoundaries[5].Y}, faceUp
					}
				case 2:
					switch facing {
					case faceRight:
						return Position{currentPos.Y - cubeMinBoundaries[2].Y + cubeMinBoundaries[1].X, cubeMaxBoundaries[1].Y}, faceUp
					case faceDown:
						return Position{currentPos.X + cubeMinBoundaries[4].X, cubeMinBoundaries[4].Y}, faceDown
					case faceLeft:
						return Position{currentPos.Y - cubeMinBoundaries[2].Y, cubeMinBoundaries[3].Y}, faceDown
					case faceUp:
						return Position{currentPos.X, cubeMaxBoundaries[0].Y}, faceUp
					}
				case 3:
					switch facing {
					case faceRight:
						return Position{cubeMinBoundaries[4].X, currentPos.Y}, faceRight
					case faceDown:
						return Position{currentPos.X, cubeMinBoundaries[5].Y}, faceDown
					case faceLeft:
						return Position{0, cubeMaxBoundaries[3].Y - currentPos.Y}, faceRight
					case faceUp:
						return Position{0, currentPos.X + cubeMinBoundaries[2].Y}, faceRight
					}
				case 4:
					switch facing {
					case faceRight:
						return Position{cubeMaxBoundaries[1].X, cubeMaxBoundaries[4].Y - currentPos.Y}, faceLeft
					case faceDown:
						return Position{cubeMaxBoundaries[5].X, currentPos.X - cubeMinBoundaries[4].X + cubeMinBoundaries[5].Y}, faceLeft
					case faceLeft:
						return Position{cubeMaxBoundaries[3].X, currentPos.Y}, faceLeft
					case faceUp:
						return Position{currentPos.X - cubeMinBoundaries[4].X, cubeMaxBoundaries[2].Y}, faceUp
					}
				case 5:
					switch facing {
					case faceRight:
						return Position{currentPos.Y - cubeMinBoundaries[5].Y + cubeMinBoundaries[4].X, cubeMaxBoundaries[4].Y}, faceUp
					case faceDown:
						return Position{currentPos.X + cubeMinBoundaries[1].X, 0}, faceDown
					case faceLeft:
						return Position{currentPos.Y - cubeMinBoundaries[5].Y, 0}, faceDown
					case faceUp:
						return Position{currentPos.X, cubeMaxBoundaries[3].Y}, faceUp
					}
				}
			}
			return naiveStep, facing
		}
	}

	// Convert map description to 2d array of tiles
	mapLines := strings.Split(desc, "\n")
	monkeyDesc.Map = make([][]string, len(mapLines))
	for y := range monkeyDesc.Map {
		mapLine := mapLines[y]
		if mapType == cubeMap {
			mapLine = strings.TrimSpace(mapLine)
		}
		tiles := strings.Split(mapLine, "")
		monkeyDesc.Map[y] = make([]string, len(tiles))
		copy(monkeyDesc.Map[y], tiles)
	}
	// Find starting point
	startPoint := Position{}
	for x := 0; x < len(monkeyDesc.Map[0]); x++ {
		if monkeyDesc.Map[0][x] == openTile {
			startPoint.X = x
			break
		}
	}
	monkeyDesc.MyPosition = startPoint
	return monkeyDesc
}

// ScanfMovement returns a list of movement instructions
func ScanfMovement(inst string) []Movement {
	movements := make([]Movement, 0)
	nextMove := Movement{}
	for len(inst) > 0 {
		nextMove.Steps = -1
		nextMove.TurnDirection = ""
		fmt.Sscanf(inst, "%d%s", &nextMove.Steps, &nextMove.TurnDirection)
		inst = nextMove.TurnDirection
		if len(nextMove.TurnDirection) > 1 {
			inst = nextMove.TurnDirection[1:]
			nextMove.TurnDirection = string(nextMove.TurnDirection[0])
		}
		if nextMove.Steps > -1 || nextMove.TurnDirection != "" {
			movements = append(movements, nextMove)
		}

	}
	return movements
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	parts := strings.Split(input, "\n\n")
	moves := ScanfMovement(parts[1])

	// Initial board setup
	// fmt.Println("==== BoardMap ====")
	// fmt.Println(boardMap)
	// fmt.Println("==== Moves ====")
	// for _, m := range moves {
	// 	fmt.Printf("%d%s, ", m.Steps, m.TurnDirection)
	// }
	// fmt.Println()

	if challengePart == 1 {
		monkeyMap := NewMonkeyDescription(parts[0], flatMap)
		for _, move := range moves {
			monkeyMap.MovePlayer(move)
			// fmt.Printf("=== Move %d%s ===\n", move.Steps, move.TurnDirection)
			// fmt.Println(boardMap)
		}
		var facingIndex int
		for facingIndex = range sortedFacing {
			if sortedFacing[facingIndex] == monkeyMap.MyFacing {
				break
			}
		}
		// fmt.Println(boardMap)
		finalPassword := (monkeyMap.MyPosition.Y+1)*1000 + (monkeyMap.MyPosition.X+1)*4 + facingIndex
		return finalPassword
	} else if challengePart == 2 {
		monkeyMap := NewMonkeyDescription(parts[0], cubeMap)
		// fmt.Println("==== BoardMap ====")
		// fmt.Println(monkeyMap)
		for _, move := range moves {
			monkeyMap.MovePlayer(move)
			// fmt.Printf("=== Move %d%s ===\n", move.Steps, move.TurnDirection)
			// fmt.Println(boardMap)
		}
		var facingIndex int
		for facingIndex = range sortedFacing {
			if sortedFacing[facingIndex] == monkeyMap.MyFacing {
				break
			}
		}
		// fmt.Println(monkeyMap)
		// fmt.Println(monkeyMap.MyPosition)
		// fmt.Println(monkeyMap.MyFacing)

		finalPassword := (monkeyMap.MyPosition.Y+1)*1000 + (monkeyMap.MyPosition.X+1)*4 + facingIndex
		return finalPassword
	}
	return -1
}

func main() {
	fmt.Println(runChallenge(2))
}
