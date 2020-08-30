/*
	多源输出接口
*/
package log

import (
	"io"
	"sync"
)

type Writer interface {
	init()
	exit()
	needColor() bool
	reflush()
	info() io.Writer
	erro() io.Writer
}

type writerChan struct {
	writers []io.Writer
}

func (w *writerChan) isEmpty() bool {
	return len(w.writers) == 0
}

func (w *writerChan) free() {
	w.writers = w.writers[:0]
	writerChanPool.Put(w)
}

func (w *writerChan) Write(p []byte) (n int, err error) {
	for _, writer := range w.writers {
		if writer == nil {
			continue
		}

		return writer.Write(p)
	}

	return
}

var writerChanPool = sync.Pool{
	New: func() interface{} {
		return new(writerChan)
	},
}

func buildIOChanByWriters(level LogLv, writers []Writer) (newChan *writerChan) {
	newChan = writerChanPool.Get().(*writerChan)

	if cap(newChan.writers) < len(writers) {
		newChan.writers = make([]io.Writer, 0, len(writers))
	}

	switch level {
	case LOG_LEVEL_ERRO, LOG_LEVEL_FATA:
		for _, writer := range writers {
			if writer == nil {
				continue
			}
			newChan.writers = append(newChan.writers, writer.erro())
		}
	default:
		for _, writer := range writers {
			if writer == nil {
				continue
			}
			newChan.writers = append(newChan.writers, writer.info())
		}
	}

	return
}

func getMajorChan(level LogLv, colored, nocolor []Writer) (major repeater) {
	return repeater{
		colored: buildIOChanByWriters(level, colored),
		nocolor: buildIOChanByWriters(level, nocolor),
	}
}

func freeMajorChan(major *repeater) {
	if major == nil {
		return
	}

	if major.colored != nil {
		major.colored.free()
	}

	if major.nocolor != nil {
		major.nocolor.free()
	}
}
