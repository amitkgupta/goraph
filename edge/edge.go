package edge

import . "github.com/amitkgupta/goraph/node"

type Edge struct {
	Node1 Node
	Node2 Node
}

type EdgeCollection []Edge

func (ec EdgeCollection) Free(node Node) bool {
	for _, e := range ec {
		if e.Node1 == node || e.Node2 == node {
			return false
		}
	}

	return true
}

func (ec EdgeCollection) Contains(edge Edge) bool {
	for _, e := range ec {
		if e == edge {
			return true
		}
	}

	return false
}

func (ec EdgeCollection) FindByNodes(node1, node2 Node) (Edge, bool) {
	for _, e := range ec {
		if (e.Node1 == node1 && e.Node2 == node2) || (e.Node1 == node2 && e.Node2 == node1) {
			return e, true
		}
	}

	return Edge{}, false
}

func (ec EdgeCollection) SymmetricDifference(ec2 EdgeCollection) EdgeCollection {
	edgesToInclude := make(map[Edge]bool)

	for _, e := range ec {
		edgesToInclude[e] = true
	}

	for _, e := range ec2 {
		edgesToInclude[e] = !edgesToInclude[e]
	}

	result := EdgeCollection{}
	for e, include := range edgesToInclude {
		if include {
			result = append(result, e)
		}
	}

	return result
}
