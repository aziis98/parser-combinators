package parcomb

import (
	"io"
	"log"
)

// ParserState ...
type ParserState interface {
	CurrentRune() rune
	Remaining() ParserState
}

// ParserResult ...
type ParserResult struct {
	Result    interface{}
	Remaining ParserState
}

// IsError - trivial
func (result ParserResult) IsError() bool {
	_, ok := result.Result.(error)
	return ok
}

// Parser rappresents a generic parser combinator
type Parser interface {
	Apply(state ParserState) *ParserResult
}

// FuncParser functional binding for the Parser interface
type FuncParser func(state ParserState) *ParserResult

// Apply - trivial
func (p FuncParser) Apply(state ParserState) *ParserResult {
	return p(state)
}

type scanner struct {
	reader io.RuneReader
	buffer []rune
	cursor int

	line, col int
}

func (s *scanner) CurrentRune() rune {
	loops := 0

	for len(s.buffer) <= s.cursor {
		r, _, err := s.reader.ReadRune()
		log.Printf(`Read(%d) "%c"`, loops, r)

		if err != nil {
			return 0
		}

		if r == '\n' {
			s.line++
			s.col = 0
		} else {
			s.col++
		}

		s.buffer = append(s.buffer, r)
		loops++
	}

	return s.buffer[s.cursor]
}

func (s *scanner) Remaining() ParserState {
	return &scanner{
		s.reader,
		s.buffer,
		s.cursor + 1,

		s.line, s.col,
	}
}

// ParseRuneReader - trivial
func ParseRuneReader(parser Parser, r io.RuneReader) (interface{}, error) {
	scanner := &scanner{r, []rune{}, 0, 0, 0}

	log.Printf(`parser.Apply()`)
	rr := parser.Apply(scanner)

	err, ok := rr.Result.(error)
	if ok {
		return nil, err
	}

	return rr.Result, nil
}
