package base32c

import "errors"

const alpha = "0123456789ABCDEFGHJKLMNPQRTVWXYZ"

var decodeMap map[rune]byte

var (
	ErrChar    = errors.New("pebble: invalid base32 character")
	ErrPadding = errors.New("pebble: non zero base32 padding")
)

func init() {
	decodeMap = make(map[rune]byte, 32)
	for i, r := range alpha {
		decodeMap[r] = byte(i)
	}
}

func Encode(p []byte) string {
	buf := make([]byte, 0, len(p)*8/5+1)
	u := uint(0)
	n := 0
	for _, b := range p {
		u |= uint(b) << n
		n += 8
		for ; n >= 5; n -= 5 {
			buf = append(buf, alpha[u&31])
			u >>= 5
		}
	}
	for ; n > 0; n -= 5 {
		buf = append(buf, alpha[u&31])
		u >>= 5
	}
	return string(buf)
}

func Decode(s string) ([]byte, error) {
	buf := make([]byte, 0, len(s)*5/8+1)
	u := uint(0)
	n := 0
	for _, c := range s {
		b, ok := decodeMap[c]
		if !ok {
			return nil, ErrChar
		}
		u |= uint(b) << n
		n += 5
		for ; n >= 8; n -= 8 {
			buf = append(buf, byte(u))
			u >>= 8
		}
	}
	if u != 0 {
		return nil, ErrPadding
	}
	return buf, nil
}
