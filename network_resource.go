package chromy

import (
	"context"
	"fmt"
)

func NewResource(opt ...ResourceOption) *Resource {
	rw := &Resource{}
	for _, o := range opt {
		o(rw)
	}

	return rw
}

type ResourceOption func(*Resource)

type Resource struct {
	matchFn func(*Request) bool
	afterFn func(*Request) error
}

func (r *Resource) Do(ctx context.Context, t *Target) error {
	if r.matchFn == nil {
		return errNoMatchFunc
	}

	errCh := make(chan error, 1)

	stopper := func(req *Request) bool {
		if !r.matchFn(req) {
			return false
		}

		// inside the func to avoid send on closed channel
		defer close(errCh)

		if r.afterFn != nil {
			errCh <- r.afterFn(req)
		}

		return true
	}

	t.domain.Network.watch(stopper)

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-errCh:
		return err
	}
}

func ResourceMatch(fn func(*Request) bool) ResourceOption {
	return func(rw *Resource) {
		rw.matchFn = fn
	}
}

func ResourceAfter(fn func(*Request) error) ResourceOption {
	return func(rw *Resource) {
		rw.afterFn = fn
	}
}

func ResourcePattern(method string, urlpattern URLPattern) ResourceOption {
	return func(rw *Resource) {
		rw.matchFn = func(r *Request) bool {
			return r.Req.Request.Method == method && urlpattern.MatchString(r.Req.Request.URL)
		}
	}
}

func ResourceDone() ResourceOption {
	return ResourceAfter(
		func(req *Request) error {
			if req.IsFailed() {
				return fmt.Errorf("requset failed: %q", req.Failed.ErrorText)
			}

			if req.Resp.Response.Status/100 != 2 {
				return fmt.Errorf("response status %d: %s", req.Resp.Response.Status, req.Resp.Response.StatusText)
			}

			return nil
		},
	)
}
