package chromy

import (
	"fmt"

	"github.com/mafredri/cdp/rpcc"
)

var (
	ErrNodeNotFound        = fmt.Errorf("dom node not found")
	ErrNoQueryFunc         = fmt.Errorf("no query function")
	ErrNoMatchFunc         = fmt.Errorf("no match function")
	ErrUnableToResolveNode = fmt.Errorf("unable to resolve node")
)

type causer interface {
	Cause() error
}

func errCause(err error) error {
	c, ok := err.(causer)
	if !ok {
		return err
	}

	return c.Cause()
}

func rpccResponseError(err error) (*rpcc.ResponseError, bool) {
	respErr, ok := errCause(err).(*rpcc.ResponseError)
	return respErr, ok
}

func nonblockErrorPush(ch chan error, err error) {
	select {
	case ch <- err:

	default:

	}
}
