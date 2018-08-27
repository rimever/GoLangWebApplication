package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちはtraceパッケージ")
		if buf.String() != "こんにちはtraceパッケージ\n" {
			t.Error("%という誤った文字列が出力されました", buf.String())
		}
	}
}