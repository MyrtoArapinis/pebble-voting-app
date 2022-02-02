package main

import (
	"bytes"
	"crypto/rand"
	"hash"

	"github.com/consensys/gnark-crypto/ecc"
	bls381_mimc "github.com/consensys/gnark-crypto/ecc/bls12-381/fr/mimc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

const CURVE_ID = ecc.BLS12_381

type ZkpParams struct {
	Depth int
	CS    frontend.CompiledConstraintSystem
	PK    groth16.ProvingKey
	VK    groth16.VerifyingKey
}

func appendUint32(slice []byte, value uint32) []byte {
	return append(slice, byte(value), byte(value>>8), byte(value>>16), byte(value>>24))
}

func getUint32(slice []byte) uint32 {
	return uint32(slice[0]) | (uint32(slice[1]) << 8) | (uint32(slice[2]) << 16) | (uint32(slice[3]) << 24)
}

func (params *ZkpParams) ToBytes() (result []byte, err error) {
	var csBuf bytes.Buffer
	if _, err = params.CS.WriteTo(&csBuf); err != nil {
		return
	}
	var pkBuf bytes.Buffer
	if _, err = params.PK.WriteTo(&pkBuf); err != nil {
		return
	}
	var vkBuf bytes.Buffer
	if _, err = params.VK.WriteTo(&vkBuf); err != nil {
		return
	}
	result = appendUint32(result, uint32(params.Depth))
	result = appendUint32(result, uint32(csBuf.Len()))
	result = appendUint32(result, uint32(pkBuf.Len()))
	result = appendUint32(result, uint32(vkBuf.Len()))
	result = append(result, csBuf.Bytes()...)
	result = append(result, pkBuf.Bytes()...)
	result = append(result, vkBuf.Bytes()...)
	return
}

func (params *ZkpParams) FromBytes(buffer []byte) (err error) {
	params.Depth = int(getUint32(buffer))
	csStart := 16
	pkStart := csStart + int(getUint32(buffer[4:]))
	vkStart := pkStart + int(getUint32(buffer[8:]))
	vkEnd := vkStart + int(getUint32(buffer[12:]))
	params.CS = groth16.NewCS(CURVE_ID)
	if _, err = params.CS.ReadFrom(bytes.NewReader(buffer[csStart:pkStart])); err != nil {
		return
	}
	params.PK = groth16.NewProvingKey(CURVE_ID)
	if _, err = params.PK.ReadFrom(bytes.NewReader(buffer[pkStart:vkStart])); err != nil {
		return
	}
	params.VK = groth16.NewVerifyingKey(CURVE_ID)
	_, err = params.VK.ReadFrom(bytes.NewReader(buffer[vkStart:vkEnd]))
	return
}

type AnonCredCircuit struct {
	MessageHash frontend.Variable `gnark:",public"`
	SerialNo    frontend.Variable `gnark:",public"`
	Signature   frontend.Variable `gnark:",public"`
	MerkleRoot  frontend.Variable `gnark:",public"`
	Secret      frontend.Variable
	Directions  []frontend.Variable
	SideHashes  []frontend.Variable
}

type AnonCredProof struct {
	MessageHash []byte
	SerialNo    []byte
	Signature   []byte
	MerkleRoot  []byte
	Proof       []byte
}

func hashVars(h mimc.MiMC, d1, d2 frontend.Variable) frontend.Variable {
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

func (circuit *AnonCredCircuit) Define(curveID ecc.ID, api frontend.API) error {
	hFunc, err := mimc.NewMiMC("anoncred1", curveID, api)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.Signature, hashVars(hFunc, circuit.MessageHash, circuit.Secret))
	sum := hashVars(hFunc, circuit.SerialNo, circuit.Secret)
	for i, dir := range circuit.Directions {
		api.AssertIsBoolean(dir)
		d1 := api.Select(dir, sum, circuit.SideHashes[i])
		d2 := api.Select(dir, circuit.SideHashes[i], sum)
		sum = hashVars(hFunc, d1, d2)
	}
	api.AssertIsEqual(sum, circuit.MerkleRoot)
	return nil
}

func SetupCircuit(depth int) (params ZkpParams, _ error) {
	var circuit AnonCredCircuit
	circuit.Directions = make([]frontend.Variable, depth)
	circuit.SideHashes = make([]frontend.Variable, depth)
	cs, err := frontend.Compile(CURVE_ID, backend.GROTH16, &circuit)
	if err != nil {
		return params, err
	}
	pk, vk, err := groth16.Setup(cs)
	params = ZkpParams{depth, cs, pk, vk}
	return params, err
}

func (params *ZkpParams) Prove(proof *AnonCredProof, secret []byte, idx int, credentials [][]byte) (err error) {
	hFunc := bls381_mimc.NewMiMC("anoncred1")
	if proof.Signature, err = hashBytes(hFunc, proof.MessageHash, secret); err != nil {
		return
	}
	var witness AnonCredCircuit
	witness.MessageHash.Assign(proof.MessageHash)
	witness.SerialNo.Assign(proof.SerialNo)
	witness.Signature.Assign(proof.Signature)
	witness.Secret.Assign(secret)
	witness.Directions = nil
	witness.SideHashes = nil

	hashes := credentials
	for j := 0; j < params.Depth; j++ {
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
	if groth16Proof, err = groth16.Prove(params.CS, params.PK, &witness); err != nil {
		return
	}
	var buffer bytes.Buffer
	if _, err = groth16Proof.WriteTo(&buffer); err != nil {
		return
	}
	proof.Proof = buffer.Bytes()
	return
}

func (params *ZkpParams) Verify(proof AnonCredProof) (err error) {
	groth16Proof := groth16.NewProof(CURVE_ID)
	if _, err = groth16Proof.ReadFrom(bytes.NewReader(proof.Proof)); err != nil {
		return
	}
	var circuit AnonCredCircuit
	circuit.MessageHash.Assign(proof.MessageHash)
	circuit.SerialNo.Assign(proof.SerialNo)
	circuit.Signature.Assign(proof.Signature)
	circuit.MerkleRoot.Assign(proof.MerkleRoot)
	circuit.Directions = make([]frontend.Variable, params.Depth)
	circuit.SideHashes = make([]frontend.Variable, params.Depth)
	return groth16.Verify(groth16Proof, params.VK, &circuit)
}

func HashMerkleTree(credentials [][]byte, depth int) (root []byte, err error) {
	hFunc := bls381_mimc.NewMiMC("anoncred1")
	hashes := credentials
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

func GenerateRandomScalar() (cred []byte, err error) {
	cred = make([]byte, 32)
	if _, err = rand.Read(cred); err != nil {
		return nil, err
	}
	hFunc := bls381_mimc.NewMiMC("rand")
	if _, err = hFunc.Write(cred); err != nil {
		return nil, err
	}
	cred = hFunc.Sum(nil)
	return
}

func GenerateCredential() (cred, serialNo, secret []byte, err error) {
	if secret, err = GenerateRandomScalar(); err != nil {
		return nil, nil, nil, err
	}
	if serialNo, err = GenerateRandomScalar(); err != nil {
		return nil, nil, nil, err
	}
	hFunc := bls381_mimc.NewMiMC("anoncred1")
	if cred, err = hashBytes(hFunc, serialNo, secret); err != nil {
		return nil, nil, nil, err
	}
	return
}
