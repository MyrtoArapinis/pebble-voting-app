package voting

import (
	"time"

	"github.com/giry-dev/pebble-voting-app/voting/structs"
)

type ElectionPhase uint8

const (
	CredGen ElectionPhase = iota
	Cast
	Tally
)

type ElectionParams struct {
	Version               uint32
	EligibilityList       *structs.EligibilityList
	CastStart, TallyStart time.Time
	VdfDifficulty         uint64
	VotingMethod          string
	Choices               []string
}

func (p *ElectionParams) Phase() ElectionPhase {
	now := time.Now()
	if now.Before(p.CastStart) {
		return CredGen
	} else if now.Before(p.TallyStart) {
		return Cast
	} else {
		return Tally
	}
}
