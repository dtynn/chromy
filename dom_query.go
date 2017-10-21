package chromy

import (
	"context"
	"time"

	"github.com/mafredri/cdp/protocol/dom"
)

type QueryOption func(q *Query)

func QuerySelector(selector string) QueryOption {
	return func(q *Query) {
		q.queryFn = func(ctx context.Context, t *Target, root *Node) ([]dom.NodeID, error) {
			reply, err := t.Client().DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(root.NodeID, selector))
			if err != nil {
				if rerr, ok := rpccResponseError(err); ok && rerr.Code == -32000 {
					return nil, nil
				}

				return nil, err
			}

			return []dom.NodeID{reply.NodeID}, nil
		}
	}
}

func QuerySelectorAll(selector string) QueryOption {
	return func(q *Query) {
		q.queryFn = func(ctx context.Context, t *Target, root *Node) ([]dom.NodeID, error) {
			reply, err := t.Client().DOM.QuerySelectorAll(ctx, dom.NewQuerySelectorAllArgs(root.NodeID, selector))
			if err != nil {
				if rerr, ok := rpccResponseError(err); ok && rerr.Code == -32000 {
					return nil, nil
				}

				return nil, err
			}

			return reply.NodeIDs, nil
		}
	}
}

func QuerySearch(query string) QueryOption {
	return func(q *Query) {
		q.queryFn = func(ctx context.Context, t *Target, root *Node) ([]dom.NodeID, error) {
			reply, err := t.Client().DOM.PerformSearch(ctx, dom.NewPerformSearchArgs(query))
			if err != nil {
				return nil, err
			}

			res, err := t.Client().DOM.GetSearchResults(ctx, dom.NewGetSearchResultsArgs(reply.SearchID, 0, reply.ResultCount))
			if err != nil {
				return nil, err
			}

			return res.NodeIDs, nil
		}
	}
}

func QueryFunc(fn func(ctx context.Context, t *Target, root *Node) ([]dom.NodeID, error)) QueryOption {
	return func(q *Query) {
		q.queryFn = fn
	}
}

func QueryNodes(ptr *[]*Node) QueryOption {
	return func(q *Query) {
		q.afterFn = func(ctx context.Context, t *Target, nodes ...*Node) error {
			*ptr = nodes
			return nil
		}
	}
}

func QueryNode(ptr *Node) QueryOption {
	return func(q *Query) {
		q.afterFn = func(ctx context.Context, t *Target, nodes ...*Node) error {
			if len(nodes) == 0 {
				return ErrNodeNotFound
			}

			*ptr = *(nodes[0])
			return nil
		}
	}
}

func QueryAfter(fn func(ctx context.Context, t *Target, node *Node) error) QueryOption {
	return func(q *Query) {
		q.afterFn = func(ctx context.Context, t *Target, nodes ...*Node) error {
			if len(nodes) == 0 {
				return ErrNodeNotFound
			}

			return fn(ctx, t, nodes[0])
		}
	}
}

func QueryAfterAll(fn func(ctx context.Context, t *Target, nodes ...*Node) error) QueryOption {
	return func(q *Query) {
		q.afterFn = fn
	}
}

func NewQuery(opt ...QueryOption) *Query {
	q := &Query{}
	for _, o := range opt {
		o(q)
	}

	return q
}

type Query struct {
	queryFn func(ctx context.Context, t *Target, root *Node) ([]dom.NodeID, error)
	afterFn func(ctx context.Context, t *Target, nodes ...*Node) error
}

func (q *Query) Do(ctx context.Context, t *Target) error {
	return actionErr("query", q.do(ctx, t))
}

func (q *Query) do(ctx context.Context, t *Target) error {
	if q.queryFn == nil {
		return ErrNoQueryFunc
	}

	found := make([]*Node, 0, 20)

	dd := t.domain.DOM
	for {
	SELECT:
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			dd.mu.RLock()
			root := dd.root
			nodes := dd.nodes
			dd.mu.RUnlock()

			if root == nil {
				break
			}

			nodeIDs, err := q.queryFn(ctx, t, root)
			if err != nil {
				return err
			}

			if len(nodeIDs) == 0 {
				break
			}

			for _, id := range nodeIDs {
				n, ok := nodes[id]
				if !ok {
					found = found[:0]
					break SELECT
				}

				found = append(found, n)
			}

			if q.afterFn != nil {
				if err = q.afterFn(ctx, t, found...); err != nil {
					return err
				}

			}

			return nil
		}

		time.Sleep(defaultWaitLoopInterval)
	}
}
