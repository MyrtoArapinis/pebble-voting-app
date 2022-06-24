package secrets

import (
	"github.com/giry-dev/pebble-voting-app/pebble-core/anoncred"
	"github.com/giry-dev/pebble-voting-app/pebble-core/pubkey"
	"github.com/giry-dev/pebble-voting-app/pebble-core/vdf"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type SecretsManager interface {
	GetPrivateKey() (pubkey.PrivateKey, error)
	GetSecretCredential(sys anoncred.CredentialSystem) (anoncred.SecretCredential, error)
	GetBallot() (structs.SignedBallot, error)
	SetBallot(ballot structs.SignedBallot) error
	GetVdfSolution() (vdf.VdfSolution, error)
	SetVdfSolution(sol vdf.VdfSolution) error
}
