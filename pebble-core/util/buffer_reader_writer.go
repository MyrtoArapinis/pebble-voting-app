package util

import (
	"io"

	"github.com/giry-dev/pebble-voting-app/pebble-core/common"
)

type BufferReader struct {
	buf []byte
}

func NewBufferReader(buf []byte) *BufferReader {
	return &BufferReader{buf}
}

func (r *BufferReader) Len() int {
	return len(r.buf)
}

func (r *BufferReader) Read(p []byte) (n int, err error) {
	if len(r.buf) < len(p) {
		err = io.ErrShortBuffer
	}
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return
}

func (r *BufferReader) ReadBytes(n int) (p []byte, err error) {
	if len(r.buf) < n {
		return nil, io.ErrShortBuffer
	}
	p = r.buf[:n]
	r.buf = r.buf[n:]
	return
}

func (r *BufferReader) Read32() (p [32]byte, err error) {
	if copy(p[:], r.buf) != 32 {
		return p, io.ErrShortBuffer
	}
	r.buf = r.buf[32:]
	return p, nil
}

func (r *BufferReader) ReadRemaining() (p []byte) {
	p = r.buf
	r.buf = nil
	return
}

func (r *BufferReader) ReadVector() (p []byte, err error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	l := int(b)
	if l > 127 {
		b, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		l = ((l & 127) << 8) + int(b)
		if l <= 127 {
			return nil, common.NewParsingError("vector", "non canonical length encoding")
		}
	}
	if l == 0 {
		return nil, nil
	}
	return r.ReadBytes(l)
}

func (r *BufferReader) ReadByte() (byte, error) {
	if len(r.buf) == 0 {
		return 0, io.ErrShortBuffer
	}
	n := r.buf[0]
	r.buf = r.buf[1:]
	return n, nil
}

func (r *BufferReader) ReadUint16() (uint16, error) {
	if len(r.buf) < 2 {
		return 0, io.ErrShortBuffer
	}
	n := (uint16(r.buf[0]) << 8) + uint16(r.buf[1])
	r.buf = r.buf[2:]
	return n, nil
}

func (r *BufferReader) ReadUint32() (uint32, error) {
	if len(r.buf) < 4 {
		return 0, io.ErrShortBuffer
	}
	n := uint32(r.buf[0])
	n = (n << 8) + uint32(r.buf[1])
	n = (n << 8) + uint32(r.buf[2])
	n = (n << 8) + uint32(r.buf[3])
	r.buf = r.buf[4:]
	return n, nil
}

func (r *BufferReader) ReadUint64() (uint64, error) {
	if len(r.buf) < 8 {
		return 0, io.ErrShortBuffer
	}
	n := uint64(r.buf[0])
	n = (n << 8) + uint64(r.buf[1])
	n = (n << 8) + uint64(r.buf[2])
	n = (n << 8) + uint64(r.buf[3])
	n = (n << 8) + uint64(r.buf[4])
	n = (n << 8) + uint64(r.buf[5])
	n = (n << 8) + uint64(r.buf[6])
	n = (n << 8) + uint64(r.buf[7])
	r.buf = r.buf[8:]
	return n, nil
}

type BufferWriter struct {
	Buffer []byte
}

func (w *BufferWriter) Len() int {
	return len(w.Buffer)
}

func (w *BufferWriter) Write(p []byte) (int, error) {
	w.Buffer = append(w.Buffer, p...)
	return len(p), nil
}

func (w *BufferWriter) Write32(p [32]byte) {
	w.Buffer = append(w.Buffer, p[:]...)
}

func (w *BufferWriter) WriteAll(ps ...[]byte) {
	for _, p := range ps {
		w.Buffer = append(w.Buffer, p...)
	}
}

func (w *BufferWriter) WriteVector(p []byte) {
	l := len(p)
	if l > 127 {
		if l > 0x7FFF {
			panic("vector too big")
		}
		w.Buffer = append(w.Buffer, byte((l>>8)|128), byte(l))
	} else {
		w.Buffer = append(w.Buffer, byte(l))
	}
	w.Buffer = append(w.Buffer, p...)
}

func (w *BufferWriter) WriteByte(c byte) error {
	w.Buffer = append(w.Buffer, c)
	return nil
}

func (w *BufferWriter) WriteUint16(n uint16) {
	w.Buffer = append(w.Buffer, byte(n>>8), byte(n))
}

func (w *BufferWriter) WriteUint32(n uint32) {
	w.Buffer = append(w.Buffer, byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

func (w *BufferWriter) WriteUint64(n uint64) {
	w.Buffer = append(w.Buffer, byte(n>>56), byte(n>>48), byte(n>>40), byte(n>>32), byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}
