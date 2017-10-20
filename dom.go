package chromy

import (
	"context"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
)

func newDOM(ctx context.Context, cli *cdp.Client) (*DOM, error) {
	err := cli.DOM.Enable(ctx)
	if err != nil {
		return nil, err
	}

	d := &DOM{
		DOM:           cli.DOM,
		ErrTracer:     nonErrTracer,
		nodes:         make(map[dom.NodeID]*Node),
		nodeModifiers: make(map[dom.NodeID]nodeModifiers),
	}

	defer func() {
		if err != nil {
			d.close()
		}
	}()

	d.events.SetChildNodesClient, err = cli.DOM.SetChildNodes(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.SetChildNodesClient)

	d.events.DocumentUpdatedClient, err = cli.DOM.DocumentUpdated(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.DocumentUpdatedClient)

	d.events.AttributeRemovedClient, err = cli.DOM.AttributeRemoved(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.AttributeRemovedClient)

	d.events.AttributeModifiedClient, err = cli.DOM.AttributeModified(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.AttributeModifiedClient)

	d.events.ChildNodeRemovedClient, err = cli.DOM.ChildNodeRemoved(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.ChildNodeRemovedClient)

	d.events.ShadowRootPoppedClient, err = cli.DOM.ShadowRootPopped(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.ShadowRootPoppedClient)

	d.events.ShadowRootPushedClient, err = cli.DOM.ShadowRootPushed(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.ShadowRootPushedClient)

	d.events.ChildNodeInsertedClient, err = cli.DOM.ChildNodeInserted(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.ChildNodeInsertedClient)

	d.events.PseudoElementAddedClient, err = cli.DOM.PseudoElementAdded(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.PseudoElementAddedClient)

	d.events.PseudoElementRemovedClient, err = cli.DOM.PseudoElementRemoved(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.PseudoElementRemovedClient)

	d.events.CharacterDataModifiedClient, err = cli.DOM.CharacterDataModified(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.CharacterDataModifiedClient)

	d.events.ChildNodeCountUpdatedClient, err = cli.DOM.ChildNodeCountUpdated(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.ChildNodeCountUpdatedClient)

	d.events.InlineStyleInvalidatedClient, err = cli.DOM.InlineStyleInvalidated(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.InlineStyleInvalidatedClient)

	d.events.DistributedNodesUpdatedClient, err = cli.DOM.DistributedNodesUpdated(ctx)
	if err != nil {
		return nil, err
	}
	d.events.clients = append(d.events.clients, d.events.DistributedNodesUpdatedClient)

	return d, nil
}

type DOM struct {
	cdp.DOM
	ErrTracer

	events struct {
		dom.SetChildNodesClient
		dom.DocumentUpdatedClient
		dom.AttributeRemovedClient
		dom.AttributeModifiedClient
		dom.ChildNodeRemovedClient
		dom.ShadowRootPoppedClient
		dom.ShadowRootPushedClient
		dom.ChildNodeInsertedClient
		dom.PseudoElementAddedClient
		dom.PseudoElementRemovedClient
		dom.CharacterDataModifiedClient
		dom.ChildNodeCountUpdatedClient
		dom.InlineStyleInvalidatedClient
		dom.DistributedNodesUpdatedClient
		clients []closer
	}

	mu            sync.RWMutex
	root          *Node
	nodes         map[dom.NodeID]*Node
	nodeModifiers map[dom.NodeID]nodeModifiers
}

func (d *DOM) loop(ctx context.Context) {
	defer d.close()

	for {
		select {
		case <-ctx.Done():
			return

		case <-d.events.SetChildNodesClient.Ready():
			if err := d.handleSetChildNodes(); err != nil {
				d.Trace(err)
			}

		case <-d.events.DocumentUpdatedClient.Ready():
			if err := d.handleDocumentUpdated(); err != nil {
				d.Trace(err)
			}

		case <-d.events.AttributeRemovedClient.Ready():
			if err := d.handleAttributeRemoved(); err != nil {
				d.Trace(err)
			}

		case <-d.events.AttributeModifiedClient.Ready():
			if err := d.handleAttributeModified(); err != nil {
				d.Trace(err)
			}

		case <-d.events.ChildNodeRemovedClient.Ready():
			if err := d.handleChildNodeRemoved(); err != nil {
				d.Trace(err)
			}

		case <-d.events.ShadowRootPoppedClient.Ready():
			if err := d.handleShadowRootPopped(); err != nil {
				d.Trace(err)
			}

		case <-d.events.ShadowRootPushedClient.Ready():
			if err := d.handleShadowRootPushed(); err != nil {
				d.Trace(err)
			}

		case <-d.events.ChildNodeInsertedClient.Ready():
			if err := d.handleChildNodeInserted(); err != nil {
				d.Trace(err)
			}

		case <-d.events.PseudoElementAddedClient.Ready():
			if err := d.handlePseudoElementAdded(); err != nil {
				d.Trace(err)
			}

		case <-d.events.PseudoElementRemovedClient.Ready():
			if err := d.handlePseudoElementRemoved(); err != nil {
				d.Trace(err)
			}

		case <-d.events.CharacterDataModifiedClient.Ready():
			if err := d.handleCharacterDataModified(); err != nil {
				d.Trace(err)
			}

		case <-d.events.ChildNodeCountUpdatedClient.Ready():
			if err := d.handleChildNodeCountUpdated(); err != nil {
				d.Trace(err)
			}

		case <-d.events.InlineStyleInvalidatedClient.Ready():
			if err := d.handleInlineStyleInvalidated(); err != nil {
				d.Trace(err)
			}

		case <-d.events.DistributedNodesUpdatedClient.Ready():
			if err := d.handleDistributedNodesUpdated(); err != nil {
				d.Trace(err)
			}

		}
	}
}

func (d *DOM) close() {
	for _, c := range d.events.clients {
		c.Close()
	}
}
