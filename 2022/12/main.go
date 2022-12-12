package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
)

const (
	startingMarker = "S"
	endMarker      = "E"
)

// Based on (0, 0) being the upper left corner.
var (
	relativeUp    = Point{x: 0, y: -1}
	relativeRight = Point{x: 1, y: 0}
	relativeDown  = Point{x: 0, y: 1}
	relativeLeft  = Point{x: -1, y: 0}
)

// convertToWeight replaces a string representation of height with a numerical.
// (a-z) a being lowest and z being highest.
func convertToWeight(heightMarker string) int {
	if heightMarker == startingMarker {
		heightMarker = "a"
	} else if heightMarker == endMarker {
		heightMarker = "z"
	}
	return int([]byte(heightMarker)[0])
}

// Point holds coordinates where (0,0) is the upper left corner.
type Point struct {
	x int
	y int
}

func (p Point) Add(d Point) Point {
	return Point{p.x + d.x, p.y + d.y}
}

// Neighbours returns point neighbours which are inside the 'w' width and 'h' height boundaries.
func (p Point) Neighbours(w, h int) []Point {
	n := make([]Point, 0)
	directions := []Point{relativeUp, relativeRight, relativeDown, relativeLeft}
	for _, d := range directions {
		t := p.Add(d)
		if t.x >= 0 && t.y >= 0 && t.x < w && t.y < h {
			n = append(n, t)
		}
	}
	return n
}

type Node struct {
	HeightMarker string
	Weight       int
	Coordinates  Point
}

type Edge struct {
	Node   *Node
	Weight int
}

type Graph struct {
	Nodes     []*Node
	Edges     map[*Node][]*Edge
	CanTravel func(n1, n2 *Node) bool
	start     *Node
	end       *Node

	visited [][]string
}

func NewGraph(canTravel func(n1, n2 *Node) bool) *Graph {
	return &Graph{
		Nodes:     make([]*Node, 0),
		Edges:     map[*Node][]*Edge{},
		CanTravel: canTravel,
	}
}

func (g *Graph) AddNode(n *Node) {
	g.Nodes = append(g.Nodes, n)
}

func (g *Graph) AddEdge(source, destination *Node, weight int) {
	edge := Edge{
		Node:   destination,
		Weight: weight,
	}
	// Edge n1 -> n2
	// fmt.Printf("Adding edge from %s(%d, %d) to %s(%d, %d) \n", source.HeightMarker, source.Coordinates.x, source.Coordinates.y, destination.HeightMarker, destination.Coordinates.x, destination.Coordinates.y)
	g.Edges[source] = append(g.Edges[source], &edge)
}

func ensureNode(g *Graph, nodes map[Point]*Node, p Point, heightMarker string) {
	if _, found := nodes[p]; !found {
		n := &Node{HeightMarker: heightMarker, Weight: convertToWeight(heightMarker), Coordinates: p}
		nodes[p] = n
		g.AddNode(n)
	}
}

// CreateGraph instantiates a graph from challenge data.
func CreateGraph(data [][]string, travelFunc func(n1, n2 *Node) bool) *Graph {
	g := NewGraph(travelFunc)
	height := len(data)
	width := len(data[0])
	nodes := make(map[Point]*Node)
	g.visited = make([][]string, height)
	for y, row := range data {
		g.visited[y] = make([]string, width)
		for x, heightMarker := range row {
			g.visited[y][x] = heightMarker
			p := Point{x, y}
			ensureNode(g, nodes, p, heightMarker)

			for _, neighbour := range p.Neighbours(width, height) {
				ensureNode(g, nodes, neighbour, data[neighbour.y][neighbour.x])
				if heightMarker == endMarker {
					fmt.Printf("%s(%d) -> %s (%d): %v \n", nodes[p].HeightMarker, nodes[p].Weight, nodes[neighbour].HeightMarker, nodes[neighbour].Weight, g.CanTravel(nodes[p], nodes[neighbour]))
				}
				if g.CanTravel(nodes[p], nodes[neighbour]) {
					g.AddEdge(nodes[p], nodes[neighbour], 1)
				}
			}
			if heightMarker == startingMarker {
				g.start = nodes[p]
			} else if heightMarker == endMarker {
				g.end = nodes[p]
			}
		}
	}
	return g
}

func getShortestPath(startNode *Node, endNode *Node, g *Graph) (map[*Node]int, int) {
	visited := make(map[*Node]bool)
	dist := make(map[*Node]int)
	prev := make(map[*Node]*Node)

	pq := NewMinimumPriorityQueue()
	start := Vertex{
		Node:     startNode,
		Priority: 0,
	}
	for _, n := range g.Nodes {
		dist[n] = math.MaxInt64
	}
	dist[startNode] = start.Priority
	pq.Enqueue(start)

	var closestA *Node

	for !pq.IsEmpty() {
		v := pq.Dequeue()

		if visited[v.Node] {
			continue
		}
		g.visited[v.Node.Coordinates.y][v.Node.Coordinates.x] = "."
		visited[v.Node] = true
		edges := g.Edges[v.Node]

		for _, edge := range edges {
			currentDistance := dist[v.Node] + edge.Weight
			if currentDistance < dist[edge.Node] {
				dist[edge.Node] = currentDistance
				prev[edge.Node] = v.Node
				closer := Vertex{
					Node:     edge.Node,
					Priority: dist[edge.Node],
				}
				pq.Enqueue(closer)
			}

		}
	}

	for _, n := range g.Nodes {
		if n.HeightMarker == "a" && (closestA == nil || dist[closestA] > dist[n]) {
			closestA = n
		}
	}
	fmt.Println(closestA, dist[closestA])
	previous := closestA
	i := 0
	for previous != nil {
		g.visited[previous.Coordinates.y][previous.Coordinates.x] = fmt.Sprintf("%d", i)
		i++
		if prev[previous] != nil {
			fmt.Printf("%s -> %s: %v \n", prev[previous].HeightMarker, previous.HeightMarker, g.CanTravel(prev[previous], previous))
		}
		previous = prev[previous]
	}
	fmt.Printf("Graph dimenisons : %d x %d \n", len(g.visited), len(g.visited[0]))
	for y := range g.visited {
		for x := range g.visited[y] {
			fmt.Printf("%s", g.visited[y][x])
		}
		fmt.Println()
	}

	return dist, dist[closestA]
}

//go:embed challenge.in
var input string

func parseInput() [][]string {
	lines := strings.Split(input, "\n")
	out := make([][]string, len(lines))
	for y, line := range lines {
		out[y] = make([]string, len(line))
		for x, c := range line {
			out[y][x] = string(c)
		}
	}
	return out
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	result := -1
	data := parseInput()
	travelFunc := func(n1, n2 *Node) bool {
		return float64(n1.Weight)-float64(n2.Weight) >= -1
	}
	if challengePart == 2 {
		travelFunc = func(n1, n2 *Node) bool {
			return float64(n1.Weight)-float64(n2.Weight) <= 1
		}
	}
	graph := CreateGraph(data, travelFunc)
	start := graph.start
	if challengePart == 2 {
		start = graph.end
	}
	m1 := "z"
	m2 := "x"
	n1 := &Node{
		HeightMarker: m1,
		Weight:       convertToWeight(m1),
	}
	n2 := &Node{
		HeightMarker: m2,
		Weight:       convertToWeight(m2),
	}
	fmt.Printf("Can travel %s -> %s: %v \n", m1, m2, graph.CanTravel(n1, n2))

	// fmt.Printf("Start = %v \n", graph.start.Coordinates)
	// fmt.Printf("End = %v \n", graph.end.Coordinates)
	// for node, edges := range graph.Edges {
	// 	val := node.HeightMarker
	// 	if val == startingMarker {
	// 		val = "START"
	// 	} else if val == endMarker {
	// 		val = "END"
	// 	}
	// 	fmt.Printf("==== Node %s (%d, %d) ====\n", val, node.Coordinates.x, node.Coordinates.y)
	// 	for _, e := range edges {
	// 		fmt.Printf("  -> Node %s (%d, %d) distance: %d \n", e.Node.HeightMarker, e.Node.Coordinates.x, e.Node.Coordinates.y, e.Weight)
	// 	}
	// 	fmt.Println()
	// }

	_, result = getShortestPath(start, graph.end, graph)

	return result
}

func main() {
	fmt.Println(runChallenge(2))
}
