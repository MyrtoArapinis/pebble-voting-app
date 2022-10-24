package voting

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/giry-dev/pebble-voting-app/pebble-core/anoncred"
	"github.com/giry-dev/pebble-voting-app/pebble-core/pubkey"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/vdf"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/methods"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type mockSecretsManager struct {
	privateKey       pubkey.PrivateKey
	secretCredential anoncred.SecretCredential
	ballot           structs.SignedBallot
	solution         vdf.VdfSolution
}

func (sm *mockSecretsManager) GetPrivateKey() (pubkey.PrivateKey, error) {
	return sm.privateKey, nil
}

func (sm *mockSecretsManager) GetSecretCredential(sys anoncred.CredentialSystem) (anoncred.SecretCredential, error) {
	return sm.secretCredential, nil
}

func (sm *mockSecretsManager) GetBallot() (structs.SignedBallot, error) {
	return sm.ballot, nil
}

func (sm *mockSecretsManager) SetBallot(ballot structs.SignedBallot) error {
	sm.ballot = ballot
	return nil
}

func (sm *mockSecretsManager) GetVdfSolution() (vdf.VdfSolution, error) {
	return sm.solution, nil
}

func (sm *mockSecretsManager) SetVdfSolution(sol vdf.VdfSolution) error {
	sm.solution = sol
	return nil
}

func generateSecretCredentials(credSys anoncred.CredentialSystem, count int) (creds []anoncred.SecretCredential, err error) {
	creds = make([]anoncred.SecretCredential, count)
	for i := range creds {
		creds[i], err = credSys.GenerateSecretCredential()
		if err != nil {
			return nil, err
		}
	}
	return
}

func generatePrivateKeys(count int) (privs []pubkey.PrivateKey, err error) {
	privs = make([]pubkey.PrivateKey, count)
	for i := range privs {
		privs[i], err = pubkey.GenerateKey(pubkey.KeyTypeEd25519)
		if err != nil {
			return nil, err
		}
	}
	return
}

func generateEligibilityList(privs []pubkey.PrivateKey) (ell *structs.EligibilityList) {
	ell = structs.NewEligibilityList()
	for _, priv := range privs {
		ell.Add(util.Hash(priv.Public()), util.HashValue{})
	}
	return ell
}

func generateElectionParams(ell *structs.EligibilityList) (params ElectionParams) {
	now := time.Now()
	params.EligibilityList = ell
	params.CastStart = now.Add(time.Second * 20)
	params.TallyStart = now.Add(time.Second * 40)
	params.TallyEnd = now.Add(time.Second * 60)
	params.VotingMethod = "Plurality"
	params.Choices = []string{"Toby Wilkinson", "Ava McLean", "Oliver Rogers"}
	return
}

func TestElection(t *testing.T) {
	ctx := context.Background()
	credSys := new(anoncred.AnonCred1)
	err := credSys.SetupCircuit(8)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	privateKeys, err := generatePrivateKeys(10)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	elligibilityList := generateEligibilityList(privateKeys)
	electionParams := generateElectionParams(elligibilityList)
	secretsManager := new(mockSecretsManager)
	broadcast := new(MockBroadcastChannel)
	var election Election
	broadcast.params = &electionParams
	election.credSys = credSys
	election.channel = broadcast
	election.secrets = secretsManager
	election.vdf = &vdf.PietrzakVdf{MaxDifficulty: 1000000, DifficultyConversion: 10000}
	election.method = &methods.PluralityVoting{}
	election.params = &electionParams
	secretCredentials, err := generateSecretCredentials(credSys, len(privateKeys))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for i := range privateKeys {
		secretsManager.privateKey = privateKeys[i]
		secretsManager.secretCredential = secretCredentials[i]
		err = election.PostCredential(ctx)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}
	for time.Now().Before(electionParams.CastStart) {
		time.Sleep(time.Second)
	}
	voterIdx := rand.Intn(len(privateKeys))
	secretsManager.secretCredential = secretCredentials[voterIdx]
	err = election.Vote(ctx, rand.Intn(len(electionParams.Choices)))
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	for time.Now().Before(electionParams.TallyStart) {
		time.Sleep(time.Second)
	}
	err = election.RevealBallotDecryption(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	_, err = election.Progress(ctx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}
