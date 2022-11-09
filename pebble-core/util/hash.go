package util

import (
	"crypto/sha256"
	"crypto/sha512"
)

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

func KDF(seed []byte, tag string) []byte {
	i := make([]byte, 0, 112)
	i = append(i, seed...)
	i = append(i, tag...)
	o := sha512.Sum512(i)
	return o[:]
}

func KDFid(seed []byte, id [32]byte, tag string) []byte {
	i := make([]byte, 0, 112)
	i = append(i, seed...)
	i = append(i, id[:]...)
	i = append(i, tag...)
	o := sha512.Sum512(i)
	return o[:]
}
