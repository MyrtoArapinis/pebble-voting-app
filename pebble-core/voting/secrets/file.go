package secrets

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"os"

	"github.com/giry-dev/pebble-voting-app/pebble-core/anoncred"
	"github.com/giry-dev/pebble-voting-app/pebble-core/base32c"
	"github.com/giry-dev/pebble-voting-app/pebble-core/pubkey"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/vdf"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type FileSecretsManager string

type secrets struct {
	Seed      []byte              `json:"seed,omitempty"`
	Elections map[string]election `json:"elections"`
}

type election struct {
	Invitation string `json:"invitation,omitempty"`
	Title      string `json:"title,omitempty"`
	Ballot     []byte `json:"ballot,omitempty"`
	Solution   []byte `json:"solution,omitempty"`
}

type BasicElectionInfo struct {
	Id         string `json:"id"`
	Invitation string `json:"invitation"`
	Title      string `json:"title"`
}

func (f FileSecretsManager) GetElections() ([]BasicElectionInfo, error) {
	s, err := loadSecrets(string(f))
	if err != nil {
		return nil, err
	}
	res := make([]BasicElectionInfo, 0, len(s.Elections))
	for id, el := range s.Elections {
		res = append(res, BasicElectionInfo{
			Id:         id,
			Invitation: el.Invitation,
			Title:      el.Title,
		})
	}
	return res, nil
}

func (f FileSecretsManager) SetElection(info BasicElectionInfo) error {
	s, err := loadSecrets(string(f))
	if err != nil {
		return err
	}
	s.Elections[info.Id] = election{
		Invitation: info.Invitation,
		Title:      info.Title,
	}
	return s.store(string(f))
}

func (f FileSecretsManager) GetPrivateKey(_ *structs.EligibilityList) (pubkey.PrivateKey, error) {
	s, err := loadSecrets(string(f))
	if err != nil {
		return pubkey.PrivateKey{}, err
	}
	return pubkey.NewKeyFromSeed(util.KDF(s.Seed, "pubkey.ed25519")[:32]), nil
}

func (f FileSecretsManager) GetAnonymitySecret(eid [32]byte, sys anoncred.CredentialSystem) (anoncred.Secret, error) {
	s, err := loadSecrets(string(f))
	if err != nil {
		return nil, err
	}
	return sys.DeriveSecret(util.KDFid(s.Seed, eid, "anoncred.anoncred1")[:32])
}

func (f FileSecretsManager) GetBallot(eid [32]byte) (ballot structs.SignedBallot, err error) {
	s, err := loadSecrets(string(f))
	if err != nil {
		return
	}
	el, ok := s.Elections[base32c.Encode(eid[:])]
	if !ok || len(el.Ballot) == 0 {
		return ballot, errors.New("pebble: no recorded ballot")
	}
	err = ballot.FromBytes(el.Ballot)
	return
}

func (f FileSecretsManager) SetBallot(eid [32]byte, ballot structs.SignedBallot) error {
	s, err := loadSecrets(string(f))
	if err != nil {
		return err
	}
	id := base32c.Encode(eid[:])
	el := s.Elections[id]
	el.Ballot = ballot.Bytes()
	s.Elections[id] = el
	return s.store(string(f))
}

func (f FileSecretsManager) GetVdfSolution(eid [32]byte) (sol vdf.VdfSolution, err error) {
	s, err := loadSecrets(string(f))
	if err != nil {
		return
	}
	el, ok := s.Elections[base32c.Encode(eid[:])]
	if !ok || len(el.Solution) == 0 {
		return sol, errors.New("pebble: no recorded solution")
	}
	r := util.NewBufferReader(el.Solution)
	sol.Input, err = r.ReadVector()
	if err != nil {
		return
	}
	sol.Output, err = r.ReadVector()
	if err != nil {
		return
	}
	sol.Proof, err = r.ReadVector()
	return
}

func (f FileSecretsManager) SetVdfSolution(eid [32]byte, sol vdf.VdfSolution) error {
	s, err := loadSecrets(string(f))
	if err != nil {
		return err
	}
	var w util.BufferWriter
	w.WriteVector(sol.Input)
	w.WriteVector(sol.Output)
	w.WriteVector(sol.Proof)
	id := base32c.Encode(eid[:])
	el := s.Elections[id]
	el.Solution = w.Buffer
	s.Elections[id] = el
	return s.store(string(f))
}

func loadSecrets(fname string) (*secrets, error) {
	s := &secrets{Elections: make(map[string]election)}
	b, err := os.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			s.Seed = make([]byte, 32)
			_, err = rand.Read(s.Seed)
			if err != nil {
				return nil, err
			}
			err = s.store(fname)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
		return nil, err
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *secrets) store(fname string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(fname, b, 0600)
}
