package chromy

import (
	"context"
	"net/http"
	"time"

	"github.com/mafredri/cdp/devtool"
)

// Connect return a connector to specified remote debugging url
func Connect(opt ...Option) *Connector {
	c := &Connector{
		remoteDebuggingURL: defaultRemoteDebuggingURL,
		timeout:            defaultConnectTimeout,
		actionTimeout:      defatulActionTimeout,
		taskStepTimeount:   defatulTaskStepTimeout,
	}

	for _, o := range opt {
		o(c)
	}

	c.devt = devtool.New(c.remoteDebuggingURL, devtool.WithClient(&http.Client{
		Timeout: c.timeout,
	}))

	return c
}

// Option connect option
type Option func(*Connector)

// RemoteDebuggingURL mo
func RemoteDebuggingURL(remoteURL string) Option {
	return func(c *Connector) {
		c.remoteDebuggingURL = remoteURL
	}
}

func ConnectTimeout(timeout time.Duration) Option {
	return func(c *Connector) {
		c.timeout = timeout
	}
}

func ActionTimeout(timeout time.Duration) Option {
	return func(c *Connector) {
		if timeout > 0 {
			c.actionTimeout = timeout
		}
	}
}

func TaskStepTimeout(timeout time.Duration) Option {
	return func(c *Connector) {
		if timeout > 0 {
			c.taskStepTimeount = timeout
		}
	}
}

type Connector struct {
	remoteDebuggingURL string
	timeout,
	actionTimeout time.Duration
	taskStepTimeount time.Duration

	devt *devtool.DevTools
}

func (c *Connector) New(ctx context.Context) (*Target, error) {
	return newTarget(ctx, c)
}
