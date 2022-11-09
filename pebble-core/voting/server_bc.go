package voting

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/giry-dev/pebble-voting-app/pebble-core/base32c"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type BroadcastClient struct {
	client                 http.Client
	paramsURI, messagesURI string
	eid                    ElectionID
}

func NewBroadcastClient(id, uri string) (*BroadcastClient, error) {
	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}
	eid, err := base32c.Decode(id)
	if err != nil {
		return nil, err
	}
	if len(eid) != 32 {
		return nil, errors.New("pebble: election id not 32 bytes long")
	}
	bc := &BroadcastClient{
		paramsURI:   uri + "params/" + id,
		messagesURI: uri + "message/" + id,
	}
	copy(bc.eid[:], eid)
	return bc, nil
}

func (bc *BroadcastClient) Id() ElectionID {
	return bc.eid
}

func (bc *BroadcastClient) Params(_ context.Context) (*ElectionParams, error) {
	resp, err := bc.client.Get(bc.paramsURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	p := new(ElectionParams)
	err = p.FromBytes(buf)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (bc *BroadcastClient) Get(_ context.Context) ([]Message, error) {
	resp, err := bc.client.Get(bc.messagesURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := util.NewBufferReader(buf)
	var msgs []Message
	for r.Len() != 0 {
		kind, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		m, err := r.ReadVector()
		if err != nil {
			return nil, err
		}
		switch ElectionPhase(kind) {
		case CredGen:
			msg := new(structs.CredentialMessage)
			err = msg.FromBytes(m)
			if err == nil {
				msgs = append(msgs, Message{Credential: msg})
			}
		case Cast:
			msg := new(structs.SignedBallot)
			err = msg.FromBytes(m)
			if err == nil {
				msgs = append(msgs, Message{SignedBallot: msg})
			}
		case Tally:
			msg := new(structs.DecryptionMessage)
			err = msg.FromBytes(m)
			if err == nil {
				msgs = append(msgs, Message{Decryption: msg})
			}
		}
	}
	return msgs, nil
}

func (bc *BroadcastClient) Post(_ context.Context, m Message) error {
	resp, err := bc.client.Post(bc.messagesURI, "application/octet-stream", bytes.NewReader(m.Bytes()))
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
