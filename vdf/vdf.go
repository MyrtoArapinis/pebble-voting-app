package vdf

type VdfSolution struct {
	Input, Output, Proof []byte
}

type VDF interface {
	Create(seconds uint64) (VdfSolution, error)
	Solve(input []byte) (VdfSolution, error)
	Verify(sol VdfSolution) error
}

type VdfError struct {
	s string
}

func (e *VdfError) Error() string {
	return e.s
}

func newError(s string) error {
	return &VdfError{s}
}
