package chromy

import (
	"context"
	"time"
)

func Sleep(d time.Duration) Action {
	return actionWrapper("sleep", ActionFunc(func(ctx context.Context, t *Target) error {
		select {
		case <-time.After(d):
			return nil

		case <-ctx.Done():
			return ctx.Err()
		}
	}))
}
