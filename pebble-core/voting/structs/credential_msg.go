package structs

import (
	"github.com/giry-dev/pebble-voting-app/pebble-core/pubkey"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
)

type CredentialMessage struct {
	Credential []byte
	PublicKey  pubkey.PublicKey
	Signature  []byte
}

func (c *CredentialMessage) Bytes() []byte {
	var w util.BufferWriter
	w.WriteVector(c.Credential)
	w.WriteVector(c.PublicKey)
	w.Write(c.Signature)
	return w.Buffer
}

func (c *CredentialMessage) FromBytes(p []byte) error {
	r := util.NewBufferReader(p)
	var err error
	c.Credential, err = r.ReadVector()
	if err != nil {
		return err
	}
	c.PublicKey, err = r.ReadVector()
	if err != nil {
		return err
	}
	c.Signature = r.ReadRemaining()
	return nil
}

func (c *CredentialMessage) Sign(k pubkey.PrivateKey, eid util.HashValue) error {
	var err error
	c.PublicKey = k.Public()
	c.Signature, err = k.Sign(util.Concat(eid[:], c.Credential))
	return err
}

func (c *CredentialMessage) Verify(eid util.HashValue) error {
	return c.PublicKey.Verify(util.Concat(eid[:], c.Credential), c.Signature)
}
