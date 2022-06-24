package structs

import (
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/vdf"
)

type DecryptionMessage struct {
	InputHash     [32]byte
	Output, Proof []byte
}

func CreateDecryptionMessage(sol vdf.VdfSolution) DecryptionMessage {
	return DecryptionMessage{util.Hash(sol.Input), sol.Output, sol.Proof}
}

func (d *DecryptionMessage) Bytes() []byte {
	var w util.BufferWriter
	w.Write32(d.InputHash)
	w.WriteVector(d.Output)
	w.Write(d.Proof)
	return w.Buffer
}

func (d *DecryptionMessage) FromBytes(p []byte) error {
	r := util.NewBufferReader(p)
	var err error
	d.InputHash, err = r.Read32()
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
