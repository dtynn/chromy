package chromy

import (
	"context"

	"github.com/mafredri/cdp/protocol/page"
)

func Navigate(URL string) Action {
	return ActionFunc(func(ctx context.Context, t *Target) error {
		_, err := t.Client().Page.Navigate(ctx, page.NewNavigateArgs(URL))
		return err
	})
}

func DocumentReady() Action {
	return ActionFunc(func(ctx context.Context, t *Target) error {
		c := t.domain.Page.events.DOMContentEventFiredClient
		select {
		case <-c.Ready():
			_, err := c.Recv()
			return err

		case <-ctx.Done():
			return ctx.Err()
		}
	})
}

func PageLoaded() Action {
	return ActionFunc(func(ctx context.Context, t *Target) error {
		c := t.domain.Page.events.LoadEventFiredClient
		select {
		case <-c.Ready():
			_, err := c.Recv()
			return err

		case <-ctx.Done():
			return ctx.Err()
		}
	})
}
