package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const (
	// dirType is an INode type for directories (input parsing)
	dirType string = "dir"
	// fileType is an INode type for files (input parsing)
	fileType string = "file"
	// commandSign is the starting character for a command (input parsing)
	commandSign string = "$"

	// totalSpace is the whole space on disk (part 2)
	totalSpace int = 70000000
	// neededSpace is the desired free capacity on disk (part 2)
	neededSpace int = 30000000
)

//go:embed challenge.in
var input string

// INode represents a node in the filesystem. Depending on nodeType some properties can be omitted
type INode struct {
	nodeType  string
	name      string
	parentDir *INode
	size      int
	contents  []*INode
}

// condition is a function that one can define for a search over INodes
type condition func(n *INode) bool

// GetSubdir returns a subdirectory of current INode that matches in name or a new empty one
func (n *INode) GetSubdir(wanted string) *INode {
	for _, s := range n.contents {
		if s.nodeType == dirType && s.name == wanted {
			return s
		}
	}
	return EmptyDirNode(wanted, n)
}

// updateParentSize will update the directory sizes of all parent directories up to '/'
func (n *INode) updateParentSize(newFileSize int) {
	if n.parentDir == nil {
		return
	}
	n.parentDir.size += newFileSize
	n.parentDir.updateParentSize(newFileSize)
}

// FindDirs recursively finds directory INodes that match given condition
func (n *INode) FindDirs(isMatch condition) []*INode {
	matches := make([]*INode, 0)
	if isMatch(n) {
		matches = append(matches, n)
	}
	for _, s := range n.contents {
		matches = append(matches, s.FindDirs(isMatch)...)
	}
	return matches
}

// PrintSubTree prints out the parsed filesystem structure
func (n *INode) PrintSubTree(level int) {
	spaces := strings.Repeat(" ", level)
	fmt.Printf("%s- %s (%s, %d)\n", spaces, n.name, n.nodeType, n.size)
	for _, s := range n.contents {
		s.PrintSubTree(level + 2)
	}
}

// NewINode parses a 'ls' command input line into a directory or file and appends it to the parent
func NewINode(desc string, parent *INode) *INode {
	n := &INode{}

	parts := strings.Split(desc, " ")
	n.name = parts[1]
	n.parentDir = parent
	parent.contents = append(parent.contents, n)

	if parts[0] == dirType {
		n.nodeType = dirType
		n.contents = make([]*INode, 0)
	} else if size, err := strconv.Atoi(parts[0]); err == nil {
		n.nodeType = fileType
		n.size = size
		n.updateParentSize(size)
	} else {
		panic("Unknown desc string: " + desc)
	}
	return n
}

// EmptyDirNode creates a new empty directory INode.
// Parameter parent is optional
func EmptyDirNode(name string, parent *INode) *INode {
	n := &INode{}
	n.nodeType = dirType
	n.name = name
	n.parentDir = parent
	n.contents = make([]*INode, 0)
	if parent != nil {
		parent.contents = append(parent.contents, n)
	}
	return n
}

// parseInput creates a hierarchical structure of INodes.
// Returns the root ('/') node.
// See example.in or challenge.in for desired input.
func parseInput(in string) *INode {
	var root *INode
	var cwd *INode
	for _, inst := range strings.Split(input, "\n") {
		parts := strings.Fields(inst)
		if parts[0] == commandSign {
			switch parts[1] {
			case "cd":
				newCwd := parts[2]
				if newCwd == ".." {
					cwd = cwd.parentDir
					continue
				}
				var newDir *INode
				if cwd != nil {
					newDir = cwd.GetSubdir(newCwd)
				} else {
					newDir = EmptyDirNode(parts[2], cwd)
				}
				if newCwd == "/" {
					root = newDir
				}
				cwd = newDir
			case "ls":
				// No clue what to do with this information ...
				continue
			}
		} else {
			// Should be contents of cwd until a '$' pops up
			NewINode(inst, cwd)
		}
	}
	return root
}

// runChallenge returns the desired output for the days challenge.
// May print additional information to stdout
func runChallenge(challengePart int) int {
	result := 0
	root := parseInput(input)

	if challengePart == 1 {
		matches := root.FindDirs(func(n *INode) bool {
			return n.nodeType == dirType && n.size <= 100000
		})
		fmt.Printf("Found %d directories that match: \n", len(matches))
		for _, m := range matches {
			fmt.Printf("%10d %s \n", m.size, m.name)
			result += m.size
		}
		return result
	}

	if challengePart == 2 {
		minimumSize := neededSpace - (totalSpace - root.size)
		if minimumSize > 0 {
			matches := root.FindDirs(func(n *INode) bool {
				return n.nodeType == dirType && n.size >= minimumSize
			})
			fmt.Printf("Found %d directories that match: \n", len(matches))
			// Find smallest directory that is > minimumSize
			smallest := -1
			for _, d := range matches {
				fmt.Printf("%10d %s \n", d.size, d.name)
				if d.size > minimumSize && (smallest < 0 || d.size < smallest) {
					smallest = d.size
				}
			}
			return smallest
		}
	}
	return -1
}

func main() {
	res := runChallenge(2)
	fmt.Println(res)
}
