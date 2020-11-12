package parcomb

import (
	"fmt"
	"io"
	"log"
)

type ParserState interface {
	CurrentRune() rune
	Remaining() ParserState
}

type ParserResult struct {
	Result    interface{}
	Remaining ParserState
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

// Success creates a successfull ParserResult
func Success(state ParserState, result interface{}) *ParserResult {
	return &ParserResult{result, state.Remaining()}
}

// Fail creates a "failing" ParserResult
func Fail(state ParserState, e error) *ParserResult {
	return &ParserResult{e, state.Remaining()}
}

// Expect creates a Parser that expects a single given character and if successfull returns a string as Result
func Expect(expected rune) Parser {
	return FuncParser(func(state ParserState) *ParserResult {
		if state.CurrentRune() != expected {
			return Fail(state, fmt.Errorf(`Expected "%c"`, expected))
		}

		return Success(state, string(expected))
	})
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
