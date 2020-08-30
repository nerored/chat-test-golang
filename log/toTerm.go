/*
	将buffer内容输出到标准IO
*/
package log

import (
	"bytes"
	"io"
	"os"
)

type termWriter struct {
	infoBuffer *bytes.Buffer
	erroBuffer *bytes.Buffer
}

func newTermWriter() Writer {
	return &termWriter{
		infoBuffer: new(bytes.Buffer),
		erroBuffer: new(bytes.Buffer),
	}
}

func (t *termWriter) init() {

}

func (t *termWriter) exit() {

}

func (t *termWriter) needColor() bool {
	return true
}

func (t *termWriter) info() io.Writer {
	return t.infoBuffer
}

func (t *termWriter) erro() io.Writer {
	return t.erroBuffer
}

func (t *termWriter) reflush() {
	if t.infoBuffer != nil && t.infoBuffer.Len() > 0 {
		_, _ = io.Copy(os.Stdout, t.infoBuffer)
		t.infoBuffer.Reset()
	}

	if t.erroBuffer != nil && t.erroBuffer.Len() > 0 {
		_, _ = io.Copy(os.Stderr, t.erroBuffer)
		t.erroBuffer.Reset()
	}
}
