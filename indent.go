package sr

import (
	"bytes"
	"io"
)

func IndentWriter(w io.Writer) io.Writer {
	return &newLinePrefixWriter{[]byte{'\t'}, w}
}

type newLinePrefixWriter struct {
	prefix []byte
	w      io.Writer
}

func (w *newLinePrefixWriter) Write(x []byte) (n int, err error) {
	for len(x) > 0 {
		k := bytes.Index(x, []byte{'\n'})
		if k < 0 {
			return w.w.Write(x)
		}
		m, err := w.w.Write(x[:k+1])
		n += m
		if err != nil {
			return n, err
		}
		_, err = w.w.Write([]byte(w.prefix))
		if err != nil {
			return n, err
		}
		x = x[k+1:]
	}
	return n, nil
}
