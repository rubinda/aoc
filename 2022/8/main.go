package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed example.in
var input string

// markCondition is used when printing the forest to mark certain trees
type markCondition func(t *Tree) bool

// direction represents 1 step into a direction (e.g. up, right, down, left)
type direction struct {
	y int
	x int
}

// Directions are based on the top left corner
var (
	up    = direction{-1, 0}
	right = direction{0, 1}
	down  = direction{1, 0}
	left  = direction{0, -1}
)

// Tree belongs to a forest
type Tree struct {
	x           int
	y           int
	height      int
	scenicScore int
	isVisible   bool
}

// isBlocking returns true if current tree is blocking the view further from given tree
func (t *Tree) isBlocking(b *Tree) bool {
	return t.height >= b.height
}

// isVisibleOver returns true if current tree can be seen behind given tree
func (t *Tree) isVisibleOver(b *Tree) bool {
	return b.isVisible && t.height > b.height
}

// Forest represents a square collection of trees
type Forest struct {
	trees        [][]*Tree
	width        int
	depth        int
	visibleTrees []*Tree
	mostScenic   *Tree
}

// Prints out the forest and marks certain trees based on given condition function
func (f *Forest) PrintMarkedTrees(condition markCondition) {
	for _, line := range f.trees {
		for _, tree := range line {
			format := " %d "
			if condition(tree) {
				format = "[%d]"
			}
			fmt.Printf(format, tree.height)
		}
		fmt.Println()
	}
}

// markVisible sets a tree to visible status
func (f *Forest) markVisible(t *Tree) {
	alreadyVisible := t.isVisible
	t.isVisible = true
	if !alreadyVisible {
		f.visibleTrees = append(f.visibleTrees, t)
	}
}

// checkVisibility manually checks if the tree is highest in given direction and thus visible
func (f *Forest) checkVisibility(tree *Tree, d direction) {
	neighbour := f.getNeighbour(tree, d)
	for neighbour != nil {
		if f.isOnEdge(neighbour) {
			f.markVisible(neighbour)
		}
		if neighbour.isBlocking(tree) {
			break
		}
		if tree.isVisibleOver(neighbour) {
			f.markVisible(tree)
			break
		}
		neighbour = f.getNeighbour(neighbour, d)
	}
}

// isOnEdge returns true if tree is on the perimiter of the forest
func (f *Forest) isOnEdge(tree *Tree) bool {
	return tree.x == 0 || tree.y == 0 || tree.x == (f.depth-1) || tree.y == (f.width-1)
}

// MarkVisibleTrees will set visible status on every tree in forest that fulfills the given conditions
// Read Advent Of Code 2022 Day 8 Part 1 for the necessary conditions
func (f *Forest) MarkVisibleTrees() {
	// Becase we traverse from the top left corner, we can store the highest neighbours in these directions
	// and use them without for instant checking if tree is visible from these directions
	topHighestLine := make([]*Tree, f.depth)
	leftHighestLine := make([]*Tree, f.width)

	// For forest width
	for y, line := range f.trees {
		// For forest depth
		for x, tree := range line {
			if tree.isVisible {
				continue
			}
			// The perimeter around the edge is always visible
			if f.isOnEdge(tree) {
				if x == 0 {
					// Instantiate the left highest tree
					leftHighestLine[y] = tree
				}
				if y == 0 {
					// Instantiate the top highest tree
					topHighestLine[x] = tree
				}
				f.markVisible(tree)
				continue
			}
			// Check left and top direction first, becase it is trivial
			if tree.height > leftHighestLine[y].height || tree.height > topHighestLine[x].height {
				f.markVisible(tree)
			} else {
				// Check into the right direction
				f.checkVisibility(tree, right)
				if !tree.isVisible {
					// Check into the bottom direction
					f.checkVisibility(tree, down)
				}
			}
			if tree.isVisible {
				// Update leftHighest and topHighest
				if !leftHighestLine[y].isBlocking(tree) {
					leftHighestLine[y] = tree
				}
				if !topHighestLine[x].isBlocking(tree) {
					topHighestLine[x] = tree
				}
			}
		}
	}
}

// viewingDistance returns how many trees are visible from given tree and direction
// A view into the abyss from the edge has a viewing distance of 0
// A view into a tree of the same height has a viewing distance of 1
func (f *Forest) viewingDistance(t *Tree, d direction) int {
	vd := 0
	neighbour := f.getNeighbour(t, d)
	for neighbour != nil {
		vd++
		if neighbour.isBlocking(t) {
			break
		}
		neighbour = f.getNeighbour(neighbour, d)
	}
	return vd
}

// outOfBounds returns true if given position is invalid for current forest
func (f *Forest) outOfBounds(x, y int) bool {
	return x < 0 || y < 0 || x >= f.depth || y >= f.width
}

// getNeighbour gets the next neighbouring tree in given direction. Returns nil if out of bounds
func (f *Forest) getNeighbour(t *Tree, d direction) *Tree {
	x := t.x + d.x
	y := t.y + d.y
	if f.outOfBounds(x, y) {
		return nil
	}
	return f.trees[y][x]
}

// calculateScenicScores adds a scenic score to each tree in forest
// Read Advent Of Code 2022 Day 8 Part 2 for details on how to calculate a score
func (f *Forest) calculateScenicScores() {
	directions := []direction{up, right, down, left}
	f.mostScenic = f.trees[0][0]

	for _, line := range f.trees {
		for _, tree := range line {
			if f.isOnEdge(tree) {
				tree.scenicScore = 0
				continue
			}
			scenicScore := 1
			for _, d := range directions {
				scenicScore *= f.viewingDistance(tree, d)
			}
			tree.scenicScore = scenicScore
			if tree.scenicScore > f.mostScenic.scenicScore {
				f.mostScenic = tree
			}
		}
	}
}

// Instantiate an empty forest of given size
func PrepareForest(xLen, yLen int) *Forest {
	forest := &Forest{}
	forest.depth = xLen
	forest.width = yLen
	forest.visibleTrees = make([]*Tree, 0)
	forest.trees = make([][]*Tree, yLen)
	for i := range forest.trees {
		forest.trees[i] = make([]*Tree, xLen)
	}
	return forest
}

// parseForest reads the challenge inputs. See example.in or challenge.in
func parseForest(desc string) *Forest {
	lines := strings.Split(desc, "\n")
	yLen := len(lines)
	xLen := len(strings.Split(lines[0], ""))
	forest := PrepareForest(xLen, yLen)
	for y, treeLine := range lines {
		for x, tH := range strings.Split(treeLine, "") {
			treeHeight, _ := strconv.Atoi(tH)
			forest.trees[y][x] = &Tree{x: x, y: y, height: treeHeight}
		}
	}
	return forest
}

// runChallenge returns the desired output for the days challenge.
// May print additional information to stdout
func runChallenge(challengePart int) int {
	result := -1
	forest := parseForest(input)
	if challengePart == 1 {
		forest.MarkVisibleTrees()
		result = len(forest.visibleTrees)
		// forest.PrintMarkedTrees(func(t *Tree) bool {
		// 	return t.isVisible
		// })
		fmt.Println("Number of visible trees", len(forest.visibleTrees))
	} else if challengePart == 2 {
		forest.calculateScenicScores()
		// forest.PrintMarkedTrees(func(t *Tree) bool {
		// 	return t.x == forest.mostScenic.x && t.y == forest.mostScenic.y
		// })
		fmt.Println("Most scenic tree has a score of ", forest.mostScenic.scenicScore)
		return forest.mostScenic.scenicScore
	}
	return result
}

func main() {
	runChallenge(2)
}
