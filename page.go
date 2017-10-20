package chromy

import (
	"context"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
)

func newPage(ctx context.Context, cli *cdp.Client) (*Page, error) {
	err := cli.Page.Enable(ctx)
	if err != nil {
		return nil, err
	}

	p := &Page{
		Page:      cli.Page,
		ErrTracer: nonErrTracer,
	}

	defer func() {
		if err != nil {
			p.close()
		}
	}()

	p.events.LoadEventFiredClient, err = cli.Page.LoadEventFired(ctx)
	if err != nil {
		return nil, err
	}
	p.events.clients = append(p.events.clients, p.events.LoadEventFiredClient)

	p.events.DOMContentEventFiredClient, err = cli.Page.DOMContentEventFired(ctx)
	if err != nil {
		return nil, err
	}
	p.events.clients = append(p.events.clients, p.events.DOMContentEventFiredClient)

	return p, nil
}

type Page struct {
	cdp.Page
	ErrTracer

	events struct {
		page.LoadEventFiredClient
		page.DOMContentEventFiredClient
		clients []closer
	}
}

func (p *Page) loop(ctx context.Context) {
	<-ctx.Done()
	p.close()
}

func (p *Page) close() {
	for _, c := range p.events.clients {
		c.Close()
	}
}
