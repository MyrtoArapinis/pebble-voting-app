package anoncred

type Secret interface {
	Commitment() (Commitment, error)
	Credential() []byte
}

type Commitment interface {
	Bytes() []byte
}

type AnonymitySet interface {
	Len() int
	Sign(secret Secret, msg []byte) ([]byte, error)
	Verify(cred, sig, msg []byte) error
}

type CredentialSystem interface {
	DeriveSecret(seed []byte) (Secret, error)
	ParseCommitment(p []byte) (Commitment, error)
	MakeAnonymitySet(commitments []Commitment) (AnonymitySet, error)
}
