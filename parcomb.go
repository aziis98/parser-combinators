package parcomb

import (
	"fmt"
	"io"
	"os"
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
}

func (s *scanner) GetLocation() (int, int) {
	lines := strings.Split(string((*s.buffer)[:s.cursor]), "\n")
	lastLine := lines[len(lines)-1]
	return len(lines) - 1, len(lastLine) + 1
}

func (s *scanner) CurrentRune() rune {
	loops := 0

	for len(*s.buffer) <= s.cursor {
		r, _, err := s.reader.ReadRune()

		if err != nil {
			return 0
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

	line, col := s.GetLocation()

	loc := fmt.Sprintf(`%d:%d`, line, col)

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
	}
	fmt.Fprintf(os.Stderr, "at %s\n", loc)
	fmt.Fprintf(os.Stderr, " | %s\n", lines[line])
	fmt.Fprintf(os.Stderr, "  %s^\n", strings.Repeat(" ", col))
}

// ParseRuneReader - trivial
func ParseRuneReader(parser Parser, r io.RuneReader) (interface{}, error) {
	s := &scanner{r, &[]rune{}, 0}

	pr, err := parser.Apply(s)
	if err != nil {
		return nil, err
	}

	return pr.Result, nil
}
