package chromy

import (
	"fmt"

	"github.com/mafredri/cdp/protocol/dom"
)

func newNode(node dom.Node) *Node {
	return &Node{
		Node: node,
	}
}

type Node struct {
	dom.Node
}

func (n *Node) Attr(attr string) string {
	for i := 0; i < len(n.Node.Attributes); i += 2 {
		if n.Node.Attributes[i] == attr {
			return n.Node.Attributes[i+1]
		}
	}

	return ""
}

type nodeModifiers []func(*Node) error

func (n nodeModifiers) apply(node *Node) {
	if len(n) == 0 {
		return
	}

	for _, m := range n {
		m(node)
	}
}

func attributeModified(attr, val string) func(*Node) error {
	return func(n *Node) error {
		for i := 0; i < len(n.Attributes); i += 2 {
			if n.Attributes[i] == attr {
				n.Attributes[i+1] = val
				return nil
			}
		}

		n.Attributes = append(n.Attributes, attr, val)
		return nil
	}
}

func attributeRemoved(attr string) func(*Node) error {
	return func(n *Node) error {
		for i := 0; i < len(n.Attributes); i += 2 {
			if n.Attributes[i] == attr {
				n.Attributes = append(n.Attributes[:i], n.Attributes[i+2:]...)
				return nil
			}
		}

		return fmt.Errorf("[DOM.AttributeRemoved] attr %s missing for %d", attr, n.NodeID)
	}
}

func setChildNode(d *DOM, nodes []dom.Node) func(*Node) error {
	return func(n *Node) error {
		n.Children = nodes
		for _, one := range nodes {
			d.walkNodes(newNode(one))
		}

		return nil
	}
}

func insertNode(prevID dom.NodeID, nodes []dom.Node, node dom.Node) []dom.Node {
	for i, n := range nodes {
		if n.NodeID == prevID {
			res := make([]dom.Node, len(nodes)+1)
			copy(res[:i+1], nodes[:i+1])
			res[i+1] = node
			copy(res[i+2:], nodes[i+1:])
			return res
		}
	}

	return append(nodes, node)
}

func removeNode(nodes []dom.Node, nodeID dom.NodeID) []dom.Node {
	for i, one := range nodes {
		if one.NodeID == nodeID {
			return append(nodes[:i], nodes[i+1:]...)
		}
	}

	return nodes
}

func shadowRootPopped(nodeID dom.NodeID) func(*Node) error {
	return func(n *Node) error {
		n.ShadowRoots = removeNode(n.ShadowRoots, nodeID)

		return nil
	}
}

func shadowRootPushed(node dom.Node) func(*Node) error {
	return func(n *Node) error {
		n.ShadowRoots = append(n.ShadowRoots, node)
		return nil
	}
}

func childNodeRemoved(childID dom.NodeID) func(*Node) error {
	return func(n *Node) error {
		n.Children = removeNode(n.Children, childID)

		return nil
	}
}

func childNodeInserted(prevID dom.NodeID, node dom.Node) func(*Node) error {
	return func(n *Node) error {
		n.Children = insertNode(prevID, n.Children, node)
		return nil
	}
}

func pseudoElementAdded(node dom.Node) func(*Node) error {
	return func(n *Node) error {
		n.PseudoElements = append(n.PseudoElements, node)
		return nil
	}
}

func pseudoElementRemoved(nodeID dom.NodeID) func(*Node) error {
	return func(n *Node) error {
		n.PseudoElements = removeNode(n.PseudoElements, nodeID)

		return nil
	}
}

func characterDataModified(value string) func(*Node) error {
	return func(n *Node) error {
		n.Value = &value
		return nil
	}
}

func childNodeCountUpdated(count int) func(*Node) error {
	return func(n *Node) error {
		n.ChildNodeCount = &count
		return nil
	}
}

func distributedNodesUpdated(nodes []dom.BackendNode) func(*Node) error {
	return func(n *Node) error {
		n.DistributedNodes = nodes
		return nil
	}
}
