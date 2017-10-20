package chromy

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

func OnNode(selector string, fn func(node *Node) error) Action {
	return NewQuery(QuerySelector(selector), QueryAfter(fn))
}

func OnNodeAll(selector string, fn func(nodes ...*Node) error) Action {
	return NewQuery(QuerySelectorAll(selector), QueryAfterAll(fn))
}
