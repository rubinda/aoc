package main

import (
	_ "embed"
	"fmt"
	"strings"
)

var (
	//go:embed example.in
	input string
	// Represents relative directions to neighbours in 3D space (6 degrees of freedom).
	neighbours6DOF = []Voxel{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{-1, 0, 0},
		{0, -1, 0},
		{0, 0, -1},
	}
)

// Element for flood fill.
const (
	Air = iota
	Water
	Lava
)

// Voxel is a 1x1x1 cube in 3D space.
type Voxel struct {
	x, y, z int
}

// add returns new voxel with the sum of coordinates.
func (ld Voxel) add(d Voxel) Voxel {
	return Voxel{ld.x + d.x, ld.y + d.y, ld.z + d.z}
}

// Stack represents a stack-like data structure.
type Stack struct {
	voxels []Voxel
}

// Push adds the item to the end of the stack.
func (q *Stack) Push(v Voxel) {
	q.voxels = append(q.voxels, v)
}

// Pop returns the first item in stack.
func (q *Stack) Pop() Voxel {
	if !q.IsNotEmpty() {
		panic("Queue is empty")
	}
	v := q.voxels[0]
	q.voxels = q.voxels[1:]
	return v
}

// IsNotEmpty returns true if stack has items left.
func (q *Stack) IsNotEmpty() bool {
	return len(q.voxels) > 0
}

// NewStack returns a new stack.
func NewStack() *Stack {
	return &Stack{
		voxels: make([]Voxel, 0),
	}
}

// Bucket is a structure that contains lava droplets.
type Bucket struct {
	contents      [][][]int
	dimensionSize int
}

// MaterialAt returns the voxel material at given position.
// Returns -1 if invalid position given.
func (b *Bucket) MaterialAt(p Voxel) int {
	if p.x < 0 || p.y < 0 || p.z < 0 || p.x >= b.dimensionSize || p.y >= b.dimensionSize || p.z >= b.dimensionSize {
		return -1
	}
	return b.contents[p.x][p.y][p.z]
}

func (b *Bucket) SetMaterial(p Voxel, material int) {
	b.contents[p.x][p.y][p.z] = material
}

// ParseDroplets reads a csv string into a list of lava droplets. Returns map with droplet coordinates as keys.
func ParseDroplets(csvCoordinates string) map[Voxel]int {
	points := strings.Split(csvCoordinates, "\n")
	occupiedSpaces := make(map[Voxel]int)
	for i := range points {
		droplet := Voxel{}
		fmt.Sscanf(points[i], "%d,%d,%d", &droplet.x, &droplet.y, &droplet.z)
		// For bucket filling, we want water to flow beneath the droplets, so any coordinate 0 should move to 1
		droplet.x++
		droplet.y++
		droplet.z++
		occupiedSpaces[droplet] = 1
	}
	return occupiedSpaces
}

// FillBucket creates a new 3D bucket filled with Lava droplets at given positions.
func FillBucket(lavaDroplets map[Voxel]int, dimenzionSize int) *Bucket {
	bucket := &Bucket{dimensionSize: dimenzionSize}
	bucket.contents = make([][][]int, dimenzionSize)
	for x := 0; x < dimenzionSize; x++ {
		bucket.contents[x] = make([][]int, dimenzionSize)
		for y := 0; y < dimenzionSize; y++ {
			bucket.contents[x][y] = make([]int, dimenzionSize)
			for z := 0; z < dimenzionSize; z++ {
				bucket.contents[x][y][z] = Air
				if _, exists := lavaDroplets[Voxel{x, y, z}]; exists {
					bucket.contents[x][y][z] = Lava
				}
			}
		}
	}
	return bucket
}

// checkNeighboursForLava check air spaces for neighbouring lava.
// Skips any position that is not air.
// Adds candidates to unvisited stack if air is encountered.
func (bucket *Bucket) checkNeigboursForLava(position Voxel, airGaps *Stack) int {
	if bucket.MaterialAt(position) != Air {
		return 0
	}
	lavaFound := 0
	for _, d := range neighbours6DOF {
		neighbour := position.add(d)
		material := bucket.MaterialAt(neighbour)
		if material < 0 || material == Water {
			continue
		}

		if material == Lava {
			lavaFound++
		} else if material == Air {
			airGaps.Push(neighbour)
		}

	}
	bucket.SetMaterial(position, Water)
	return lavaFound
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	lavaDroplets := ParseDroplets(input)
	surface := 0
	if challengePart == 1 {
		for droplet := range lavaDroplets {
			neighbours := 0
			for _, dir := range neighbours6DOF {
				if _, exists := lavaDroplets[droplet.add(dir)]; exists {
					neighbours++
				}
			}
			lavaDroplets[droplet] = 6 - neighbours
			surface += lavaDroplets[droplet]
		}
	} else if challengePart == 2 {
		bucket := FillBucket(lavaDroplets, 25)
		unvisited := NewStack()
		unvisited.Push(Voxel{0, 0, 0})
		for unvisited.IsNotEmpty() {
			current := unvisited.Pop()
			surface += bucket.checkNeigboursForLava(current, unvisited)
		}
	}
	return surface
}

func main() {
	fmt.Println(runChallenge(2))
}
