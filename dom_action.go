package chromy

import (
	"context"

	"github.com/mafredri/cdp/protocol/dom"
)

func WaitNode(selector string) Action {
	return actionWrapper("waitNode", NewQuery(QuerySelector(selector)))
}

func WaitNodeAll(selector string) Action {
	return actionWrapper("waitNodeAll", NewQuery(QuerySelectorAll(selector)))
}

func GetNode(selector string, node *Node) Action {
	return actionWrapper("getNode", NewQuery(QuerySelector(selector), QueryNode(node)))
}

func GetNodeAll(selector string, nodes *[]*Node) Action {
	return actionWrapper("getNodeAll", NewQuery(QuerySelectorAll(selector), QueryNodes(nodes)))
}

func OnNode(selector string, fn func(ctx context.Context, t *Target, node *Node) error) Action {
	return actionWrapper("onNode", NewQuery(QuerySelector(selector), QueryAfter(fn)))
}

func OnNodeAll(selector string, fn func(ctx context.Context, t *Target, nodes ...*Node) error) Action {
	return actionWrapper("onNodeAll", NewQuery(QuerySelectorAll(selector), QueryAfterAll(fn)))
}

func Focus(selector string) Action {
	return actionWrapper("focus", OnNode(selector, func(ctx context.Context, t *Target, node *Node) error {
		err := t.Client().DOM.Focus(ctx, dom.NewFocusArgs().SetNodeID(node.NodeID))
		if err != nil {
			return err
		}

		return nil
	}))
}

func Blur(selector string) Action {
	return actionWrapper("blur", OnNode(selector, func(ctx context.Context, t *Target, node *Node) error {
		cli := t.Client()
		objectID, err := nodeIDToRemoteObjectID(ctx, cli, node.NodeID)
		if err != nil {
			return err
		}

		return callFuncOnRemoteObject(ctx, cli, objectID, "function() { this.click() }", nil, nil)
	}))
}

func Click(selector string) Action {
	return actionWrapper("click", OnNode(selector, func(ctx context.Context, t *Target, node *Node) error {
		cli := t.Client()
		objectID, err := nodeIDToRemoteObjectID(ctx, cli, node.NodeID)
		if err != nil {
			return err
		}

		return callFuncOnRemoteObject(ctx, cli, objectID, "function() { this.click() }", nil, nil)
	}))
}
