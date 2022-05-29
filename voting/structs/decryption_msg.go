package structs

import "github.com/giry-dev/pebble-voting-app/util"

type DecryptionMessage struct {
	InputHash, Output, Proof []byte
}

func (d *DecryptionMessage) Bytes() []byte {
	var w util.BufferWriter
	w.Write(d.InputHash)
	w.WriteVector(d.Output)
	w.Write(d.Proof)
	return w.Buffer
}

func (d *DecryptionMessage) FromBytes(p []byte) error {
	r := util.NewBufferReader(p)
	var err error
	d.InputHash, err = r.ReadBytes(32)
	if err != nil {
		return err
	}
	d.Output, err = r.ReadVector()
	if err != nil {
		return err
	}
	d.Proof = r.ReadRemaining()
	return nil
}
