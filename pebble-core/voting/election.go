package voting

import (
	"context"
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

	ErrDecryptionNotFound = errors.New("pebble: ballot decryption not found")
)

type ElectionID = [32]byte

type Election struct {
	credSys anoncred.CredentialSystem
	channel BroadcastChannel
	secrets secrets.SecretsManager
	vdf     vdf.VDF
	method  methods.VotingMethod
	params  *ElectionParams
}

type ElectionProgress struct {
	Phase        ElectionPhase
	Count, Total int
	Tally        methods.Tally
}

func NewElection(ctx context.Context, bc BroadcastChannel, sec secrets.SecretsManager) (*Election, error) {
	if anoncred.AnonCred1Instance == nil {
		return nil, errors.New("pebble: anoncred.AnonCred1Instance is nil")
	}
	params, err := bc.Params(ctx)
	if err != nil {
		return nil, err
	}
	method, err := methods.Get(params.VotingMethod, len(params.Choices))
	if err != nil {
		return nil, err
	}
	vdf := &vdf.PietrzakVdf{
		MaxDifficulty:        params.MaxVdfDifficulty,
		DifficultyConversion: uint64(float64(params.MaxVdfDifficulty) / params.TallyStart.Sub(params.CastStart).Seconds()),
	}
	return &Election{
		credSys: anoncred.AnonCred1Instance,
		channel: bc,
		secrets: sec,
		vdf:     vdf,
		method:  method,
		params:  params,
	}, nil
}

func (e *Election) Params() *ElectionParams {
	return e.params
}

func (e *Election) Phase() ElectionPhase {
	return e.params.Phase()
}

func (e *Election) Id() ElectionID {
	return e.channel.Id()
}

func (e *Election) Channel() BroadcastChannel {
	return e.channel
}

func (e *Election) PostCredentialCommitment(ctx context.Context) error {
	if e.params.Phase() != CredGen {
		return ErrWrongPhase
	}
	priv, err := e.secrets.GetPrivateKey(e.params.EligibilityList)
	if err != nil {
		return err
	}
	sec, err := e.secrets.GetAnonymitySecret(e.Id(), e.credSys)
	if err != nil {
		return err
	}
	com, err := sec.Commitment()
	if err != nil {
		return err
	}
	msg := new(structs.CredentialMessage)
	msg.Commitment = com.Bytes()
	err = msg.Sign(priv, e.Id())
	if err != nil {
		return err
	}
	return e.channel.Post(ctx, Message{Credential: msg})
}

func (e *Election) GetAnonymitySet(ctx context.Context) (anoncred.AnonymitySet, error) {
	if e.params.Phase() <= CredGen {
		return nil, ErrWrongPhase
	}
	msgs, err := e.channel.Get(ctx)
	if err != nil {
		return nil, err
	}
	creds := make(map[util.HashValue]anoncred.Commitment)
	for _, msg := range msgs {
		if msg.Credential == nil {
			continue
		}
		if msg.Credential.Verify(e.Id()) != nil {
			continue
		}
		cred, err := e.credSys.ParseCommitment(msg.Credential.Commitment)
		if err != nil {
			continue
		}
		creds[util.Hash(msg.Credential.PublicKey)] = cred
	}
	var list []anoncred.Commitment
	for _, c := range creds {
		list = append(list, c)
	}
	return e.credSys.MakeAnonymitySet(list)
}

func (e *Election) Vote(ctx context.Context, choices ...int) error {
	if e.params.Phase() != Cast {
		return ErrWrongPhase
	}
	set, err := e.GetAnonymitySet(ctx)
	if err != nil {
		return err
	}
	sol, err := e.vdf.Create(e.puzzleDuration())
	if err != nil {
		return err
	}
	err = e.secrets.SetVdfSolution(e.Id(), sol)
	if err != nil {
		return err
	}
	sec, err := e.secrets.GetAnonymitySecret(e.Id(), e.credSys)
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
	err = e.secrets.SetBallot(e.Id(), signBallot)
	if err != nil {
		return err
	}
	return e.channel.Post(ctx, Message{SignedBallot: &signBallot})
}

func (e *Election) puzzleDuration() uint64 {
	return uint64(e.params.TallyStart.Sub(e.params.CastStart).Seconds())
}

func (e *Election) RevealBallotDecryption(ctx context.Context) error {
	sol, err := e.secrets.GetVdfSolution(e.Id())
	if err != nil {
		return err
	}
	return e.PostBallotDecryption(ctx, sol)
}

func (e *Election) PostBallotDecryption(ctx context.Context, sol vdf.VdfSolution) error {
	if e.params.Phase() != Tally {
		return ErrWrongPhase
	}
	msg := structs.CreateDecryptionMessage(sol)
	return e.channel.Post(ctx, Message{Decryption: &msg})
}

func (e *Election) Progress(ctx context.Context) (p ElectionProgress, err error) {
	p.Phase = e.params.Phase()
	if p.Phase <= CredGen {
		return
	}
	set, err := e.GetAnonymitySet(ctx)
	if err != nil {
		return
	}
	msgs, err := e.channel.Get(ctx)
	if err != nil {
		return
	}
	var signBallots []structs.SignedBallot
	var decMsgs []structs.DecryptionMessage
	for _, msg := range msgs {
		if msg.SignedBallot != nil {
			signBallots = append(signBallots, *msg.SignedBallot)
		} else if msg.Decryption != nil {
			decMsgs = append(decMsgs, *msg.Decryption)
		}
	}
	var usedCredentials util.BytesSet
	var decBallots []structs.Ballot
	validSignBallots := 0
	validDecBallots := 0
	invalidDecBallots := 0
	for _, signBallot := range signBallots {
		if usedCredentials.Contains(signBallot.Credential) {
			continue
		}
		err = signBallot.Verify(set)
		if err != nil {
			continue
		}
		validSignBallots++
		if p.Phase >= Tally {
			ballot, err := decryptBallot(signBallot.EncryptedBallot, decMsgs, e.vdf)
			if err != nil {
				if err != ErrDecryptionNotFound {
					invalidDecBallots++
				}
				continue
			}
			decBallots = append(decBallots, ballot)
			validDecBallots++
		}
	}
	if p.Phase == Cast {
		p.Total = set.Len()
		p.Count = validSignBallots
	} else if p.Phase == Tally {
		p.Total = validSignBallots - invalidDecBallots
		p.Count = validDecBallots
		p.Tally = e.method.Tally(decBallots)
	} else {
		p.Total = validSignBallots
		p.Count = validDecBallots
		p.Tally = e.method.Tally(decBallots)
	}
	return p, nil
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
				return nil, err
			}
			return ballot, nil
		}
	}
	return nil, ErrDecryptionNotFound
}
