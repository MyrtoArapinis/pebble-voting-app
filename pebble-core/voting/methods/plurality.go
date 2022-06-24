package methods

import "github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"

type PluralityVoting struct {
	choices int
}

func (m *PluralityVoting) Vote(choices ...int) structs.Ballot {
	if len(choices) != 1 {
		panic("more than one choice in plurality voting")
	}
	return []byte{byte(choices[0])}
}

func (m *PluralityVoting) Tally(ballots []structs.Ballot) Tally {
	tally := make(Tally, m.choices)
	for i := 0; i < m.choices; i++ {
		tally[i].Index = i
	}
	for _, b := range ballots {
		if len(b) != 1 {
			continue
		}
		c := int(b[0])
		if c > m.choices {
			continue
		}
		tally[c].Count++
	}
	return tally
}
