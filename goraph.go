package goraph

import "errors"
import "fmt"

type Node struct {
	id    int
	Label string
}

var nilNode = Node{id: -1}

func (n Node) isNil() bool {
	return n.id == -1
}

type Edge struct {
	node1 Node
	node2 Node
}

type BipartiteGraph struct {
	Left  []Node
	Right []Node
	Edges []Edge
}

func NewBipartiteGraph(leftValues, rightValues []interface{}, neighbours func(interface{}, interface{}) (bool, error)) (*BipartiteGraph, error) {
	if len(leftValues) != len(rightValues) {
		return nil, errors.New(fmt.Sprintf("left and right values have mismatched lengths: %d and %d", len(leftValues), len(rightValues)))
	}

	left := []Node{}
	for i, _ := range leftValues {
		left = append(left, Node{i, "left"})
	}

	right := []Node{}
	for j, _ := range rightValues {
		right = append(right, Node{j, "right"})
	}

	edges := []Edge{}
	for i, leftValue := range leftValues {
		for j, rightValue := range rightValues {
			neighbours, err := neighbours(leftValue, rightValue)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error determining adjacency for %v and %v: %s", leftValue, rightValue, err.Error()))
			}

			if neighbours {
				edges = append(edges, Edge{left[i], right[j]})
			}
		}
	}

	return &BipartiteGraph{left, right, edges}, nil
}

func (bg *BipartiteGraph) Neighbours(n1, n2 Node) bool {
	for _, edge := range bg.Edges {
		if (edge.node1 == n1 && edge.node2 == n2) || (edge.node1 == n2 && edge.node2 == n1) {
			return true
		}
	}

	return false
}

func (bg *BipartiteGraph) bfs(leftRight, rightLeft map[Node]Node, dist map[Node]int) bool {
	queue := []Node{}

	for _, v := range bg.Left {
		if leftRight[v].isNil() {
			dist[v] = 0
			queue = append(queue, v)
		} else {
			dist[v] = -1
		}
	}
	dist[nilNode] = -1

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		if dist[v] != -1 && (dist[v] < dist[nilNode] || dist[nilNode] == -1) {
			for _, u := range bg.Right {
				w := rightLeft[u]

				if bg.Neighbours(v, u) && dist[w] == -1 {
					dist[w] = dist[v] + 1
					queue = append(queue, w)
				}
			}
		}
	}

	return dist[nilNode] != -1
}

func (bg *BipartiteGraph) dfs(v Node, leftRight, rightLeft map[Node]Node, dist map[Node]int) bool {
	if !v.isNil() {
		for _, u := range bg.Right {
			w := rightLeft[u]

			if bg.Neighbours(v, u) && dist[w] == dist[v]+1 && bg.dfs(w, leftRight, rightLeft, dist) {
				rightLeft[u] = v
				leftRight[v] = u
				return true
			}
		}

		dist[v] = -1
		return false
	}

	return true
}

func (bg *BipartiteGraph) LargestMatchingSize() int {
	leftRight := make(map[Node]Node)
	rightLeft := make(map[Node]Node)
	dist := make(map[Node]int)

	for _, v := range bg.Left {
		leftRight[v] = nilNode
	}
	for _, u := range bg.Right {
		rightLeft[u] = nilNode
	}
	leftRight[nilNode] = nilNode
	rightLeft[nilNode] = nilNode

	matching := 0

	for bg.bfs(leftRight, rightLeft, dist) {
		for _, v := range bg.Left {
			if leftRight[v].isNil() && bg.dfs(v, leftRight, rightLeft, dist) {
				matching = matching + 1
			}
		}
	}

	return matching
}
