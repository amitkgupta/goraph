package goraph

import "errors"
import "fmt"
import "math"

type Node struct {
	id    int
	Label string
}

type Edge struct {
	node1 Node
	node2 Node
}

type BipartiteGraph struct {
	Left  []Node
	Right []Node
	Edges []Edge // all edges go from Left to Right nodes
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

func free(node Node, matching []Edge) bool {
	for _, edge := range matching {
		if edge.node1 == node || edge.node2 == node {
			return false
		}
	}

	return true
}

func (bg *BipartiteGraph) partition(matching []Edge) [][]Node {
	layers := [][]Node{}
	used := make(map[Node]bool)
	done := false

	currentLayer := []Node{}
	for _, node := range bg.Left {
		if free(node, matching) {
			used[node] = true
			currentLayer = append(currentLayer, node)
		}
	}
	layers = append(layers, currentLayer)

	for !done {
		lastLayer := currentLayer
		currentLayer = []Node{}

		if math.Mod(float64(len(layers)), 2.0) == 1.0 {
			for _, leftNode := range lastLayer {
				for _, rightNode := range bg.Right {
					if used[rightNode] {
						continue
					}

					edge, found := bg.findEdge(leftNode, rightNode)
					if !found || edgeInMatching(edge, matching) {
						continue
					}

					currentLayer = append(currentLayer, rightNode)
					used[rightNode] = true

					if free(rightNode, matching) {
						done = true
					}
				}
			}
		} else {
			for _, rightNode := range lastLayer {
				for _, leftNode := range bg.Left {
					if used[leftNode] {
						continue
					}

					edge, found := bg.findEdge(leftNode, rightNode)
					if !found || !edgeInMatching(edge, matching) {
						continue
					}

					currentLayer = append(currentLayer, leftNode)
					used[leftNode] = true
				}
			}

		}

		layers = append(layers, currentLayer)
	}

	return layers
}

func edgeInMatching(edge Edge, matching []Edge) bool {
	for _, e := range matching {
		if edge == e {
			return true
		}
	}

	return false
}

func (bg *BipartiteGraph) findEdge(node1, node2 Node) (Edge, bool) {
	for _, edge := range bg.Edges {
		if (edge.node1 == node1 && edge.node2 == node2) || (edge.node1 == node2 && edge.node2 == node1) {
			return edge, true
		}
	}

	return Edge{}, false
}

func (bg *BipartiteGraph) findDisjointSLAPHelper(currentNode Node, currentSLAP []Edge, currentLevel int, matching []Edge, layers [][]Node, used map[Node]bool) ([]Edge, bool) {
	used[currentNode] = true

	if currentLevel == 0 {
		return currentSLAP, true
	}

	for _, nextNode := range layers[currentLevel-1] {
		if used[nextNode] {
			continue
		}

		edge, found := bg.findEdge(currentNode, nextNode)
		if !found {
			continue
		}

		if edgeInMatching(edge, matching) == (math.Mod(float64(currentLevel), 2.0) == 1.0) {
			continue
		}

		currentSLAP = append(currentSLAP, edge)
		slap, found := bg.findDisjointSLAPHelper(nextNode, currentSLAP, currentLevel-1, matching, layers, used)
		if found {
			return slap, true
		}
		currentSLAP = currentSLAP[:len(currentSLAP)-1]
	}

	used[currentNode] = false
	return nil, false
}

func (bg *BipartiteGraph) findDisjointSLAP(start Node, matching []Edge, layers [][]Node, used map[Node]bool) ([]Edge, bool) {
	return bg.findDisjointSLAPHelper(start, []Edge{}, len(layers)-1, matching, layers, used)
}

func (bg *BipartiteGraph) maximalDisjointSLAPCollection(matching []Edge) [][]Edge {
	layers := bg.partition(matching)
	used := make(map[Node]bool)
	result := [][]Edge{}

	for _, u := range layers[len(layers)-1] {
		slap, found := bg.findDisjointSLAP(u, matching, layers, used)
		if found {
			for _, edge := range slap {
				used[edge.node1] = true
				used[edge.node2] = true
			}
			result = append(result, slap)
		}
	}

	return result
}

// assumes each input slice has no repeat elements
func symmetricDifference(edges1, edges2 []Edge) []Edge {
	edgesToInclude := make(map[Edge]bool)

	for _, edge := range edges1 {
		edgesToInclude[edge] = true
	}

	for _, edge := range edges2 {
		edgesToInclude[edge] = !edgesToInclude[edge]
	}

	edges := []Edge{}
	for edge, include := range edgesToInclude {
		if include {
			edges = append(edges, edge)
		}
	}

	return edges
}

func (bg *BipartiteGraph) LargestMatching() []Edge {
	matching := []Edge{}
	paths := bg.maximalDisjointSLAPCollection(matching)

	for len(paths) > 0 {
		for _, path := range paths {
			matching = symmetricDifference(matching, path)
		}
		paths = bg.maximalDisjointSLAPCollection(matching)
	}

	return matching
}

func (bg *BipartiteGraph) LargestMatchingSize() int {
	return len(bg.LargestMatching())
}
