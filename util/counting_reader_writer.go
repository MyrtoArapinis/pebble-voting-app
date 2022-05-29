package util

import "io"

type CountingWriter struct {
	w     io.Writer
	count int64
}

func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{w, 0}
}

func (w *CountingWriter) ResetCount() {
	w.count = 0
}

func (w *CountingWriter) Count() int64 {
	return w.count
}

func (w *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)
	w.count += int64(n)
	return
}

func (w *CountingWriter) WriteAll(ps ...[]byte) (n int, err error) {
	var t int
	for _, p := range ps {
		t, err = w.Write(p)
		n += t
		if err != nil {
			return
		}
	}
	return
}

type CountingReader struct {
	r     io.Reader
	count int64
}

func NewCountingReader(r io.Reader) *CountingReader {
	return &CountingReader{r, 0}
}

func (w *CountingReader) ResetCount() {
	w.count = 0
}

func (w *CountingReader) Count() int64 {
	return w.count
}

func (r *CountingReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.count += int64(n)
	return
}

func (r *CountingReader) ReadAll(ps ...[]byte) (n int, err error) {
	var t int
	for _, p := range ps {
		t, err = r.Read(p)
		n += t
		if err != nil {
			return
		}
	}
	return
}
