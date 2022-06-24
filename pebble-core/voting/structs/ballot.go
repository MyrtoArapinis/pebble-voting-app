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
	SerialNo        []byte
	Signature       []byte
}

func (b *SignedBallot) Bytes() []byte {
	var w util.BufferWriter
	w.WriteVector(b.SerialNo)
	w.WriteVector(b.Signature)
	w.Write(b.EncryptedBallot.Bytes())
	return w.Buffer
}

func (b *SignedBallot) FromBytes(p []byte) error {
	r := util.NewBufferReader(p)
	var err error
	b.SerialNo, err = r.ReadVector()
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
		return nil, errors.New("mismatched VDF solution")
	}
	cipher, err := createCipher(sol.Input)
	if err != nil {
		return nil, err
	}
	return cipher.Open(nil, eb.Payload[:12], eb.Payload[12:], nil)
}

func (eb *EncryptedBallot) Sign(set anoncred.CredentialSet, cred anoncred.SecretCredential) (sb SignedBallot, err error) {
	sb.EncryptedBallot = *eb
	sb.SerialNo = cred.SerialNo()
	sb.Signature, err = set.Sign(cred, eb.Bytes())
	return
}

func (b *SignedBallot) Verify(set anoncred.CredentialSet) error {
	return set.Verify(b.SerialNo, b.Signature, b.Bytes())
}
