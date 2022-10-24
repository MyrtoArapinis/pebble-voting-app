package anoncred

import "testing"

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
	var secret SecretCredential
	var credentials []PublicCredential
	for i := 0; i < credNum; i++ {
		sec, err := params.GenerateSecretCredential()
		if err != nil {
			t.Error("Error generating credential")
		}
		pub, err := sec.Public()
		if err != nil {
			t.Error("Error generating credential")
		}
		credentials = append(credentials, pub)
		if i == credPos {
			secret = sec
		}
	}
	set, err := params.MakeCredentialSet(credentials)
	if err != nil {
		t.Error("Error making credential set")
	}
	sig, err := set.Sign(secret, msg)
	if err != nil {
		t.Errorf("Error signing: %s", err.Error())
	}
	err = set.Verify(secret.SerialNo(), sig, msg)
	if err != nil {
		t.Errorf("Error verifying proof: %s", err.Error())
	}
}
