package main

import "C"

//export FfiGenerateCredential
func FfiGenerateCredential(buffer []byte) int32 {
	cred, serialNo, secret, err := GenerateCredential()
	if err != nil {
		return -1
	}
	copy(buffer, cred)
	copy(buffer[32:], serialNo)
	copy(buffer[64:], secret)
	return 0
}

//export FfiHashMerkleTree
func FfiHashMerkleTree(destRoot, credentialsConcat []byte, depth int) int32 {
	credentials := make([][]byte, len(credentialsConcat)/32)
	for i := 0; i < len(credentials); i++ {
		credentials[i] = credentialsConcat[i*32 : (i+1)*32]
	}
	root, err := HashMerkleTree(credentials, depth)
	if err != nil {
		return -1
	}
	copy(destRoot, root)
	return 0
}

//export FfiProve
func FfiProve(out, paramsBytes, messageHash, serialNo, secret []byte,
	idx int, credentialsConcat []byte) int32 {
	var params ZkpParams
	if err := params.FromBytes(paramsBytes); err != nil {
		return -1
	}
	credentials := make([][]byte, len(credentialsConcat)/32)
	for i := 0; i < len(credentials); i++ {
		credentials[i] = credentialsConcat[i*32 : (i+1)*32]
	}
	var proof AnonCredProof
	proof.MessageHash = messageHash
	proof.SerialNo = serialNo
	if err := params.Prove(&proof, secret, idx, credentials); err != nil {
		return -1
	}
	copy(out[:32], proof.Signature)
	copy(out[32:], proof.Proof)
	return int32(len(proof.Signature) + len(proof.Proof))
}

//export FfiVerify
func FfiVerify(paramsBytes, messageHash, serialNo, signature, merkleRoot []byte) int32 {
	var params ZkpParams
	err := params.FromBytes(paramsBytes)
	if err != nil {
		return -1
	}
	err = params.Verify(AnonCredProof{messageHash, serialNo, signature[:32], merkleRoot, signature[32:]})
	if err != nil {
		return -1
	}
	return 0
}
