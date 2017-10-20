package chromy

import (
	"time"
)

const (
	// connection
	defaultRemoteDebuggingURL = "http://127.0.0.1:9222"
	defaultConnectTimeout     = 5 * time.Second
	defatulActionTimeout      = 1 * time.Minute
	defatulTaskStepTimeout    = 10 * time.Second
)

const (
	// action
	defaultWaitLoopInterval = 50 * time.Millisecond
)
