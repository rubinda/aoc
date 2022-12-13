package main

// Relative directions based on (0, 0) being the upper left corner.
var (
	relativeUp    = Point{x: 0, y: -1}
	relativeRight = Point{x: 1, y: 0}
	relativeDown  = Point{x: 0, y: 1}
	relativeLeft  = Point{x: -1, y: 0}
)

var (
	// 4 DOF movement
	directions4 = []Point{relativeUp, relativeRight, relativeDown, relativeLeft}
)

// Point holds coordinates where (0,0) is the upper left corner.
type Point struct {
	x int
	y int
}

// Add returns the coordinate-wise sum of two points.
func (p Point) Add(d Point) Point {
	return Point{p.x + d.x, p.y + d.y}
}

// Neighbours returns point 4 way neighbours which are inside the 'w' width and 'h' height boundaries.
func (p Point) Neighbours(w, h int) []Point {
	n := make([]Point, 0)
	for _, d := range directions4 {
		t := p.Add(d)
		if t.x >= 0 && t.y >= 0 && t.x < w && t.y < h {
			n = append(n, t)
		}
	}
	return n
}

// Node (or vertice) is a point belonging to a graph.
type Node struct {
	HeightMarker string
	Weight       int
	Coordinates  Point
}

// Edge represents a weighted graph connection to Node.
type Edge struct {
	Node   *Node
	Weight int
}

// Graph is a collection of Nodes connected via edges.
type Graph struct {
	// Nodes is a collection of points in the graph.
	Nodes []*Node
	// Edges contains a list of all edges originating from a given Node.
	Edges map[*Node][]*Edge
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make([]*Node, 0),
		Edges: map[*Node][]*Edge{},
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
	g.Edges[source] = append(g.Edges[source], &edge)
}
