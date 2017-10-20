package chromy

import (
	"context"

	"github.com/mafredri/cdp/protocol/dom"
)

func (d *DOM) appendNodeModifer(nodeID dom.NodeID, m func(*Node) error) {
	d.nodeModifiers[nodeID] = append(d.nodeModifiers[nodeID], m)
}

func (d *DOM) applyNodeModifiers(node *Node) {
	mod, ok := d.nodeModifiers[node.NodeID]
	if ok && len(mod) > 0 {
		for _, m := range mod {
			if err := m(node); err != nil {
				d.Trace(err)
			}
		}

		delete(d.nodeModifiers, node.NodeID)
	}
}

func (d *DOM) tryModifyNode(nodeID dom.NodeID, m func(*Node) error) error {
	node, ok := d.nodes[nodeID]
	if ok {
		return m(node)
	}

	d.appendNodeModifer(nodeID, m)
	return nil
}

func (d *DOM) removeNode(nodeID dom.NodeID) {
	if _, ok := d.nodes[nodeID]; ok {
		delete(d.nodes, nodeID)
	}

	if _, ok := d.nodeModifiers[nodeID]; ok {
		delete(d.nodeModifiers, nodeID)
	}
}

func (d *DOM) handleSetChildNodes() error {
	reply, err := d.events.SetChildNodesClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()
	err = d.tryModifyNode(reply.ParentID, setChildNode(d, reply.Nodes))
	d.mu.Unlock()

	return err
}

func (d *DOM) handleDocumentUpdated() error {
	_, err := d.events.DocumentUpdatedClient.Recv()
	if err != nil {
		return err
	}

	return d.updateDocument()

}

func (d *DOM) updateDocument() error {
	// reset nodes
	d.mu.Lock()
	d.root = nil
	d.nodes = make(map[dom.NodeID]*Node)
	d.nodeModifiers = make(map[dom.NodeID]nodeModifiers)
	d.mu.Unlock()

	reply, err := d.GetDocument(context.Background(), dom.NewGetDocumentArgs().SetDepth(-1).SetPierce(true))
	if err != nil {
		return err
	}

	d.mu.Lock()
	d.root = newNode(reply.Root)
	d.walkNodes(d.root)
	d.mu.Unlock()

	return nil
}

func (d *DOM) walkNodes(top *Node) {
	if top == nil {
		return
	}

	for _, one := range top.Children {
		d.walkNodes(newNode(one))
	}

	for _, one := range top.PseudoElements {
		d.walkNodes(newNode(one))
	}

	for _, one := range top.ShadowRoots {
		d.walkNodes(newNode(one))
	}

	if top.ContentDocument != nil {
		d.walkNodes(newNode(*top.ContentDocument))
	}

	if top.ImportedDocument != nil {
		d.walkNodes(newNode(*top.ImportedDocument))
	}

	if top.TemplateContent != nil {
		d.walkNodes(newNode(*top.TemplateContent))
	}

	d.nodes[top.NodeID] = top
	d.applyNodeModifiers(top)
}

func (d *DOM) handleAttributeRemoved() error {
	reply, err := d.events.AttributeRemovedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()
	err = d.tryModifyNode(reply.NodeID, attributeRemoved(reply.Name))
	d.mu.Unlock()

	return err
}

func (d *DOM) handleAttributeModified() error {
	reply, err := d.events.AttributeModifiedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()
	err = d.tryModifyNode(reply.NodeID, attributeModified(reply.Name, reply.Value))
	d.mu.Unlock()

	return err
}

func (d *DOM) handleChildNodeRemoved() error {
	reply, err := d.events.ChildNodeRemovedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()

	d.removeNode(reply.NodeID)

	err = d.tryModifyNode(reply.ParentNodeID, childNodeRemoved(reply.NodeID))

	d.mu.Unlock()

	return err
}

func (d *DOM) handleShadowRootPopped() error {
	reply, err := d.events.ShadowRootPoppedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()

	d.removeNode(reply.RootID)

	d.tryModifyNode(reply.HostID, shadowRootPopped(reply.RootID))

	d.mu.Unlock()

	return err
}

func (d *DOM) handleShadowRootPushed() error {
	reply, err := d.events.ShadowRootPushedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()
	d.walkNodes(newNode(reply.Root))
	err = d.tryModifyNode(reply.HostID, shadowRootPushed(reply.Root))
	d.mu.Unlock()

	return err
}

func (d *DOM) handleChildNodeInserted() error {
	reply, err := d.events.ChildNodeInsertedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()

	d.walkNodes(newNode(reply.Node))

	err = d.tryModifyNode(reply.ParentNodeID, childNodeInserted(reply.PreviousNodeID, reply.Node))

	d.mu.Unlock()

	return nil
}

func (d *DOM) handlePseudoElementAdded() error {
	reply, err := d.events.PseudoElementAddedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()
	d.walkNodes(newNode(reply.PseudoElement))

	err = d.tryModifyNode(reply.ParentID, pseudoElementAdded(reply.PseudoElement))

	d.mu.Unlock()

	return err
}

func (d *DOM) handlePseudoElementRemoved() error {
	reply, err := d.events.PseudoElementRemovedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()
	d.removeNode(reply.PseudoElementID)

	err = d.tryModifyNode(reply.ParentID, pseudoElementRemoved(reply.PseudoElementID))

	d.mu.Unlock()

	return err
}

func (d *DOM) handleCharacterDataModified() error {
	reply, err := d.events.CharacterDataModifiedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()

	err = d.tryModifyNode(reply.NodeID, characterDataModified(reply.CharacterData))

	d.mu.Unlock()

	return err
}

func (d *DOM) handleChildNodeCountUpdated() error {
	reply, err := d.events.ChildNodeCountUpdatedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()

	err = d.tryModifyNode(reply.NodeID, childNodeCountUpdated(reply.ChildNodeCount))

	d.mu.Unlock()

	return err
}

func (d *DOM) handleInlineStyleInvalidated() error {
	return nil
}

func (d *DOM) handleDistributedNodesUpdated() error {
	reply, err := d.events.DistributedNodesUpdatedClient.Recv()
	if err != nil {
		return err
	}

	d.mu.Lock()

	err = d.tryModifyNode(reply.InsertionPointID, distributedNodesUpdated(reply.DistributedNodes))

	d.mu.Unlock()

	return err
}
