package chromy

import (
	"github.com/mafredri/cdp/protocol/network"
)

type watcherStopper func(req *Request) bool

type Request struct {
	RequestID network.RequestID
	Req       *network.RequestWillBeSentReply
	Resp      *network.ResponseReceivedReply
	Failed    *network.LoadingFailedReply
	Finished  *network.LoadingFinishedReply
	FromCache bool
}

func (r *Request) Done() bool {
	return r.IsFailed() || r.IsFinished() && r.Req != nil
}

func (r *Request) IsFailed() bool {
	return r.Failed != nil
}

func (r *Request) IsFinished() bool {
	return r.Resp != nil && r.Finished != nil
}

func applyRequestWillBeSentReply(reply *network.RequestWillBeSentReply) func(r *Request) bool {
	return func(r *Request) bool {
		r.Req = reply
		return true
	}
}

func applyResponseReceivedReply(reply *network.ResponseReceivedReply) func(r *Request) bool {
	return func(r *Request) bool {
		r.Resp = reply
		return true
	}
}

func applyLoadingFailedReply(reply *network.LoadingFailedReply) func(r *Request) bool {
	return func(r *Request) bool {
		r.Failed = reply
		return true
	}
}

func applyLoadingFinishedReply(reply *network.LoadingFinishedReply) func(r *Request) bool {
	return func(r *Request) bool {
		r.Finished = reply
		return true
	}
}

func applyRequestServedFromCache() func(r *Request) bool {
	return func(r *Request) bool {
		r.FromCache = true
		return false
	}
}
