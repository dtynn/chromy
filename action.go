package chromy

import (
	"context"
)

// Action
type Action interface {
	Do(context.Context, *Target) error
}

// ActionFunc
type ActionFunc func(context.Context, *Target) error

func (a ActionFunc) Do(ctx context.Context, t *Target) error {
	return a(ctx, t)
}

type Task []Action

func (t Task) Do(ctx context.Context, tar *Target) error {
	for _, a := range t {
		if err := step(ctx, tar, a); err != nil {
			return err
		}
	}

	return nil
}

func step(ctx context.Context, t *Target, a Action) error {
	ctx, cancel := context.WithTimeout(ctx, t.c.taskStepTimeount)
	defer cancel()

	return a.Do(ctx, t)
}
