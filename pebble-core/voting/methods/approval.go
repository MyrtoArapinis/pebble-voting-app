package methods

import "github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"

type ApprovalVoting struct {
	choices int
}

func (m *ApprovalVoting) Vote(choices ...int) structs.Ballot {
	b := make(structs.Ballot, m.choices)
	for _, c := range choices {
		b[c] = 1
	}
	return b
}

func (m *ApprovalVoting) Tally(ballots []structs.Ballot) Tally {
	tally := make(Tally, m.choices)
	for i := 0; i < m.choices; i++ {
		tally[i].Index = i
	}
loop:
	for _, b := range ballots {
		if len(b) != m.choices {
			continue
		}
		for _, approval := range b {
			if approval > 1 {
				continue loop
			}
		}
		for i, approval := range b {
			tally[i].Count += uint64(approval)
		}
	}
	return tally
}
