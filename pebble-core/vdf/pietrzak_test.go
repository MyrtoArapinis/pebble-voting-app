package vdf

import (
	"math/rand"
	"testing"
)

func TestSolvePietrzak(t *testing.T) {
	var vdf VDF = &PietrzakVdf{1 << 63, 10000}
	puz, err := vdf.Create(uint64(1 + rand.Int31n(10)))
	if err != nil {
		t.Error(err.Error())
	}
	err = vdf.Verify(puz)
	if err != nil {
		t.Error(err.Error())
	}
	sol, err := vdf.Solve(puz.Input)
	if err != nil {
		t.Error(err.Error())
	}
	err = vdf.Verify(sol)
	if err != nil {
		t.Error(err.Error())
	}
}
