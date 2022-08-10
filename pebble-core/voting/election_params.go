package voting

import (
	"time"

	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type ElectionPhase uint8

const (
	CredGen ElectionPhase = iota
	Cast
	Tally
	End
)

type ElectionID = [32]byte

type ElectionParams struct {
	Version                         uint32
	Id                              ElectionID
	EligibilityList                 *structs.EligibilityList
	CastStart, TallyStart, TallyEnd time.Time
	VdfDifficulty                   uint64
	VotingMethod                    string
	Choices                         []string
}

func (p *ElectionParams) Phase() ElectionPhase {
	now := time.Now()
	if now.Before(p.CastStart) {
		return CredGen
	} else if now.Before(p.TallyStart) {
		return Cast
	} else if now.Before(p.TallyEnd) {
		return Tally
	} else {
		return End
	}
}
