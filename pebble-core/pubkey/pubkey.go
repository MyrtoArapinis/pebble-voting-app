package pubkey

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"strings"

	"blockwatch.cc/tzgo/tezos"
	"github.com/giry-dev/pebble-voting-app/pebble-core/base32c"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
)

type PublicKey []byte

type PrivateKey struct {
	p PublicKey
	s []byte
}

type KeyType byte

const (
	KeyTypeUnknown KeyType = iota
	KeyTypeEd25519
	KeyTypeTezos
)

var (
	ErrInvalidKeyLength = errors.New("pebble: invalid key length")

	ErrUnknownKeyType = errors.New("pebble: unknown key type")

	ErrInvalidSignature = errors.New("pebble: invalid signature")
)

func newPublicKey(t KeyType, k []byte) PublicKey {
	p := make(PublicKey, len(k)+1)
	p[0] = byte(t)
	copy(p[1:], k)
	return p
}

func (k PublicKey) Type() KeyType {
	if len(k) < 1 {
		return KeyTypeUnknown
	}
	return KeyType(k[0])
}

func (k PrivateKey) Type() KeyType {
	return k.p.Type()
}

func (k PrivateKey) Public() PublicKey {
	return k.p
}

func (k PrivateKey) Secret() []byte {
	return k.s
}

func NewKeyFromSeed(seed []byte) PrivateKey {
	priv := ed25519.NewKeyFromSeed(seed)
	return PrivateKey{
		p: newPublicKey(KeyTypeEd25519, priv[32:]),
		s: priv[:32],
	}
}

func GenerateKey(keyType KeyType) (PrivateKey, error) {
	switch keyType {
	case KeyTypeEd25519:
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return PrivateKey{}, err
		}
		return PrivateKey{newPublicKey(keyType, pub), priv.Seed()}, nil
	case KeyTypeTezos:
		priv, err := tezos.GenerateKey(tezos.KeyTypeEd25519)
		if err != nil {
			return PrivateKey{}, err
		}
		return PrivateKey{
			p: newPublicKey(keyType, priv.Public().Address().Bytes()),
			s: []byte(priv.String()),
		}, nil
	default:
		return PrivateKey{}, ErrUnknownKeyType
	}
}

func (k PrivateKey) Sign(msg []byte) ([]byte, error) {
	switch k.Type() {
	case KeyTypeEd25519:
		return ed25519.Sign(ed25519.NewKeyFromSeed(k.s), msg), nil
	case KeyTypeTezos:
		key, err := tezos.ParsePrivateKey(string(k.s))
		if err != nil {
			return nil, err
		}
		hash := util.Hash(msg)
		sig, err := key.Sign(hash[:])
		if err != nil {
			return nil, err
		}
		pub := key.Public().Bytes()
		res := make([]byte, 1, 120)
		res[0] = byte(len(pub))
		res = append(res, pub...)
		res = append(res, sig.Bytes()...)
		return res, nil
	default:
		return nil, ErrUnknownKeyType
	}
}

func (k PublicKey) Verify(msg, sig []byte) error {
	if len(k) == 0 {
		return ErrInvalidKeyLength
	}
	switch KeyType(k[0]) {
	case KeyTypeEd25519:
		pk := ed25519.PublicKey(k[1:])
		if len(pk) != ed25519.PublicKeySize {
			return ErrInvalidKeyLength
		}
		if !ed25519.Verify(pk, msg, sig) {
			return ErrInvalidSignature
		}
		return nil
	case KeyTypeTezos:
		var addr tezos.Address
		err := addr.UnmarshalBinary(k[1:])
		if err != nil {
			return err
		}
		if len(sig) == 0 {
			return ErrInvalidSignature
		}
		pubLen := int(sig[1])
		if len(sig) <= pubLen {
			return ErrInvalidSignature
		}
		pk, err := tezos.DecodeKey(sig[1 : pubLen+1])
		if err != nil {
			return ErrInvalidSignature
		}
		var tzsig tezos.Signature
		err = tzsig.UnmarshalBinary(sig[pubLen+1:])
		if err != nil {
			return ErrInvalidSignature
		}
		hash := util.Hash(msg)
		err = pk.Verify(hash[:], tzsig)
		if err != nil {
			return ErrInvalidSignature
		}
		return nil
	default:
		return ErrUnknownKeyType
	}
}

func (k PublicKey) String() (string, error) {
	if len(k) == 0 {
		return "", ErrInvalidKeyLength
	}
	switch KeyType(k[0]) {
	case KeyTypeEd25519:
		p := make([]byte, 2, len(k)+1)
		p[0] = 238
		p[1] = 78
		p = append(p, k[1:]...)
		return base32c.CheckEncode(p), nil
	case KeyTypeTezos:
		pk, err := tezos.DecodeKey(k[1:])
		if err != nil {
			return "", err
		}
		return pk.String(), nil
	default:
		return "", ErrUnknownKeyType
	}
}

func ParsePublicKey(s string) (PublicKey, error) {
	if strings.HasPrefix(s, "EPK") {
		p, err := base32c.CheckDecode(s)
		if err != nil {
			return nil, err
		}
		if len(p) < 3 || p[0] != 238 || p[1] != 78 {
			return nil, ErrUnknownKeyType
		}
		return PublicKey(p[2:]), nil
	} else if strings.HasPrefix(s, "tz") {
		addr, err := tezos.ParseAddress(s)
		if err != nil {
			return nil, err
		}
		return newPublicKey(KeyTypeTezos, addr.Bytes()), nil
	}
	return nil, ErrUnknownKeyType
}
