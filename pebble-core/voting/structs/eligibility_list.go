package structs

import (
	"errors"

	"github.com/giry-dev/pebble-voting-app/pebble-core/common"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
)

const structName = "EligibilityList"
const magic = 0x454c4c01

var ErrDuplicateKey = errors.New("pebble: duplicate key")

type EligibilityList struct {
	publicKeyHashes []util.HashValue
	idCommitments   map[util.HashValue]util.HashValue
}

func (list *EligibilityList) Add(pkh, idCom util.HashValue) bool {
	if _, exists := list.idCommitments[pkh]; exists {
		return false
	}
	list.publicKeyHashes = append(list.publicKeyHashes, pkh)
	list.idCommitments[pkh] = idCom
	return true
}

func (list *EligibilityList) IdCommitment(pkh util.HashValue) (idCom util.HashValue, ok bool) {
	idCom, ok = list.idCommitments[pkh]
	return
}

func (list *EligibilityList) Contains(pkh util.HashValue) bool {
	_, ok := list.idCommitments[pkh]
	return ok
}

func (list *EligibilityList) Bytes() []byte {
	var w util.BufferWriter
	w.WriteUint32(magic)
	for _, pkh := range list.publicKeyHashes {
		w.Write(pkh[:])
		c, _ := list.IdCommitment(pkh)
		w.Write(c[:])
	}
	return w.Buffer
}

func (list *EligibilityList) FromBytes(p []byte) error {
	r := util.NewBufferReader(p)
	m, err := r.ReadUint32()
	if err != nil {
		return err
	}
	if m != magic {
		return common.NewParsingError(structName, "unknown magic")
	}
	list.publicKeyHashes = nil
	list.idCommitments = nil
	for r.Len() != 0 {
		pkh, err := r.Read32()
		if err != nil {
			return err
		}
		idCom, err := r.Read32()
		if err != nil {
			return err
		}
		if !list.Add(pkh, idCom) {
			return ErrDuplicateKey
		}
	}
	return nil
}
