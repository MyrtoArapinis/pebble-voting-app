package vdf

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

const delta = 4096
const modulusBits = 1024

type PietrzakVdf struct {
	MaxDifficulty        uint64
	DifficultyConversion uint64
}

type intSerializer struct {
	Buffer []byte
}

func (s *intSerializer) WriteUint64(x uint64) {
	s.Buffer = append(s.Buffer, byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
}

func (s *intSerializer) Write(x *big.Int) {
	s.Buffer = append(s.Buffer, x.FillBytes(make([]byte, modulusBits/8))...)
}

func (s *intSerializer) ReadUint64() uint64 {
	if len(s.Buffer) < 8 {
		return ^uint64(0)
	}
	var res uint64 = 0
	for i := 0; i < 8; i++ {
		res = (res << 8) + uint64(s.Buffer[i])
	}
	s.Buffer = s.Buffer[8:]
	return res
}

func (s *intSerializer) Read() *big.Int {
	const intBytes = modulusBits / 8
	if len(s.Buffer) < intBytes {
		return nil
	}
	res := new(big.Int)
	res.SetBytes(s.Buffer[:intBytes])
	s.Buffer = s.Buffer[intBytes:]
	return res
}

type squarer interface {
	Eval(x *big.Int, t uint64) *big.Int
}

type repeatedSquarer struct {
	n *big.Int
}

func (s *repeatedSquarer) Eval(x *big.Int, t uint64) (r *big.Int) {
	e := new(big.Int)
	two := big.NewInt(2)
	e.Exp(two, big.NewInt(delta), nil)
	r = new(big.Int)
	r.Set(x)
	for t >= delta {
		r.Exp(r, e, s.n)
		t -= delta
	}
	if t != 0 {
		e.Exp(two, big.NewInt(int64(t)), nil)
		r.Exp(r, e, s.n)
	}
	return r
}

type trapdoorSquarer struct {
	n, phi *big.Int
}

func newTrapdoorSquarer(p, q *big.Int) (s *trapdoorSquarer) {
	s = new(trapdoorSquarer)
	s.n = new(big.Int)
	s.phi = new(big.Int)
	s.n.Mul(p, q)
	var p1, q1 big.Int
	one := big.NewInt(1)
	p1.Sub(p, one)
	q1.Sub(q, one)
	s.phi.Mul(&p1, &q1)
	return s
}

func (s *trapdoorSquarer) Eval(x *big.Int, t uint64) (r *big.Int) {
	var e big.Int
	r = new(big.Int)
	e.Exp(big.NewInt(2), big.NewInt(int64(t)), s.phi)
	r.Exp(x, &e, s.n)
	return r
}

func (vdf *PietrzakVdf) Create(seconds uint64) (sol VdfSolution, err error) {
	t := vdf.DifficultyConversion * seconds
	if t%2 != 0 {
		t++
	}
	if t > vdf.MaxDifficulty {
		t = vdf.MaxDifficulty
	}
	p, err := rand.Prime(rand.Reader, modulusBits/2)
	if err != nil {
		return sol, err
	}
	q, err := rand.Prime(rand.Reader, modulusBits/2)
	if err != nil {
		return sol, err
	}
	n := new(big.Int)
	n.Mul(p, q)
	x, err := rand.Int(rand.Reader, n)
	if err != nil {
		return sol, err
	}
	sqr := newTrapdoorSquarer(p, q)
	y := sqr.Eval(x, t)
	var w intSerializer
	w.WriteUint64(t)
	w.Write(n)
	w.Write(x)
	sol.Input = w.Buffer
	sol.Output = y.FillBytes(make([]byte, modulusBits/8))
	sol.Proof = sqr.phi.FillBytes(make([]byte, modulusBits/8))
	return
}

type transcript struct {
	hash [32]byte
}

func (tr *transcript) Init(t uint64) {
	b := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		b[i] = byte(t)
		t >>= 8
	}
	tr.hash = sha256.Sum256(b)
}

func (tr *transcript) Add(x *big.Int) {
	b := tr.hash[:]
	b = append(b, x.FillBytes(make([]byte, 256))...)
	tr.hash = sha256.Sum256(b)
}

func (tr *transcript) Challenge() (r *big.Int) {
	r = new(big.Int)
	r.SetBytes(tr.hash[:16])
	return
}

// computes a^b * c % n
func expAndMul(a, b, c, n *big.Int) (r *big.Int) {
	r = new(big.Int)
	r.Exp(a, b, n)
	r.Mul(r, c)
	r.Rem(r, n)
	return
}

func (vdf *PietrzakVdf) Solve(input []byte) (VdfSolution, error) {
	return vdf.solve(input, nil)
}

func (vdf *PietrzakVdf) solve(input []byte, sqr squarer) (sol VdfSolution, err error) {
	ser := intSerializer{input}
	t := ser.ReadUint64()
	n := ser.Read()
	x := ser.Read()
	if t > vdf.MaxDifficulty || n == nil || x == nil {
		err = newError("failed to parse vdf input")
		return
	}
	x.Rem(x, n)
	if sqr == nil {
		sqr = &repeatedSquarer{n}
	}
	y := sqr.Eval(x, t)
	sol.Input = input
	sol.Output = y.FillBytes(make([]byte, modulusBits/8))
	var tr transcript
	tr.Init(t)
	tr.Add(n)
	tr.Add(x)
	tr.Add(y)
	ser.Buffer = nil
	for t > delta {
		halfT := t / 2
		muRoot := sqr.Eval(x, halfT-1)
		ser.Write(muRoot)
		tr.Add(muRoot)
		r := tr.Challenge()
		mu := new(big.Int)
		mu.Mul(muRoot, muRoot)
		mu.Rem(mu, n)
		x = expAndMul(x, r, mu, n)
		y = expAndMul(mu, r, y, n)
		if halfT%2 == 0 {
			t = halfT
		} else {
			t = halfT + 1
			y.Mul(y, y)
			y.Rem(y, n)
		}
	}
	sol.Proof = ser.Buffer
	return
}

func (vdf *PietrzakVdf) Verify(sol VdfSolution) error {
	r := intSerializer{sol.Input}
	t := r.ReadUint64()
	n := r.Read()
	x := r.Read()
	if t > vdf.MaxDifficulty || n == nil || x == nil {
		return newError("failed to parse vdf input")
	}
	if t%2 != 0 {
		return newError("time difficulty not even")
	}
	x.Rem(x, n)
	y := new(big.Int)
	y.SetBytes(sol.Output)
	if y.Cmp(n) >= 0 {
		return newError("output greater than modulous")
	}
	if len(sol.Proof) == modulusBits/8 {
		phi := new(big.Int)
		phi.SetBytes(sol.Proof)
		if phi.IsUint64() && phi.Uint64() < 2 {
			return newError("invalid order given")
		}
		var one big.Int
		one.Exp(x, phi, n)
		if !(one.IsUint64() && one.Uint64() == 1) {
			return newError("invalid order given")
		}
		if (&trapdoorSquarer{n, phi}).Eval(x, t).Cmp(y) != 0 {
			return newError("trapdoor evaluation does not match output")
		}
		return nil
	}
	var tr transcript
	tr.Init(t)
	tr.Add(n)
	tr.Add(x)
	tr.Add(y)
	r.Buffer = sol.Proof
	mu := new(big.Int)
	for i := 0; t > delta; i++ {
		halfT := t / 2
		muRoot := r.Read()
		if muRoot == nil {
			return newError("failed to parse proof")
		}
		tr.Add(muRoot)
		r := tr.Challenge()
		mu.Mul(muRoot, muRoot)
		mu.Rem(mu, n)
		x = expAndMul(x, r, mu, n)
		y = expAndMul(mu, r, y, n)
		if halfT%2 == 0 {
			t = halfT
		} else {
			t = halfT + 1
			y.Mul(y, y)
			y.Rem(y, n)
		}
	}
	if (&repeatedSquarer{n}).Eval(x, t).Cmp(y) != 0 {
		return newError("final evaluation check failed")
	}
	return nil
}
