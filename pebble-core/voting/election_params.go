package voting

import (
	"errors"
	"time"

	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

var errUnknownVersion = errors.New("pebble: unknown ElectionParams version")

type ElectionPhase uint8

const (
	Setup ElectionPhase = iota
	CredGen
	Cast
	Tally
	End
)

type ElectionParams struct {
	Version                         uint32
	CastStart, TallyStart, TallyEnd time.Time
	MaxVdfDifficulty                uint64
	VotingMethod                    string
	Title, Description              string
	Choices                         []string
	EligibilityList                 *structs.EligibilityList
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

func (p *ElectionParams) Bytes() []byte {
	var w util.BufferWriter
	w.WriteUint32(p.Version)
	w.WriteUint64(uint64(p.CastStart.Unix()))
	w.WriteUint64(uint64(p.TallyStart.Unix()))
	w.WriteUint64(uint64(p.TallyEnd.Unix()))
	w.WriteUint64(p.MaxVdfDifficulty)
	w.WriteVector([]byte(p.VotingMethod))
	w.WriteVector([]byte(p.Title))
	w.WriteVector([]byte(p.Description))
	w.WriteByte(byte(len(p.Choices)))
	for _, c := range p.Choices {
		w.WriteVector([]byte(c))
	}
	w.Write(p.EligibilityList.Bytes())
	return w.Buffer
}

func (p *ElectionParams) FromBytes(b []byte) (err error) {
	r := util.NewBufferReader(b)
	p.Version, err = r.ReadUint32()
	if err != nil {
		return err
	}
	if p.Version != 0 {
		return errUnknownVersion
	}
	if err != nil {
		return err
	}
	t, err := r.ReadUint64()
	if err != nil {
		return err
	}
	p.CastStart = time.Unix(int64(t), 0)
	t, err = r.ReadUint64()
	if err != nil {
		return err
	}
	p.TallyStart = time.Unix(int64(t), 0)
	t, err = r.ReadUint64()
	if err != nil {
		return err
	}
	p.TallyEnd = time.Unix(int64(t), 0)
	p.MaxVdfDifficulty, err = r.ReadUint64()
	if err != nil {
		return err
	}
	b, err = r.ReadVector()
	if err != nil {
		return err
	}
	p.VotingMethod = string(b)
	b, err = r.ReadVector()
	if err != nil {
		return err
	}
	p.Title = string(b)
	b, err = r.ReadVector()
	if err != nil {
		return err
	}
	p.Description = string(b)
	numChoices, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Choices = make([]string, numChoices)
	for i := range p.Choices {
		b, err = r.ReadVector()
		if err != nil {
			return err
		}
		p.Choices[i] = string(b)
	}
	p.EligibilityList = structs.NewEligibilityList()
	err = p.EligibilityList.FromBytes(r.ReadRemaining())
	return err
}
