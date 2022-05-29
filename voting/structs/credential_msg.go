package structs

import "github.com/giry-dev/pebble-voting-app/util"

type CredentialMessage struct {
	Credential, PublicKeyHash, Signature []byte
}

func (c *CredentialMessage) Bytes() []byte {
	var w util.BufferWriter
	w.WriteVector(c.Credential)
	w.Write(c.PublicKeyHash)
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
	c.PublicKeyHash, err = r.ReadBytes(32)
	if err != nil {
		return err
	}
	c.Signature = r.ReadRemaining()
	return nil
}
