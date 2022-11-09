package structs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"github.com/giry-dev/pebble-voting-app/pebble-core/anoncred"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/vdf"
)

type Ballot []byte

type EncryptedBallot struct {
	VdfInput, Payload []byte
}

var (
	ErrMismatchedVdfSolution = errors.New("pebble: mismatched VDF solution")
	ErrPayloadTooShort       = errors.New("pebble: ballot payload too short")
)

func (b *EncryptedBallot) Bytes() []byte {
	var w util.BufferWriter
	w.WriteVector(b.VdfInput)
	w.Write(b.Payload)
	return w.Buffer
}

func (b *EncryptedBallot) FromBytes(p []byte) error {
	r := util.NewBufferReader(p)
	var err error
	b.VdfInput, err = r.ReadVector()
	if err != nil {
		return err
	}
	b.Payload = r.ReadRemaining()
	return nil
}

type SignedBallot struct {
	EncryptedBallot EncryptedBallot
	Credential      []byte
	Signature       []byte
}

func (b *SignedBallot) Bytes() []byte {
	var w util.BufferWriter
	w.WriteVector(b.Credential)
	w.WriteVector(b.Signature)
	w.Write(b.EncryptedBallot.Bytes())
	return w.Buffer
}

func (b *SignedBallot) FromBytes(p []byte) error {
	r := util.NewBufferReader(p)
	var err error
	b.Credential, err = r.ReadVector()
	if err != nil {
		return err
	}
	b.Signature, err = r.ReadVector()
	if err != nil {
		return err
	}
	b.EncryptedBallot.FromBytes(r.ReadRemaining())
	return nil
}

func createCipher(vdfInput []byte) (cipher.AEAD, error) {
	key := sha256.Sum256(vdfInput)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func (b Ballot) Encrypt(sol vdf.VdfSolution) (eb EncryptedBallot, err error) {
	cipher, err := createCipher(sol.Input)
	if err != nil {
		return
	}
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}
	eb.Payload = util.Concat(nonce, cipher.Seal(nil, nonce, b, nil))
	eb.VdfInput = sol.Input
	return
}

func (eb *EncryptedBallot) Decrypt(sol vdf.VdfSolution) (Ballot, error) {
	if !bytes.Equal(sol.Input, eb.VdfInput) {
		return nil, ErrMismatchedVdfSolution
	}
	cipher, err := createCipher(sol.Input)
	if err != nil {
		return nil, err
	}
	if len(eb.Payload) < 12 {
		return nil, ErrPayloadTooShort
	}
	return cipher.Open(nil, eb.Payload[:12], eb.Payload[12:], nil)
}

func (eb *EncryptedBallot) Sign(set anoncred.AnonymitySet, secret anoncred.Secret) (sb SignedBallot, err error) {
	sb.EncryptedBallot = *eb
	sb.Credential = secret.Credential()
	sb.Signature, err = set.Sign(secret, eb.Bytes())
	return
}

func (b *SignedBallot) Verify(set anoncred.AnonymitySet) error {
	return set.Verify(b.Credential, b.Signature, b.Bytes())
}
