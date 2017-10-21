package chromy

import (
	"fmt"
	"strings"

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

func actionErr(action string, err error) error {
	if err == nil {
		return nil
	}

	ae, ok := err.(*ActionError)
	if ok {
		ae.action = append(ae.action, "")
		copy(ae.action[1:], ae.action)
		ae.action[0] = action
		return ae
	}

	return &ActionError{
		action: []string{action},
		err:    err,
	}
}

type ActionError struct {
	action []string
	err    error
}

func (a *ActionError) Error() string {
	return fmt.Sprintf("[%s] %s", strings.Join(a.action, "; "), a.err)
}

func (a *ActionError) Cause() error {
	return a.err
}
