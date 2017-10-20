package chromy

import (
	"context"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
)

func newTarget(ctx context.Context, c *Connector) (*Target, error) {
	dt, err := c.devt.Create(ctx)
	if err != nil {
		return nil, err
	}

	// TODO dial options
	conn, err := rpcc.DialContext(ctx, dt.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}

	t := &Target{
		c:    c,
		dt:   dt,
		conn: conn,
		cli:  cdp.NewClient(conn),
	}

	t.domain.Page, err = newPage(ctx, t.cli)
	if err != nil {
		return nil, err
	}

	t.domain.DOM, err = newDOM(ctx, t.cli)
	if err != nil {
		return nil, err
	}

	t.domain.Network, err = newNetwork(ctx, t.cli)
	if err != nil {
		return nil, err
	}

	t.wg.Add(3)

	loopCtx, loopCancel := context.WithCancel(context.Background())
	t.cancel = loopCancel

	go func() {
		defer t.wg.Done()

		t.domain.Page.loop(loopCtx)
	}()

	go func() {
		defer t.wg.Done()

		t.domain.DOM.loop(loopCtx)
	}()

	go func() {
		defer t.wg.Done()

		t.domain.Network.loop(loopCtx)
	}()

	return t, nil
}

// Target devtool target
type Target struct {
	c    *Connector
	dt   *devtool.Target
	conn *rpcc.Conn
	cli  *cdp.Client

	cancel context.CancelFunc

	wg sync.WaitGroup

	domain struct {
		*Page
		*DOM
		*Network
	}
}

func (t *Target) Run(ctx context.Context, action Action) error {
	ctx, cancel := context.WithTimeout(ctx, t.c.actionTimeout)
	defer cancel()

	return action.Do(ctx, t)
}

func (t *Target) Close() error {
	if t.cancel != nil {
		t.cancel()
	}

	t.conn.Close()
	return t.c.devt.Close(context.Background(), t.dt)
}

func (t *Target) Wait() {
	t.wg.Wait()
}

func (t *Target) Client() *cdp.Client {
	return t.cli
}
