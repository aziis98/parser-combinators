package parcomb

import (
	"fmt"
	"io"
	"strings"
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

// Parser rappresents a generic parser combinator
type Parser interface {
	Apply(state ParserState) (*ParserResult, error)
}

// FuncParser functional binding for the Parser interface
type FuncParser func(state ParserState) (*ParserResult, error)

// Apply - trivial
func (p FuncParser) Apply(state ParserState) (*ParserResult, error) {
	return p(state)
}

type scanner struct {
	reader io.RuneReader
	buffer *[]rune
	cursor int

	line, col int
}

func (s *scanner) CurrentRune() rune {
	loops := 0

	for len(*s.buffer) <= s.cursor {
		r, _, err := s.reader.ReadRune()
		// log.Printf(`Read(%d) "%c"`, loops, r)

		if err != nil {
			return 0
		}

		if r == '\n' {
			s.line++
			s.col = 0
		} else {
			s.col++
		}

		*s.buffer = append(*s.buffer, r)
		loops++
	}

	return (*s.buffer)[s.cursor]
}

func (s *scanner) Remaining() ParserState {
	return &scanner{
		s.reader,
		s.buffer,
		s.cursor + 1,

		s.line, s.col,
	}
}

func (s *scanner) PrintErrorMessage(e error) {
	var r rune
	var err error

	r, _, err = s.reader.ReadRune()

	for err == nil {
		*s.buffer = append(*s.buffer, r)
		r, _, err = s.reader.ReadRune()
	}

	lines := strings.Split(string(*s.buffer), "\n")
	loc := fmt.Sprintf(`%d:%d`, s.line, s.col)

	fmt.Printf("%v\n", e)
	fmt.Printf("at %s\n", loc)
	fmt.Printf(" | %s\n", lines[s.line])
	fmt.Printf("  %s^\n", strings.Repeat(" ", s.col))
}

// ParseRuneReader - trivial
func ParseRuneReader(parser Parser, r io.RuneReader) (interface{}, error) {
	s := &scanner{r, &[]rune{}, 0, 0, 0}

	pr, err := parser.Apply(s)
	if err != nil {
		errScanner := pr.Remaining.(*scanner)
		errScanner.PrintErrorMessage(err)

		return nil, err
	}

	return pr.Result, nil
}
