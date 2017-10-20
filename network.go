package chromy

import (
	"context"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/network"
)

func newNetwork(ctx context.Context, cli *cdp.Client) (*Network, error) {
	err := cli.Network.Enable(ctx, network.NewEnableArgs())
	if err != nil {
		return nil, err
	}

	n := &Network{
		Network:   cli.Network,
		ErrTracer: nonErrTracer,

		requests: make(map[network.RequestID]*Request),
	}

	defer func() {
		if err != nil {
			n.close()
		}
	}()

	n.events.LoadingFailedClient, err = cli.Network.LoadingFailed(ctx)
	if err != nil {
		return nil, err
	}
	n.events.clients = append(n.events.clients, n.events.LoadingFailedClient)

	n.events.LoadingFinishedClient, err = cli.Network.LoadingFinished(ctx)
	if err != nil {
		return nil, err
	}
	n.events.clients = append(n.events.clients, n.events.LoadingFinishedClient)

	n.events.ResponseReceivedClient, err = cli.Network.ResponseReceived(ctx)
	if err != nil {
		return nil, err
	}
	n.events.clients = append(n.events.clients, n.events.ResponseReceivedClient)

	n.events.RequestWillBeSentClient, err = cli.Network.RequestWillBeSent(ctx)
	if err != nil {
		return nil, err
	}
	n.events.clients = append(n.events.clients, n.events.RequestWillBeSentClient)

	n.events.RequestServedFromCacheClient, err = cli.Network.RequestServedFromCache(ctx)
	if err != nil {
		return nil, err
	}
	n.events.clients = append(n.events.clients, n.events.RequestServedFromCacheClient)

	return n, nil
}

type Network struct {
	cdp.Network
	ErrTracer

	events struct {
		network.LoadingFailedClient
		network.LoadingFinishedClient
		network.ResponseReceivedClient
		network.RequestWillBeSentClient
		network.RequestServedFromCacheClient
		clients []closer
	}

	mu       sync.RWMutex
	requests map[network.RequestID]*Request
	watcher  []watcherStopper
}

func (n *Network) watch(w watcherStopper) {
	if w == nil {
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	requests := n.requests
	for _, req := range requests {
		if req.Done() && w(req) {
			return
		}
	}

	n.watcher = append(n.watcher, w)
	return

}

func (n *Network) loop(ctx context.Context) {
	defer n.close()

	for {
		select {
		case <-ctx.Done():
			return

		case <-n.events.LoadingFailedClient.Ready():
			if err := n.handleLoadingFailed(); err != nil {
				n.Trace(err)
			}

		case <-n.events.LoadingFinishedClient.Ready():
			if err := n.handleLoadingFinished(); err != nil {
				n.Trace(err)
			}

		case <-n.events.ResponseReceivedClient.Ready():
			if err := n.handleResponseReceived(); err != nil {
				n.Trace(err)
			}

		case <-n.events.RequestWillBeSentClient.Ready():
			if err := n.handleRequestWillBeSent(); err != nil {
				n.Trace(err)
			}

		case <-n.events.RequestServedFromCacheClient.Ready():
			if err := n.handleRequestServedFromCache(); err != nil {
				n.Trace(err)
			}

		}

	}
}

func (n *Network) close() {
	for _, c := range n.events.clients {
		c.Close()
	}
}
