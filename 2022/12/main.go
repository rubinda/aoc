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

// canTravel returns true if Edge from source -> destionation exists.
func canTravel(source, destination *Node, challengePart int) bool {
	if challengePart == 2 {
		return float64(source.Weight)-float64(destination.Weight) <= 1
	}
	return float64(source.Weight)-float64(destination.Weight) >= -1
}

// ensureNode creates a new Node for given point if it doesn't exist yet.
func ensureNode(g *Graph, nodes map[Point]*Node, p Point, heightMarker string) {
	if _, found := nodes[p]; !found {
		n := &Node{HeightMarker: heightMarker, Weight: convertToWeight(heightMarker), Coordinates: p}
		nodes[p] = n
		g.AddNode(n)
	}
}

// CreateGraph instantiates a graph from challenge data.
func CreateGraph(data [][]string, challengePart int) (graph *Graph, startNode, endNode *Node) {
	graph = NewGraph()
	height := len(data)
	width := len(data[0])
	nodes := make(map[Point]*Node)
	for y, row := range data {
		for x, heightMarker := range row {
			p := Point{x, y}
			ensureNode(graph, nodes, p, heightMarker)
			for _, neighbour := range p.Neighbours(width, height) {
				ensureNode(graph, nodes, neighbour, data[neighbour.y][neighbour.x])
				if canTravel(nodes[p], nodes[neighbour], challengePart) {
					graph.AddEdge(nodes[p], nodes[neighbour], 1)
				}
			}
			if heightMarker == startingMarker {
				startNode = nodes[p]
			} else if heightMarker == endMarker {
				endNode = nodes[p]
			}
		}
	}
	return
}

// ShortestPath find the shortest path between source and sink (Dijkstra).
// Returns pathing map and distance map.
func ShortestPath(source *Node, sink *Node, g *Graph) (map[*Node]*Node, int, map[*Node]int) {
	visited := make(map[*Node]bool)
	dist := make(map[*Node]int)
	prev := make(map[*Node]*Node)

	pq := NewMinimumPriorityQueue()
	start := Vertex{
		Node:     source,
		Priority: 0,
	}
	for _, n := range g.Nodes {
		dist[n] = math.MaxInt64
	}
	dist[source] = start.Priority

	pq.Enqueue(start)
	for !pq.IsEmpty() {
		v := pq.Dequeue()

		if visited[v.Node] {
			continue
		}
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
	return prev, dist[sink], dist
}

//go:embed example.in
var input string

// parseInput reads the input string and returns DEM-like grid.
func parseInput() [][]string {
	lines := strings.Split(input, "\n")
	dem := make([][]string, len(lines))
	for y, line := range lines {
		dem[y] = make([]string, len(line))
		for x, c := range line {
			dem[y][x] = string(c)
		}
	}
	return dem
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) int {
	result := -1
	data := parseInput()
	graph, startNode, endNode := CreateGraph(data, challengePart)
	if challengePart == 1 {
		_, result, _ = ShortestPath(startNode, endNode, graph)
	} else if challengePart == 2 {
		startNode = endNode
		_, _, distances := ShortestPath(endNode, nil, graph)
		// Find which Node that fulfills elevation requirement is closest to endNode
		var closest *Node
		wantedElevation := convertToWeight("a")
		for _, n := range graph.Nodes {
			if n.Weight != wantedElevation {
				continue
			}
			if closest == nil || distances[n] < distances[closest] {
				closest = n
			}
		}
		result = distances[closest]
	}

	return result
}

func main() {
	fmt.Println(runChallenge(2))
}
