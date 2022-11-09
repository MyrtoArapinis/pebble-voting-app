package anoncred

import (
	"crypto/rand"
	"os"
	"testing"
)

const depth = 8

func TestSetupCircuit(t *testing.T) {
	var params AnonCred1
	err := params.SetupCircuit(depth)
	if err != nil {
		t.Errorf("Error creating params: %s", err.Error())
		return
	}
	bytes, err := params.ToBytes()
	if err != nil {
		t.Errorf("Error serializing circuit: %s", err.Error())
		return
	}
	err = params.FromBytes(bytes)
	if err != nil {
		t.Errorf("Error deserializing circuit: %s", err.Error())
		return
	}
	file, err := os.OpenFile("anoncred1-params.bin", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Log("Failed opening anoncred1-params.bin for writing")
		return
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		t.Log(err.Error())
	}
}

func TestProveVerify(t *testing.T) {
	const credNum = 17
	const credPos = 5
	var params AnonCred1
	err := params.SetupCircuit(depth)
	if err != nil {
		t.Error("Error creating params")
	}
	msg := make([]byte, 5)
	var secret Secret
	var commitments []Commitment
	for i := 0; i < credNum; i++ {
		var seed [32]byte
		rand.Reader.Read(seed[:])
		sec, err := params.DeriveSecret(seed[:])
		if err != nil {
			t.Error("Error generating credential")
		}
		com, err := sec.Commitment()
		if err != nil {
			t.Error("Error generating credential")
		}
		commitments = append(commitments, com)
		if i == credPos {
			secret = sec
		}
	}
	set, err := params.MakeAnonymitySet(commitments)
	if err != nil {
		t.Error("Error making anonymity set")
	}
	sig, err := set.Sign(secret, msg)
	if err != nil {
		t.Errorf("Error signing: %s", err.Error())
	}
	err = set.Verify(secret.Credential(), sig, msg)
	if err != nil {
		t.Errorf("Error verifying proof: %s", err.Error())
	}
}
