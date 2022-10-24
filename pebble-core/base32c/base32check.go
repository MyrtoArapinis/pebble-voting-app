package base32c

import (
	"bytes"
	"crypto/sha256"
	"errors"
)

var (
	ErrLen   = errors.New("pebble: invalid base32 length")
	ErrCheck = errors.New("pebble: invalid base32 checksum")
)

func CheckEncode(p []byte) string {
	h := sha256.Sum256(p)
	b := make([]byte, 0, len(p)+4)
	b = append(b, p...)
	b = append(b, h[0], h[1], h[2], h[3])
	return Encode(b)
}

func CheckDecode(s string) ([]byte, error) {
	b, err := Decode(s)
	if err != nil {
		return nil, err
	}
	if len(b) < 4 {
		return nil, ErrLen
	}
	p := b[:len(b)-4]
	h := sha256.Sum256(p)
	if !bytes.Equal(h[:4], b[len(p):]) {
		return nil, ErrCheck
	}
	return p, nil
}
