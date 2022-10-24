package methods

import (
	"errors"
	"sort"

	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

var ErrUnknownVotingMethod = errors.New("pebble: unknown voting method")

func Get(method string, numChoices int) (VotingMethod, error) {
	switch method {
	case "Approval":
		return &ApprovalVoting{numChoices}, nil
	case "Plurality":
		return &PluralityVoting{numChoices}, nil
	default:
		return nil, ErrUnknownVotingMethod
	}
}

type TallyCount struct {
	Index int
	Count uint64
}

type Tally []TallyCount

type VotingMethod interface {
	Vote(choices ...int) structs.Ballot
	Tally(ballots []structs.Ballot) Tally
}

func (t Tally) Sort() {
	sort.Sort(t)
}

func (t Tally) Len() int {
	return len(t)
}

func (t Tally) Less(i, j int) bool {
	if t[i].Count < t[j].Count {
		return true
	}
	return t[i].Index < t[j].Index
}

func (t Tally) Swap(i, j int) {
	tmp := t[i]
	t[i] = t[j]
	t[j] = tmp
}
