package main

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type electionData struct {
	params                            ElectionParams
	credentials, ballots, decryptions []Message
}

type mockServer struct {
	elections map[string]electionData
}

func NewMockServer() Server {
	return &mockServer{make(map[string]electionData)}
}

func (srv *mockServer) GetElection(id string) (*ElectionParams, *ServerError) {
	election, ok := srv.elections[id]
	if !ok {
		return nil, &ServerError{404, "Election not found"}
	}
	return &election.params, nil
}

func (srv *mockServer) CreateElection(params ElectionParams) (string, *ServerError) {
	var randBytes [32]byte
	_, err := rand.Read(randBytes[:])
	if err != nil {
		return "", &ServerError{500, err.Error()}
	}
	id := hex.EncodeToString(randBytes[:])
	srv.elections[id] = electionData{params, nil, nil, nil}
	return id, nil
}

func (srv *mockServer) GetMessages(id string, kind string) ([]Message, *ServerError) {
	election, ok := srv.elections[id]
	if !ok {
		return nil, &ServerError{404, "Election not found"}
	}
	switch kind {
	case "credential":
		return election.credentials, nil
	case "ballot":
		return election.ballots, nil
	case "decryption":
		return election.decryptions, nil
	default:
		return nil, &ServerError{500, "Message kind not supported"}
	}
}

func (srv *mockServer) PostMessage(id string, msg Message) *ServerError {
	election, ok := srv.elections[id]
	if !ok {
		return &ServerError{404, "Election not found"}
	}
	voteStart := election.params.VoteStart
	tallyStart := election.params.TallyStart
	now := time.Now()
	switch msg.Kind {
	case "credential":
		if !now.Before(voteStart) {
			return &ServerError{400, "Credential generation period ended"}
		}
		election.credentials = append(election.credentials, msg)
	case "ballot":
		if now.Before(voteStart) {
			return &ServerError{400, "Voting period not started"}
		}
		if !now.Before(tallyStart) {
			return &ServerError{400, "Voting period ended"}
		}
		election.ballots = append(election.ballots, msg)
	case "decryption":
		if now.Before(tallyStart) {
			return &ServerError{400, "Tallying period not started"}
		}
		election.decryptions = append(election.decryptions, msg)
	default:
		return &ServerError{500, "Message kind not supported"}
	}
	srv.elections[id] = election
	return nil
}
