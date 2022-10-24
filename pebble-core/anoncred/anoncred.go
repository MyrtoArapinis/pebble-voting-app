package anoncred

type SecretCredential interface {
	Bytes() []byte
	Public() (PublicCredential, error)
	SerialNo() []byte
}

type PublicCredential interface {
	Bytes() []byte
}

type CredentialSet interface {
	Len() int
	Sign(secret SecretCredential, msg []byte) ([]byte, error)
	Verify(serialNo, sig, msg []byte) error
}

type CredentialSystem interface {
	GenerateSecretCredential() (SecretCredential, error)
	ReadSecretCredential(p []byte) (SecretCredential, error)
	ReadPublicCredential(p []byte) (PublicCredential, error)
	MakeCredentialSet(credentials []PublicCredential) (CredentialSet, error)
}
