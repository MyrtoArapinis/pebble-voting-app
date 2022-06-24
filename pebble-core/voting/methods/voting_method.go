package methods

import (
	"sort"

	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

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
