package common

type ParsingError struct {
	Struct string
	Msg    string
}

func (e *ParsingError) Error() string {
	return "error parsing " + e.Struct + ": " + e.Msg
}

func NewParsingError(structName, msg string) error {
	return &ParsingError{structName, msg}
}
