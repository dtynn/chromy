package chromy

import (
	"github.com/mafredri/cdp/protocol/network"
)

func (n *Network) updateRequest(requestID network.RequestID, fn func(r *Request) bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	request, ok := n.requests[requestID]
	if !ok {
		request = &Request{
			RequestID: requestID,
		}

		n.requests[requestID] = request
	}

	if checkDone := fn(request); !checkDone {
		return
	}

	if !request.Done() {
		return
	}

	i := 0
	for {
		if i == len(n.watcher) {
			break
		}

		w := n.watcher[i]

		if w(request) {
			n.watcher = append(n.watcher[:i], n.watcher[i+1:]...)
			continue
		}

		i++
	}
}

func (n *Network) handleLoadingFailed() error {
	reply, err := n.events.LoadingFailedClient.Recv()
	if err != nil {
		return err
	}

	n.updateRequest(reply.RequestID, applyLoadingFailedReply(reply))
	return nil
}

func (n *Network) handleLoadingFinished() error {
	reply, err := n.events.LoadingFinishedClient.Recv()
	if err != nil {
		return err
	}

	n.updateRequest(reply.RequestID, applyLoadingFinishedReply(reply))
	return nil
}

func (n *Network) handleResponseReceived() error {
	reply, err := n.events.ResponseReceivedClient.Recv()
	if err != nil {
		return err
	}

	n.updateRequest(reply.RequestID, applyResponseReceivedReply(reply))
	return nil
}

func (n *Network) handleRequestWillBeSent() error {
	reply, err := n.events.RequestWillBeSentClient.Recv()
	if err != nil {
		return err
	}

	n.updateRequest(reply.RequestID, applyRequestWillBeSentReply(reply))
	return nil
}

func (n *Network) handleRequestServedFromCache() error {
	reply, err := n.events.RequestServedFromCacheClient.Recv()
	if err != nil {
		return err
	}

	n.updateRequest(reply.RequestID, applyRequestServedFromCache())
	return nil
}
