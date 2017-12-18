package trace

import (
	"io"
	"fmt"
)
type tracer struct {
	out	io.Writer
}

type nilTracer struct {

}

type Tracer interface {
	Trace(...interface{})
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func (t *nilTracer) Trace(a...interface{}) {}

func Off() Tracer {
	return &nilTracer{}
}