package parcomb

import (
	"fmt"
	"log"
	"strings"
)

// Success creates a successfull ParserResult
func Success(state ParserState, result interface{}) (*ParserResult, error) {
	return &ParserResult{result, state.Remaining()}, nil
}

// Fail creates a "failing" ParserResult
func Fail(state ParserState, e error) (*ParserResult, error) {
	return &ParserResult{nil, state}, e
}

// Expect creates a Parser that expects a single given character and if successfull returns a string as Result
func Expect(expected rune) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		if state.CurrentRune() != expected {
			return Fail(state, fmt.Errorf(`Expected "%c"`, expected))
		}

		return Success(state, string(expected))
	})
}

// ExpectAny creates a Parser that expects any rune from a given list
func ExpectAny(expectedList []rune) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		for _, expected := range expectedList {
			if state.CurrentRune() == expected {
				return Success(state, string(expected))
			}
		}

		return Fail(state, fmt.Errorf(`Expected one of %v`, strings.Join(strings.Split(string(expectedList), ""), ", ")))
	})
}

// ExpectString creates a Parser that expects all the runes from a given list
func ExpectString(expectedList []rune) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		currentState := state

		for _, expected := range expectedList {
			log.Printf(`%c == %c`, currentState.CurrentRune(), expected)
			if currentState.CurrentRune() != expected {
				return Fail(currentState, fmt.Errorf(`Expected "%c"`, expected))
			}

			currentState = currentState.Remaining()
		}

		return Success(currentState, string(expectedList))
	})
}

// Seq combines parsers in a seequence
func Seq(parsers ...Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		currentState := state
		results := []interface{}{}

		for _, parser := range parsers {
			pr, err := parser.Apply(currentState)
			if err != nil {
				return Fail(currentState, err)
			}

			results = append(results, pr.Result)
			currentState = pr.Remaining
		}

		return Success(currentState, results)
	})
}

// AnyOf must match one of the given parsers
func AnyOf(parsers ...Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		errors := []string{}

		for _, parser := range parsers {
			pr, err := parser.Apply(state)
			log.Printf(`%+v`, pr)

			if err == nil {
				return Success(state, pr.Result)
			}

			errors = append(errors, fmt.Sprintf(" - %v", err))
		}

		return Fail(state, fmt.Errorf("All cases failed:\n%s", strings.Join(errors, "\n")))
	})
}

func OneOrMore(parser Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		currentState := state
		results := []interface{}{}

		pr, err := parser.Apply(currentState)
		if err != nil {
			return Fail(state, err)
		}

		for err == nil {
			results = append(results, pr.Result)
			currentState = pr.Remaining

			pr, err = parser.Apply(currentState)
		}

		return Success(currentState, results)
	})
}
