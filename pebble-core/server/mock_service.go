package server

import (
	"context"
	"crypto/sha256"
	"errors"

	"github.com/giry-dev/pebble-voting-app/pebble-core/base32c"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting"
)

var (
	errNotFound = errors.New("pebble: election not found")
	errExists   = errors.New("pebble: election id already exists")
)

type mockService struct {
	elections map[string]*voting.Election
	ids       map[string]string
	url       string
}

func NewMockServer(url string, passHash []byte) *Server {
	if len(passHash) != 0 && len(passHash) != sha256.Size {
		panic("server: invalid password hash length")
	}
	return &Server{
		srv: &mockService{
			elections: make(map[string]*voting.Election),
			ids:       make(map[string]string),
			url:       url,
		},
		passHash: passHash,
		create:   true,
		post:     true,
	}
}

func (s *mockService) Create(spar ElectionSetupParams) error {
	if _, exists := s.ids[spar.AdminId]; exists {
		return errExists
	}
	epar, err := spar.Params()
	if err != nil {
		return err
	}
	id, err := util.RandomId()
	if err != nil {
		return err
	}
	bc := voting.NewMockBroadcastChannel(id, epar)
	election, err := voting.NewElection(context.Background(), bc, nil)
	if err != nil {
		return err
	}
	eid := base32c.Encode(id[:])
	s.elections[eid] = election
	s.ids[spar.AdminId] = eid
	return nil
}

func (s *mockService) Setup(adminId string) (info SetupInfo) {
	backendId, ok := s.ids[adminId]
	if !ok {
		info.Status = SetupError
		info.Error = "Election not found"
		return
	}
	if _, ok = s.elections[backendId]; !ok {
		info.Status = SetupError
		info.Error = "Election not found"
		return
	}
	var inv voting.Invitation
	inv.Network = "mock"
	inv.Address = []byte(backendId)
	inv.Servers = append(inv.Servers, s.url)
	info.Status = SetupDone
	info.BackendId = backendId
	info.Invitation = inv.String()
	return
}

func (s *mockService) Election(id string) (*voting.Election, error) {
	if el, ok := s.elections[id]; ok {
		return el, nil
	}
	return nil, errNotFound
}
