package voting

import (
	"errors"

	"github.com/giry-dev/pebble-voting-app/pebble-core/anoncred"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/vdf"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/methods"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/secrets"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

var (
	ErrWrongPhase = errors.New("pebble: wrong election phase")

	ErrBallotNotDecrypted = errors.New("pebble: ballot not decrypted")
)

type Election struct {
	credSys anoncred.CredentialSystem
	channel BroadcastChannel
	secrets secrets.SecretsManager
	vdf     vdf.VDF
	method  methods.VotingMethod
	params  *ElectionParams
}

func (e *Election) PostCredential() error {
	if e.params.Phase() != CredGen {
		return ErrWrongPhase
	}
	priv, err := e.secrets.GetPrivateKey()
	if err != nil {
		return err
	}
	sec, err := e.secrets.GetSecretCredential(e.credSys)
	if err != nil {
		return err
	}
	pub, err := sec.Public()
	if err != nil {
		return err
	}
	var msg structs.CredentialMessage
	msg.Credential = pub.Bytes()
	err = msg.Sign(priv, e.params.Id)
	if err != nil {
		return err
	}
	return e.channel.PostCredential(msg)
}

func (e *Election) GetCredentialSet() (anoncred.CredentialSet, error) {
	if e.params.Phase() <= CredGen {
		return nil, ErrWrongPhase
	}
	msgs, err := e.channel.GetCredentials()
	if err != nil {
		return nil, err
	}
	creds := make(map[util.HashValue]anoncred.PublicCredential)
	for _, msg := range msgs {
		if msg.Verify(e.params.Id) != nil {
			continue
		}
		cred, err := e.credSys.ReadPublicCredential(msg.Credential)
		if err != nil {
			continue
		}
		creds[util.Hash(msg.PublicKey)] = cred
	}
	var list []anoncred.PublicCredential
	for _, c := range creds {
		list = append(list, c)
	}
	return e.credSys.MakeCredentialSet(list)
}

func (e *Election) Vote(choices ...int) error {
	if e.params.Phase() != Cast {
		return ErrWrongPhase
	}
	set, err := e.GetCredentialSet()
	if err != nil {
		return err
	}
	sol, err := e.vdf.Create(e.puzzleDuration())
	if err != nil {
		return err
	}
	err = e.secrets.SetVdfSolution(sol)
	if err != nil {
		return err
	}
	sec, err := e.secrets.GetSecretCredential(e.credSys)
	if err != nil {
		return err
	}
	ballot := e.method.Vote(choices...)
	encBallot, err := ballot.Encrypt(sol)
	if err != nil {
		return err
	}
	signBallot, err := encBallot.Sign(set, sec)
	if err != nil {
		return err
	}
	err = e.secrets.SetBallot(signBallot)
	if err != nil {
		return err
	}
	return e.channel.PostSignedBallot(signBallot)
}

func (e *Election) puzzleDuration() uint64 {
	return uint64(e.params.TallyStart.Sub(e.params.CastStart).Seconds())
}

func (e *Election) RevealBallotDecryption() error {
	sol, err := e.secrets.GetVdfSolution()
	if err != nil {
		return err
	}
	return e.PostBallotDecryption(sol)
}

func (e *Election) PostBallotDecryption(sol vdf.VdfSolution) error {
	if e.params.Phase() != Tally {
		return ErrWrongPhase
	}
	msg := structs.CreateDecryptionMessage(sol)
	return e.channel.PostBallotDecryption(msg)
}

func (e *Election) Tally() (methods.Tally, error) {
	if e.params.Phase() != Tally {
		return nil, ErrWrongPhase
	}
	set, err := e.GetCredentialSet()
	if err != nil {
		return nil, err
	}
	signBallots, err := e.channel.GetSignedBallots()
	if err != nil {
		return nil, err
	}
	decMsgs, err := e.channel.GetBallotDecryptions()
	if err != nil {
		return nil, err
	}
	var serialNos util.BytesSet
	var decBallots []structs.Ballot
	for _, signBallot := range signBallots {
		if serialNos.Contains(signBallot.SerialNo) {
			continue
		}
		err = signBallot.Verify(set)
		if err != nil {
			continue
		}
		ballot, err := decryptBallot(signBallot.EncryptedBallot, decMsgs, e.vdf)
		if err != nil {
			continue
		}
		decBallots = append(decBallots, ballot)
	}
	return e.method.Tally(decBallots), nil
}

func decryptBallot(encBallot structs.EncryptedBallot, msgs []structs.DecryptionMessage, ivdf vdf.VDF) (structs.Ballot, error) {
	vdfInputHash := util.Hash(encBallot.VdfInput)
	for _, msg := range msgs {
		if msg.InputHash == vdfInputHash {
			sol := vdf.VdfSolution{Input: encBallot.VdfInput, Output: msg.Output, Proof: msg.Proof}
			err := ivdf.Verify(sol)
			if err != nil {
				continue
			}
			ballot, err := encBallot.Decrypt(sol)
			if err != nil {
				continue
			}
			return ballot, nil
		}
	}
	return nil, ErrBallotNotDecrypted
}
