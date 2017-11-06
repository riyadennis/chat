package trace

import (
	"io"
	"fmt"
)

type Tracer interface {
	Trace(...interface{})
}
type nilTracer struct {}

type trace struct {
	out io.Writer
}
func (t trace) Trace(a ...interface{}){
	fmt.Fprintln(t.out, a...)
}
func (st nilTracer) Trace(a ...interface{}){}

func Off() Tracer{
	return &nilTracer{}
}

func New(w io.Writer) Tracer {
	return trace{
		out:w,
	}
}