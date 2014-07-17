package bipartitegraph

import "errors"
import "fmt"

import . "github.com/amitkgupta/goraph/node"
import . "github.com/amitkgupta/goraph/edge"
import "github.com/amitkgupta/goraph/util"

type BipartiteGraph struct {
	Left  NodeCollection
	Right NodeCollection
	Edges EdgeCollection // all edges go from Left to Right nodes
}

func NewBipartiteGraph(leftValues, rightValues []interface{}, neighbours func(interface{}, interface{}) (bool, error)) (*BipartiteGraph, error) {
	left := NodeCollection{}
	for i, _ := range leftValues {
		left = append(left, Node{i})
	}

	right := NodeCollection{}
	for j, _ := range rightValues {
		right = append(right, Node{j + len(left)})
	}

	edges := EdgeCollection{}
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

func (bg *BipartiteGraph) LargestMatchingSize() int {
	return len(bg.LargestMatching())
}

func (bg *BipartiteGraph) LargestMatching() EdgeCollection {
	matching := EdgeCollection{}
	paths := bg.maximalDisjointSLAPCollection(matching)

	for len(paths) > 0 {
		for _, path := range paths {
			matching = matching.SymmetricDifference(path)
		}
		paths = bg.maximalDisjointSLAPCollection(matching)
	}

	return matching
}

func (bg *BipartiteGraph) maximalDisjointSLAPCollection(matching EdgeCollection) []EdgeCollection {
	layers := bg.partition(matching)
	used := make(map[Node]bool)
	result := []EdgeCollection{}

	for _, u := range layers[len(layers)-1] {
		slap, found := bg.findDisjointSLAP(u, matching, layers, used)
		if found {
			for _, edge := range slap {
				used[edge.Node1] = true
				used[edge.Node2] = true
			}
			result = append(result, slap)
		}
	}

	return result
}

func (bg *BipartiteGraph) findDisjointSLAP(
	start Node,
	matching EdgeCollection,
	layers []NodeCollection,
	used map[Node]bool,
) ([]Edge, bool) {
	return bg.findDisjointSLAPHelper(start, EdgeCollection{}, len(layers)-1, matching, layers, used)
}

func (bg *BipartiteGraph) findDisjointSLAPHelper(
	currentNode Node,
	currentSLAP EdgeCollection,
	currentLevel int,
	matching EdgeCollection,
	layers []NodeCollection,
	used map[Node]bool,
) (EdgeCollection, bool) {
	used[currentNode] = true

	if currentLevel == 0 {
		return currentSLAP, true
	}

	for _, nextNode := range layers[currentLevel-1] {
		if used[nextNode] {
			continue
		}

		edge, found := bg.Edges.FindByNodes(currentNode, nextNode)
		if !found {
			continue
		}

		if matching.Contains(edge) == util.Odd(currentLevel) {
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

func (bg *BipartiteGraph) partition(matching EdgeCollection) []NodeCollection {
	layers := []NodeCollection{}
	used := make(map[Node]bool)
	done := false

	currentLayer := NodeCollection{}
	for _, node := range bg.Left {
		if matching.Free(node) {
			used[node] = true
			currentLayer = append(currentLayer, node)
		}
	}
	layers = append(layers, currentLayer)

	for !done {
		lastLayer := currentLayer
		currentLayer = NodeCollection{}

		if util.Odd(len(layers)) {
			for _, leftNode := range lastLayer {
				for _, rightNode := range bg.Right {
					if used[rightNode] {
						continue
					}

					edge, found := bg.Edges.FindByNodes(leftNode, rightNode)
					if !found || matching.Contains(edge) {
						continue
					}

					currentLayer = append(currentLayer, rightNode)
					used[rightNode] = true

					if matching.Free(rightNode) {
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

					edge, found := bg.Edges.FindByNodes(leftNode, rightNode)
					if !found || !matching.Contains(edge) {
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
