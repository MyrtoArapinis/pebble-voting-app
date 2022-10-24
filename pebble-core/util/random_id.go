package util

import "crypto/rand"

func RandomId() (id [32]byte, err error) {
	_, err = rand.Read(id[:])
	return
}
