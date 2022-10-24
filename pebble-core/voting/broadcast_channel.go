package voting

import (
	"context"
	"errors"

	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type Message struct {
	ElectionParams *ElectionParams
	Credential     *structs.CredentialMessage
	SignedBallot   *structs.SignedBallot
	Decryption     *structs.DecryptionMessage
}

type BroadcastChannel interface {
	Id() ElectionID
	Params(ctx context.Context) (*ElectionParams, error)
	Get(ctx context.Context) ([]Message, error)
	Post(ctx context.Context, m Message) error
}

var (
	ErrInvalidMessageType = errors.New("pebble: invalid message type")
	ErrInvalidMessageSize = errors.New("pebble: invalid message size")
)

func (m Message) Bytes() []byte {
	var phase ElectionPhase
	var p []byte
	if m.ElectionParams != nil {
		phase = Setup
		p = m.ElectionParams.Bytes()
	} else if m.Credential != nil {
		phase = CredGen
		p = m.Credential.Bytes()
	} else if m.SignedBallot != nil {
		phase = Cast
		p = m.SignedBallot.Bytes()
	} else if m.Decryption != nil {
		phase = Tally
		p = m.Decryption.Bytes()
	} else {
		panic("pebble: invalid message type")
	}
	r := make([]byte, 1, len(p)+1)
	r[0] = byte(phase)
	r = append(r, p...)
	return r
}

func MessageFromBytes(p []byte) (m Message, err error) {
	if len(p) < 1 {
		return m, ErrInvalidMessageSize
	}
	switch ElectionPhase(p[0]) {
	case Setup:
		m.ElectionParams = new(ElectionParams)
		err = m.ElectionParams.FromBytes(p[1:])
	case CredGen:
		m.Credential = new(structs.CredentialMessage)
		err = m.Credential.FromBytes(p[1:])
	case Cast:
		m.SignedBallot = new(structs.SignedBallot)
		err = m.SignedBallot.FromBytes(p[1:])
	case Tally:
		m.Decryption = new(structs.DecryptionMessage)
		err = m.Decryption.FromBytes(p[1:])
	default:
		return m, ErrInvalidMessageType
	}
	return
}

type MockBroadcastChannel struct {
	messages []Message
	params   *ElectionParams
	id       ElectionID
}

func NewMockBroadcastChannel(id ElectionID, params *ElectionParams) *MockBroadcastChannel {
	return &MockBroadcastChannel{
		params: params,
		id:     id,
	}
}

func (bc *MockBroadcastChannel) Id() ElectionID {
	return bc.id
}

func (bc *MockBroadcastChannel) Params(ctx context.Context) (*ElectionParams, error) {
	return bc.params, nil
}

func (bc *MockBroadcastChannel) Get(ctx context.Context) ([]Message, error) {
	return bc.messages, nil
}

func (bc *MockBroadcastChannel) Post(ctx context.Context, m Message) error {
	bc.messages = append(bc.messages, m)
	return nil
}
