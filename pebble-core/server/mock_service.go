package server

import (
	"context"
	"errors"

	"github.com/giry-dev/pebble-voting-app/pebble-core/base32c"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting"
)

var errNotFound = errors.New("pebble: election not found")

type MockService struct {
	elections map[string]*voting.Election
	ids       map[string]string
	url       string
}

func NewMockService() *MockService {
	s := new(MockService)
	s.elections = make(map[string]*voting.Election)
	s.ids = make(map[string]string)
	return s
}

func (s *MockService) Create(spar ElectionSetupParams) error {
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

func (s *MockService) Setup(adminId string) (info SetupInfo) {
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

func (s *MockService) Election(id string) (*voting.Election, error) {
	if el, ok := s.elections[id]; ok {
		return el, nil
	}
	return nil, errNotFound
}
