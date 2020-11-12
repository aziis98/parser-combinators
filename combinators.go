package parcomb

import (
	"fmt"
)

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

func ExpectAny(expectedList []rune) Parser {
	return FuncParser(func(state ParserState) *ParserResult {
		for _, expected := range expectedList {
			if state.CurrentRune() == expected {
				return Success(state, string(expected))
			}
		}

		return Fail(state, fmt.Errorf(`Expected one of "%v"`, string(expectedList)))
	})
}

func Seq(parsers ...Parser) {
	return FuncParser(func(state ParserState) *ParserResult {
		for _, parser := range parsers {
			pr := parser.Apply(state)
			if pr.IsError() {
				return Fail(state, error(pr.Result))
			}

			state = pr.Remaining
		}
	})
}
