package voting

import "github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"

type BroadcastChannel interface {
	PostCredential(cred structs.CredentialMessage) error
	PostSignedBallot(ballot structs.SignedBallot) error
	PostBallotDecryption(dec structs.DecryptionMessage) error
	GetCredentials() ([]structs.CredentialMessage, error)
	GetSignedBallots() ([]structs.SignedBallot, error)
	GetBallotDecryptions() ([]structs.DecryptionMessage, error)
}
