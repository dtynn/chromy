package chromy

import (
	"io"
)

var nonErrTracer = &NonErrTracer{}

type closer interface {
	io.Closer
}

type ErrTracer interface {
	Trace(err error)
}

type NonErrTracer struct {
}

func (n *NonErrTracer) Trace(err error) {
	return
}
