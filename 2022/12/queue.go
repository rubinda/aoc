package main

import "sync"

// Vertex is an item in queue that has a given priority.
type Vertex struct {
	Node     *Node
	Priority int
}

// PriorityQueue implements a minimum priority queue.
// Items are sorted ascending according to priority (minimum priority first).
type PriorityQueue struct {
	Items []Vertex
	mutex sync.RWMutex
}

// Enqueue inserts an item into the queue in regard to it's priority.
func (pq *PriorityQueue) Enqueue(v Vertex) {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	if len(pq.Items) == 0 {
		pq.Items = append(pq.Items, v)
		return
	}
	inserted := false
	for i, u := range pq.Items {
		if v.Priority < u.Priority {
			pq.Items = append(pq.Items[:i+1], pq.Items[i:]...)
			pq.Items[i] = v
			inserted = true
			break
		}
	}
	if !inserted {
		pq.Items = append(pq.Items, v)
	}
}

// Dequeue returns the first (lowest priority) item in the queue.
func (pq *PriorityQueue) Dequeue() *Vertex {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	first := pq.Items[0]
	pq.Items = pq.Items[1:]
	return &first
}

// IsEmpty returns true if no items in queue.
func (pq *PriorityQueue) IsEmpty() bool {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	return len(pq.Items) == 0
}

// Size returns number of items in queue.
func (pq *PriorityQueue) Size() int {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	return len(pq.Items)
}

// NewPriorityQueue initializes a new minimum priority queue.
func NewMinimumPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		Items: make([]Vertex, 0),
	}
}
