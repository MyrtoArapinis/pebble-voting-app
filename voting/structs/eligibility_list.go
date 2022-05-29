package structs

import (
	"bytes"

	"github.com/giry-dev/pebble-voting-app/common"
	"github.com/giry-dev/pebble-voting-app/hashmap"
	"github.com/giry-dev/pebble-voting-app/util"
)

const structName = "EligibilityList"
const magic = 0x454c4c01

type PublicKey []byte

func (p PublicKey) Equals(other hashmap.Key) bool {
	switch o := other.(type) {
	case PublicKey:
		return bytes.Equal(p, o)
	default:
		return false
	}
}

func (p PublicKey) Hash() int {
	return hashmap.HashBytes(p)
}

type EligibilityList struct {
	publicKeys    []PublicKey
	idCommitments hashmap.Map // PublicKey->[]byte
}

func NewEligibilityList() *EligibilityList {
	return &EligibilityList{}
}

func (list *EligibilityList) Add(publicKey PublicKey, idCom []byte) bool {
	if _, exists := list.idCommitments.Get(publicKey); exists {
		return false
	}
	list.publicKeys = append(list.publicKeys, publicKey)
	list.idCommitments.Put(publicKey, idCom)
	return true
}

func (list *EligibilityList) IdCommitment(publicKey PublicKey) ([]byte, bool) {
	i, ok := list.idCommitments.Get(publicKey)
	if ok {
		switch v := i.(type) {
		case []byte:
			return v, true
		}
	}
	return nil, false
}

func (list *EligibilityList) Contains(publicKey PublicKey) bool {
	_, ok := list.idCommitments.Get(publicKey)
	return ok
}

func (list *EligibilityList) Bytes() []byte {
	var w util.BufferWriter
	w.WriteUint32(magic)
	for _, pk := range list.publicKeys {
		w.WriteVector(pk)
		c, _ := list.IdCommitment(pk)
		w.WriteVector(c)
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
	list.publicKeys = nil
	list.idCommitments.Clear()
	for r.Len() != 0 {
		var pk PublicKey
		pk, err = r.ReadVector()
		if err != nil {
			return err
		}
		c, err := r.ReadVector()
		if err != nil {
			return err
		}
		if !list.Add(pk, c) {
			return common.NewParsingError(structName, "duplicate public key")
		}
	}
	return nil
}
