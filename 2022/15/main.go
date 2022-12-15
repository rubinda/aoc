package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
)

var (
	//go:embed challenge.in
	input string
)

// Relative directions based on (0,0) being the upper left corner
var (
	relativeUpRight   = Point{1, -1}
	relativeDownRight = Point{1, 1}
	relativeDownLeft  = Point{-1, 1}
	relativeUpLeft    = Point{-1, -1}

	diamondDirections = []Point{relativeUpRight, relativeDownRight, relativeDownLeft, relativeUpLeft}
)

// Challenge constants (can differ from example and challenge)
const (
	challenge1Y               = 2000000
	challenge2SearchMin       = 0
	challenge2SearchMax       = 4000000
	tuningFrequencyMultiplier = 4000000
)

// max returns the highest value from given.
func max(items ...int) int {
	highest := items[0]
	for i := 1; i < len(items); i++ {
		if items[i] > highest {
			highest = items[i]
		}
	}
	return highest
}

// min returns the lowest value from given.
func min(items ...int) int {
	lowest := items[0]
	for i := 1; i < len(items); i++ {
		if items[i] < lowest {
			lowest = items[i]
		}
	}
	return lowest
}

// Point represents a point in a coordinate system
type Point struct {
	x int
	y int
}

// DistanceTo returns the Manhattan distance between given points.
func (a Point) DistanceTo(b Point) int {
	return int(math.Abs(float64(a.x)-float64(b.x)) + math.Abs(float64(a.y)-float64(b.y)))
}

// add sums each coordinate of points.
func (a Point) add(b Point) Point {
	return Point{a.x + b.x, a.y + b.y}
}

// Sensor represents a cave object that can detect the closest beacon.
type Sensor struct {
	Location         Point
	ClosestBeacon    Point
	DistanceToBeacon int
}

// NewSensor parses a challenge input line into a sensor.
func NewSensor(sensorDesc string) Sensor {
	// e.g. sensorDesc = "Sensor at x=2, y=18: closest beacon is at x=-2, y=15"
	sensor := Sensor{}
	sensor.ClosestBeacon = Point{}
	sensor.Location = Point{}
	fmt.Sscanf(sensorDesc, `Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d`, &sensor.Location.x, &sensor.Location.y, &sensor.ClosestBeacon.x, &sensor.ClosestBeacon.y)
	sensor.DistanceToBeacon = sensor.Location.DistanceTo(sensor.ClosestBeacon)
	return sensor
}

// HasCoverageOver returns if sensor has coverage over point - distance is <= to closestBeacon.
func (s Sensor) HasCoverageOver(p Point) bool {
	return s.Location.DistanceTo(p) <= s.DistanceToBeacon
}

// Cave represents an imaginary 2D overview of a cave that contains sensors and beacons.
type Cave struct {
	Sensors []Sensor
	// ShallowestY represents lowest Y at least one sensor can cover
	ShallowestY int
	// Width represents highest Y at least one sensor can cover
	Depth int
	// ShallowestX represents lowest X at least one sensor can cover
	ShallowestX int
	// Width represents highest X at least one sensor can cover
	Width int
	// IsOccupied contains a map of occupied points
	IsOccupied map[Point]bool
}

// ParseCave returns a cave with sensors and beacons from the challenge input.
func ParseCave(challengeInput string) Cave {
	cave := Cave{}
	lines := strings.Split(challengeInput, "\n")
	cave.Sensors = make([]Sensor, len(lines))
	cave.IsOccupied = make(map[Point]bool)

	var initX, initY bool
	for i := range lines {
		s := NewSensor(lines[i])
		cave.Sensors[i] = s
		cave.IsOccupied[s.Location] = true
		cave.IsOccupied[s.ClosestBeacon] = true
		maxX := max(s.ClosestBeacon.x, s.Location.x, s.Location.x+s.DistanceToBeacon)
		minX := min(s.ClosestBeacon.x, s.Location.x, s.Location.x-s.DistanceToBeacon)
		maxY := max(s.ClosestBeacon.y, s.Location.y, s.Location.y+s.DistanceToBeacon)
		minY := min(s.ClosestBeacon.y, s.Location.y, s.Location.y-s.DistanceToBeacon)
		if maxX > cave.Width {
			cave.Width = maxX
		}
		if minX < cave.ShallowestX || !initX {
			cave.ShallowestX = minX
			initX = true
		}
		if maxY > cave.Depth {
			cave.Depth = maxY
		}
		if minY < cave.ShallowestY || !initY {
			cave.ShallowestY = minY
			initY = true
		}

	}
	return cave
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	cave := ParseCave(input)
	if challengePart == 1 {
		definitelyBeaconless := 0
		occuppied := 0
		noCoverage := 0
		for x := cave.ShallowestX; x <= cave.Depth; x++ {
			cavePoint := Point{x, challenge1Y}
			if _, isOccupied := cave.IsOccupied[cavePoint]; isOccupied {
				occuppied++
				continue
			}
			hasCoverage := false
			for _, s := range cave.Sensors {
				if s.HasCoverageOver(cavePoint) {
					hasCoverage = true
					break
				}
			}
			if hasCoverage {
				definitelyBeaconless++
			} else {
				noCoverage++
			}
		}
		fmt.Printf("=== Depth %d ===\nNo coverage: %d \n  No beacon: %d \n   Occupied: %d\n", challenge1Y, noCoverage, definitelyBeaconless, occuppied)
		return definitelyBeaconless
	} else if challengePart == 2 {
		covered := 0
		occuppied := 0
		for sI, sensor := range cave.Sensors {
			// Number of steps to take on expanded perimeter
			steps := 2*sensor.DistanceToBeacon + 3
			// Down left from leftmost edge (so we can call perimeter.add at beginning)
			perimeterPoint := sensor.Location.add(Point{-(sensor.DistanceToBeacon + 2), 1})
			for _, dir := range diamondDirections {
				for i := 0; i < steps; i++ {
					perimeterPoint = perimeterPoint.add(dir)
					if perimeterPoint.x < challenge2SearchMin || perimeterPoint.y < challenge2SearchMin || perimeterPoint.x > challenge2SearchMax || perimeterPoint.y > challenge2SearchMax {
						// Ship perimiter position if it's out of search area
						continue
					}
					if _, ok := cave.IsOccupied[perimeterPoint]; ok {
						// Sensor or beacon already exists at point
						occuppied++
						continue
					}
					hasCoverage := false
					for sI2, s := range cave.Sensors {
						if sI == sI2 {
							// We already know current sensor can't reach
							continue
						}
						if s.HasCoverageOver(perimeterPoint) {
							hasCoverage = true
							break
						}
					}
					if !hasCoverage {
						// Point has no coverage! since challenge requires only one such point the search is done
						frequency := perimeterPoint.x*tuningFrequencyMultiplier + perimeterPoint.y
						fmt.Printf("  Covered: %d \n Occupied: %d\n", covered, occuppied)
						return frequency
					} else {
						covered++
					}
				}
			}
		}
	}
	return -1
}

func main() {
	fmt.Println(runChallenge(2))
}
