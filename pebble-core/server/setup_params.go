package server

import (
	"time"

	"github.com/giry-dev/pebble-voting-app/pebble-core/pubkey"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting"
	"github.com/giry-dev/pebble-voting-app/pebble-core/voting/structs"
)

type ElectionSetupVoter struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type ElectionSetupParams struct {
	AdminId       string               `json:"adminId"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	VoteStart     string               `json:"voteStart"`
	VoteEnd       string               `json:"voteEnd"`
	VdfDifficulty string               `json:"vdfDifficulty"`
	Method        string               `json:"method"`
	Choices       []string             `json:"choices"`
	Voters        []ElectionSetupVoter `json:"voters"`
}

func (sp *ElectionSetupParams) Params() (*voting.ElectionParams, error) {
	castStart, err := time.Parse(time.RFC3339, sp.VoteStart)
	if err != nil {
		return nil, err
	}
	tallyStart, err := time.Parse(time.RFC3339, sp.VoteEnd)
	if err != nil {
		return nil, err
	}
	tallyEnd := tallyStart.Add(tallyStart.Sub(castStart) * 100)
	ep := &voting.ElectionParams{
		Version:         0,
		CastStart:       castStart,
		TallyStart:      tallyStart,
		TallyEnd:        tallyEnd,
		Title:           sp.Title,
		Description:     sp.Description,
		VotingMethod:    sp.Method,
		Choices:         sp.Choices,
		EligibilityList: structs.NewEligibilityList(),
	}
	for _, voter := range sp.Voters {
		pk, err := pubkey.Parse(voter.Key)
		if err != nil {
			continue
		}
		idCom := util.Hash([]byte(voter.Id))
		ep.EligibilityList.Add(util.Hash(pk), idCom)
	}
	return ep, nil
}
