package anoncred

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"sort"

	"github.com/consensys/gnark-crypto/ecc"
	bls381_mimc "github.com/consensys/gnark-crypto/ecc/bls12-381/fr/mimc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
)

const curveID = ecc.BLS12_381

var (
	errIncompatibleSecret     = errors.New("anoncred1: secret not compatible with system")
	errIncompatibleCommitment = errors.New("anoncred1: commitment not compatible with system")
)

func appendUint32(slice []byte, value uint32) []byte {
	return append(slice, byte(value), byte(value>>8), byte(value>>16), byte(value>>24))
}

func getUint32(slice []byte) uint32 {
	return uint32(slice[0]) | (uint32(slice[1]) << 8) | (uint32(slice[2]) << 16) | (uint32(slice[3]) << 24)
}

type AnonCred1 struct {
	depth int
	cs    frontend.CompiledConstraintSystem
	pk    groth16.ProvingKey
	vk    groth16.VerifyingKey
}

func (params *AnonCred1) ToBytes() (result []byte, err error) {
	var csBuf bytes.Buffer
	if _, err = params.cs.WriteTo(&csBuf); err != nil {
		return
	}
	var pkBuf bytes.Buffer
	if _, err = params.pk.WriteTo(&pkBuf); err != nil {
		return
	}
	var vkBuf bytes.Buffer
	if _, err = params.vk.WriteTo(&vkBuf); err != nil {
		return
	}
	result = appendUint32(result, uint32(params.depth))
	result = appendUint32(result, uint32(csBuf.Len()))
	result = appendUint32(result, uint32(pkBuf.Len()))
	result = appendUint32(result, uint32(vkBuf.Len()))
	result = append(result, csBuf.Bytes()...)
	result = append(result, pkBuf.Bytes()...)
	result = append(result, vkBuf.Bytes()...)
	return
}

func (params *AnonCred1) FromBytes(buffer []byte) (err error) {
	params.depth = int(getUint32(buffer))
	csStart := 16
	pkStart := csStart + int(getUint32(buffer[4:]))
	vkStart := pkStart + int(getUint32(buffer[8:]))
	vkEnd := vkStart + int(getUint32(buffer[12:]))
	params.cs = groth16.NewCS(curveID)
	if _, err = params.cs.ReadFrom(bytes.NewReader(buffer[csStart:pkStart])); err != nil {
		return
	}
	params.pk = groth16.NewProvingKey(curveID)
	if _, err = params.pk.ReadFrom(bytes.NewReader(buffer[pkStart:vkStart])); err != nil {
		return
	}
	params.vk = groth16.NewVerifyingKey(curveID)
	_, err = params.vk.ReadFrom(bytes.NewReader(buffer[vkStart:vkEnd]))
	return
}

type anonCred1Circuit struct {
	MessageHash frontend.Variable `gnark:",public"`
	SerialNo    frontend.Variable `gnark:",public"`
	Signature   frontend.Variable `gnark:",public"`
	MerkleRoot  frontend.Variable `gnark:",public"`
	Secret      frontend.Variable
	Directions  []frontend.Variable
	SideHashes  []frontend.Variable
}

type anonCred1Proof struct {
	MessageHash []byte
	SerialNo    []byte
	Signature   []byte
	MerkleRoot  []byte
	Proof       []byte
}

func hashMsg(data []byte) []byte {
	res := sha256.Sum256(data)
	return res[:]
}

func hashVars(h *mimc.MiMC, d1, d2 frontend.Variable) frontend.Variable {
	h.Reset()
	h.Write(d1, d2)
	return h.Sum()
}

func hashBytes(h hash.Hash, d1, d2 []byte) (b []byte, err error) {
	h.Reset()
	if _, err = h.Write(d1); err != nil {
		return
	}
	if _, err = h.Write(d2); err != nil {
		return
	}
	b = h.Sum(b)
	return
}

func (circuit *anonCred1Circuit) Define(curveID ecc.ID, api frontend.API) error {
	hFunc, err := mimc.NewMiMC("anoncred1", curveID, api)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.Signature, hashVars(&hFunc, circuit.MessageHash, circuit.Secret))
	sum := hashVars(&hFunc, circuit.SerialNo, circuit.Secret)
	for i, dir := range circuit.Directions {
		api.AssertIsBoolean(dir)
		d1 := api.Select(dir, sum, circuit.SideHashes[i])
		d2 := api.Select(dir, circuit.SideHashes[i], sum)
		sum = hashVars(&hFunc, d1, d2)
	}
	api.AssertIsEqual(sum, circuit.MerkleRoot)
	return nil
}

func (params *AnonCred1) SetupCircuit(depth int) error {
	var circuit anonCred1Circuit
	circuit.Directions = make([]frontend.Variable, depth)
	circuit.SideHashes = make([]frontend.Variable, depth)
	cs, err := frontend.Compile(curveID, backend.GROTH16, &circuit)
	if err != nil {
		return err
	}
	pk, vk, err := groth16.Setup(cs)
	params.depth = depth
	params.cs = cs
	params.pk = pk
	params.vk = vk
	return err
}

func (params *AnonCred1) prove(proof *anonCred1Proof, secret []byte, idx int, credentials [][]byte) (err error) {
	hFunc := bls381_mimc.NewMiMC("anoncred1")
	if proof.Signature, err = hashBytes(hFunc, proof.MessageHash, secret); err != nil {
		return
	}
	var witness anonCred1Circuit
	witness.MessageHash.Assign(proof.MessageHash)
	witness.SerialNo.Assign(proof.SerialNo)
	witness.Signature.Assign(proof.Signature)
	witness.Secret.Assign(secret)
	witness.Directions = nil
	witness.SideHashes = nil

	hashes := credentials
	for j := 0; j < params.depth; j++ {
		var newHashes [][]byte
		var hash []byte
		newIdx := 0
		for i := 0; i < len(hashes)-1; i += 2 {
			if i == idx || i+1 == idx {
				newIdx = len(newHashes)
				if i == idx {
					witness.Directions = append(witness.Directions, frontend.Value(1))
					witness.SideHashes = append(witness.SideHashes, frontend.Value(hashes[i+1]))
				} else {
					witness.Directions = append(witness.Directions, frontend.Value(0))
					witness.SideHashes = append(witness.SideHashes, frontend.Value(hashes[i]))
				}
			}
			if hash, err = hashBytes(hFunc, hashes[i], hashes[i+1]); err != nil {
				return
			}
			newHashes = append(newHashes, hash)
		}
		if len(hashes)%2 != 0 {
			hash = hashes[len(hashes)-1]
			if idx == len(hashes)-1 {
				newIdx = len(newHashes)
				witness.Directions = append(witness.Directions, frontend.Value(0))
				witness.SideHashes = append(witness.SideHashes, frontend.Value(hash))
			}
			if hash, err = hashBytes(hFunc, hash, hash); err != nil {
				return
			}
			newHashes = append(newHashes, hash)
		}
		hashes = newHashes
		idx = newIdx
	}

	proof.MerkleRoot = hashes[0]
	witness.MerkleRoot.Assign(proof.MerkleRoot)

	var groth16Proof groth16.Proof
	if groth16Proof, err = groth16.Prove(params.cs, params.pk, &witness); err != nil {
		return
	}
	var buffer bytes.Buffer
	if _, err = groth16Proof.WriteTo(&buffer); err != nil {
		return
	}
	proof.Proof = buffer.Bytes()
	return
}

func (params *AnonCred1) verify(proof anonCred1Proof) (err error) {
	groth16Proof := groth16.NewProof(curveID)
	if _, err = groth16Proof.ReadFrom(bytes.NewReader(proof.Proof)); err != nil {
		return
	}
	var circuit anonCred1Circuit
	circuit.MessageHash.Assign(proof.MessageHash)
	circuit.SerialNo.Assign(proof.SerialNo)
	circuit.Signature.Assign(proof.Signature)
	circuit.MerkleRoot.Assign(proof.MerkleRoot)
	circuit.Directions = make([]frontend.Variable, params.depth)
	circuit.SideHashes = make([]frontend.Variable, params.depth)
	return groth16.Verify(groth16Proof, params.vk, &circuit)
}

func hashMerkleTree(hashes [][]byte, depth int) (root []byte, err error) {
	hFunc := bls381_mimc.NewMiMC("anoncred1")
	for j := 0; j < depth; j++ {
		var newHashes [][]byte
		var hash []byte
		for i := 0; i < len(hashes)-1; i += 2 {
			if hash, err = hashBytes(hFunc, hashes[i], hashes[i+1]); err != nil {
				return
			}
			newHashes = append(newHashes, hash)
		}
		if len(hashes)%2 != 0 {
			hash = hashes[len(hashes)-1]
			if hash, err = hashBytes(hFunc, hash, hash); err != nil {
				return
			}
			newHashes = append(newHashes, hash)
		}
		hashes = newHashes
	}
	return hashes[0], nil
}

type anonCred1Secret struct {
	credential, secret []byte
}

func (s *anonCred1Secret) Commitment() (Commitment, error) {
	hFunc := bls381_mimc.NewMiMC("anoncred1")
	digest, err := hashBytes(hFunc, s.credential, s.secret)
	if err != nil {
		return nil, err
	}
	return &anonCred1Commitment{digest}, nil
}

func (s *anonCred1Secret) Credential() []byte {
	return s.credential
}

type anonCred1Commitment struct {
	bytes []byte
}

func (c *anonCred1Commitment) Bytes() []byte {
	return c.bytes
}

func (*AnonCred1) DeriveSecret(seed []byte) (Secret, error) {
	hFunc := bls381_mimc.NewMiMC("anoncred1.kdf")
	cred, err := hashBytes(hFunc, []byte{0}, seed)
	if err != nil {
		return nil, err
	}
	secret, err := hashBytes(hFunc, []byte{1}, seed)
	if err != nil {
		return nil, err
	}
	return &anonCred1Secret{cred, secret}, nil
}

func (*AnonCred1) ParseCommitment(p []byte) (Commitment, error) {
	if len(p) != 32 {
		return nil, errors.New("anoncred1: commitment must be 32 bytes")
	}
	return &anonCred1Commitment{p}, nil
}

type anonCred1Set struct {
	params *AnonCred1
	creds  [][]byte
	root   []byte
}

func (set *anonCred1Set) Len() int {
	return len(set.creds)
}

func (set *anonCred1Set) Less(i, j int) bool {
	a := set.creds[i]
	b := set.creds[j]
	for i = 0; i < 32; i++ {
		if a[i] >= b[i] {
			return false
		}
	}
	return true
}

func (set *anonCred1Set) Swap(i, j int) {
	t := set.creds[i]
	set.creds[i] = set.creds[j]
	set.creds[j] = t
}

func (params *AnonCred1) MakeAnonymitySet(commitments []Commitment) (AnonymitySet, error) {
	set := new(anonCred1Set)
	set.params = params
	if len(commitments) == 0 {
		return set, nil
	}
	for _, item := range commitments {
		com, ok := item.(*anonCred1Commitment)
		if !ok {
			return nil, errIncompatibleCommitment
		}
		set.creds = append(set.creds, com.bytes)
	}
	sort.Sort(set)
	creds := make([][]byte, 0, len(commitments))
	creds = append(creds, set.creds[0])
	for i := 1; i < len(set.creds); i++ {
		if !bytes.Equal(set.creds[i-1], set.creds[i]) {
			creds = append(creds, set.creds[i])
		}
	}
	set.creds = creds
	root, err := hashMerkleTree(creds, params.depth)
	if err != nil {
		return nil, err
	}
	set.root = root
	return set, nil
}

func (set *anonCred1Set) Sign(secret Secret, msg []byte) ([]byte, error) {
	sec, ok := secret.(*anonCred1Secret)
	if !ok {
		return nil, errIncompatibleSecret
	}
	pub, err := sec.Commitment()
	if err != nil {
		return nil, err
	}
	pubBytes := pub.Bytes()
	var proof anonCred1Proof
	proof.MessageHash = hashMsg(msg)
	proof.SerialNo = sec.credential
	proof.MerkleRoot = set.root
	idx := -1
	for i, b := range set.creds {
		if bytes.Equal(b, pubBytes) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, fmt.Errorf("credential not in credential set")
	}
	err = set.params.prove(&proof, sec.secret, idx, set.creds)
	if err != nil {
		return nil, err
	}
	return util.Concat(proof.Signature, proof.Proof), nil
}

func (set *anonCred1Set) Verify(serialNo, sig, msg []byte) error {
	if len(serialNo) != 32 {
		return fmt.Errorf("len(serialNo) != 32")
	}
	if len(sig) <= 32 {
		return fmt.Errorf("len(sig) <= 32")
	}
	var proof anonCred1Proof
	proof.MessageHash = hashMsg(msg)
	proof.SerialNo = serialNo
	proof.Signature = sig[:32]
	proof.MerkleRoot = set.root
	proof.Proof = sig[32:]
	return set.params.verify(proof)
}
