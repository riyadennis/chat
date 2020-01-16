package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Tracer should not return null")
	} else {
		Trace("Hello Trace Package")
		if buf.String() != "Hello Trace Package\n" {
			t.Errorf("Trace should not write %s", buf.String())
		}
	}
}
func TestOff(t *testing.T) {
	var slientTracer = Off()
	Trace("Nothing")
}
