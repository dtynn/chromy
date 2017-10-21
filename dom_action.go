package chromy

import (
	"context"

	"github.com/mafredri/cdp/protocol/dom"
)

func WaitNode(selector string) Action {
	return NewQuery(QuerySelector(selector))
}

func WaitNodeAll(selector string) Action {
	return NewQuery(QuerySelectorAll(selector))
}

func GetNode(selector string, node *Node) Action {
	return NewQuery(QuerySelector(selector), QueryNode(node))
}

func GetNodeAll(selector string, nodes *[]*Node) Action {
	return NewQuery(QuerySelectorAll(selector), QueryNodes(nodes))
}

func OnNode(selector string, fn func(ctx context.Context, t *Target, node *Node) error) Action {
	return NewQuery(QuerySelector(selector), QueryAfter(fn))
}

func OnNodeAll(selector string, fn func(ctx context.Context, t *Target, nodes ...*Node) error) Action {
	return NewQuery(QuerySelectorAll(selector), QueryAfterAll(fn))
}

func Focus(selector string) Action {
	return OnNode(selector, func(ctx context.Context, t *Target, node *Node) error {
		err := t.Client().DOM.Focus(ctx, dom.NewFocusArgs().SetNodeID(node.NodeID))
		if err != nil {
			return err
		}

		return nil
	})
}

func Blur(selector string) Action {
	return OnNode(selector, func(ctx context.Context, t *Target, node *Node) error {
		cli := t.Client()
		objectID, err := nodeIDToRemoteObjectID(ctx, cli, node.NodeID)
		if err != nil {
			return err
		}

		return callFuncOnRemoteObject(ctx, cli, objectID, "function() { this.click() }", nil, nil)
	})
}

func Click(selector string) Action {
	return OnNode(selector, func(ctx context.Context, t *Target, node *Node) error {
		cli := t.Client()
		objectID, err := nodeIDToRemoteObjectID(ctx, cli, node.NodeID)
		if err != nil {
			return err
		}

		return callFuncOnRemoteObject(ctx, cli, objectID, "function() { this.click() }", nil, nil)
	})
}
