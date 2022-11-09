package secrets

import (
	"github.com/giry-dev/pebble-voting-app/pebble-core/anoncred"
	"github.com/giry-dev/pebble-voting-app/pebble-core/pubkey"
	"github.com/giry-dev/pebble-voting-app/pebble-core/vdf"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type SecretsManager interface {
	GetPrivateKey(ell *structs.EligibilityList) (pubkey.PrivateKey, error)
	GetAnonymitySecret(eid [32]byte, sys anoncred.CredentialSystem) (anoncred.Secret, error)
	GetBallot(eid [32]byte) (structs.SignedBallot, error)
	SetBallot(eid [32]byte, ballot structs.SignedBallot) error
	GetVdfSolution(eid [32]byte) (vdf.VdfSolution, error)
	SetVdfSolution(eid [32]byte, sol vdf.VdfSolution) error
}
