package chromy

import (
	"bytes"
	"context"
	"io"

	"github.com/mafredri/cdp/protocol/page"
)

func Navigate(URL string) Action {
	return actionWrapper("navigate", ActionFunc(func(ctx context.Context, t *Target) error {
		_, err := t.Client().Page.Navigate(ctx, page.NewNavigateArgs(URL))
		return err
	}))
}

func DocumentReady() Action {
	return actionWrapper("documentReady", ActionFunc(func(ctx context.Context, t *Target) error {
		c := t.domain.Page.events.DOMContentEventFiredClient
		select {
		case <-c.Ready():
			_, err := c.Recv()
			return err

		case <-ctx.Done():
			return ctx.Err()
		}
	}))
}

func PageLoaded() Action {
	return actionWrapper("pageLoaded", ActionFunc(func(ctx context.Context, t *Target) error {
		c := t.domain.Page.events.LoadEventFiredClient
		select {
		case <-c.Ready():
			_, err := c.Recv()
			return err

		case <-ctx.Done():
			return ctx.Err()
		}
	}))
}

func CaptureScreenshot(w io.Writer) Action {
	return actionWrapper("captureScreenshot", ActionFunc(func(ctx context.Context, t *Target) error {
		reply, err := t.Client().Page.CaptureScreenshot(ctx, page.NewCaptureScreenshotArgs())
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(reply.Data)
		_, err = io.Copy(w, buf)
		return err
	}))
}
