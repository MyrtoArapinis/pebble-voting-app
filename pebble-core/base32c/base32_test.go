package base32c

import (
	"bytes"
	"testing"
)

func TestBase32(t *testing.T) {
	// b, _ := Decode("EPK0")
	// s := hex.EncodeToString(b)
	// t.Log(s)
	a := make([]byte, 256)
	for i := range a {
		a[i] = byte(i)
	}
	for i := range a {
		s := Encode(a[i:])
		b, err := Decode(s)
		if err != nil {
			t.Error(err)
			return
		}
		if !bytes.Equal(a[i:], b) {
			t.Error("decoded base32 payload different from original")
			return
		}
	}
}
