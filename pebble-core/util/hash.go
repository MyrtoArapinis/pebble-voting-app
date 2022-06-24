package util

import "crypto/sha256"

type HashValue = [32]byte

func Hash(data []byte) HashValue {
	return sha256.Sum256(data)
}

func HashAll(data ...[]byte) (h HashValue) {
	f := sha256.New()
	for _, p := range data {
		f.Write(p)
	}
	f.Sum(h[:0])
	return
}
