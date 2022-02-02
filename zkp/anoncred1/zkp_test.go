package main

import "testing"

const depth = 8

func TestSetupCircuit(t *testing.T) {
	circuit, err := SetupCircuit(depth)
	if err != nil {
		t.Errorf("Error creating params: %s", err.Error())
	}
	bytes, err := circuit.ToBytes()
	if err != nil {
		t.Errorf("Error serializing circuit: %s", err.Error())
	}
	circuit = ZkpParams{}
	err = circuit.FromBytes(bytes)
	if err != nil {
		t.Errorf("Error deserializing circuit: %s", err.Error())
	}
}

func TestProveVerify(t *testing.T) {
	const credNum = 17
	const credPos = 5
	params, err := SetupCircuit(depth)
	if err != nil {
		t.Error("Error creating params")
	}
	var proof AnonCredProof
	proof.MessageHash, err = GenerateRandomScalar()
	if err != nil {
		t.Error("Error generating random message hash")
	}
	var secret []byte
	var credentials [][]byte
	for i := 0; i < credNum; i++ {
		cred, ser, sec, err := GenerateCredential()
		if err != nil {
			t.Error("Error generating credential")
		}
		credentials = append(credentials, cred)
		if i == credPos {
			secret = sec
			proof.SerialNo = ser
		}
	}
	err = params.Prove(&proof, secret, credPos, credentials)
	if err != nil {
		t.Errorf("Error generating proof: %s", err.Error())
	}
	proof.MerkleRoot, err = HashMerkleTree(credentials, params.Depth)
	if err != nil {
		t.Errorf("Error hashing Merkel tree: %s", err.Error())
	}
	err = params.Verify(proof)
	if err != nil {
		t.Errorf("Error verifying proof: %s", err.Error())
	}
	t.Logf("len(proof.Proof) = %d", len(proof.Proof))
}
