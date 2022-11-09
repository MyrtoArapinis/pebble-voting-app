package voting

import (
	"context"
	"errors"

	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/secrets"
)

var (
	ErrUnknownNetwork = errors.New("pebble: unknown network")
	ErrNoServers      = errors.New("pebble: no servers in invitation")
	ErrInvalidAddress = errors.New("pebble: invalid address")
)

func NewElectionFromInvitation(ctx context.Context, inv Invitation, sec secrets.SecretsManager) (*Election, error) {
	switch inv.Network {
	case "mock":
		if len(inv.Servers) == 0 {
			return nil, ErrNoServers
		}
		bc, err := NewBroadcastClient(string(inv.Address), inv.Servers[0])
		if err != nil {
			return nil, err
		}
		return NewElection(ctx, bc, sec)
	default:
		return nil, ErrUnknownNetwork
	}
}
