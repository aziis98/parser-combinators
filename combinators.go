package parcomb

import (
	"fmt"
	// "log"
	"strings"
)

// Success creates a successfull ParserResult
func Success(state ParserState, result interface{}) (*ParserResult, error) {
	return &ParserResult{result, state}, nil
}

// Fail creates a "failing" ParserResult
func Fail(state ParserState, e error) (*ParserResult, error) {
	return &ParserResult{nil, state}, e
}

// Expect creates a Parser that expects a single given character and if successfull returns a string as Result
func Expect(expected rune) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {

		// state.(*scanner).PrintErrorMessage(fmt.Errorf(`Expect() :: "%c" == "%c"`, state.CurrentRune(), expected))

		if state.CurrentRune() != expected {
			return Fail(state, fmt.Errorf(`Expected "%c"`, expected))
		}

		return Success(state.Remaining(), string(expected))
	})
}

// ExpectPredicate creates a Parser for a rune based on given predicate function
func ExpectPredicate(predicate func(rune) bool, descriptions ...string) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		if predicate(state.CurrentRune()) {
			return Success(state.Remaining(), string(state.CurrentRune()))
		}

		return Fail(state, fmt.Errorf(`Expected "%v"`, descriptions))
	})
}

// ExpectAny creates a Parser that expects any rune from a given list
func ExpectAny(expectedList []rune) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		for _, expected := range expectedList {
			if state.CurrentRune() == expected {
				return Success(state.Remaining(), string(expected))
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
			if currentState.CurrentRune() != expected {
				return Fail(currentState, fmt.Errorf(`Expected "%c"`, expected))
			}

			currentState = currentState.Remaining()
		}

		return Success(currentState, string(expectedList))
	})
}

// SeqOf combines parsers in a seequence
func SeqOf(parsers ...Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		currentState := state
		results := []interface{}{}

		for _, parser := range parsers {

			// log.Printf(`SeqOf(%v)`, i)
			// currentState.(*scanner).PrintErrorMessage(nil)
			pr, err := parser.Apply(currentState)
			// pr.Remaining.(*scanner).PrintErrorMessage(nil)

			if err != nil {
				// log.Printf(`SeqOf() :: End`)
				return Fail(currentState, err)
			}

			results = append(results, pr.Result)
			currentState = pr.Remaining
		}

		// log.Printf(`SeqOf() :: End`)

		return Success(currentState, results)
	})
}

// AnyOf must match one of the given parsers
func AnyOf(parsers ...Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		errors := []string{}

		for _, parser := range parsers {

			// log.Printf(`AnyOf(%v)`, i)
			// state.(*scanner).PrintErrorMessage(nil)
			pr, err := parser.Apply(state)

			if err == nil {
				// log.Printf(`AnyOf() :: End`)
				return &ParserResult{pr.Result, pr.Remaining}, nil
			}

			// pr.Remaining.(*scanner).PrintErrorMessage(nil)

			errors = append(errors, fmt.Sprintf(" - %v", err))
		}

		// log.Printf(`AnyOf() :: End`)
		return Fail(state, fmt.Errorf("All cases failed:\n%s", strings.Join(errors, "\n")))
	})
}

// OneOrMore matches one or more of a given parser
func OneOrMore(parser Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		currentState := state
		results := []interface{}{}

		pr, err := parser.Apply(currentState)
		if err != nil {
			return Fail(currentState, err)
		}

		currentState = pr.Remaining

		for err == nil {
			results = append(results, pr.Result)
			pr, err = parser.Apply(currentState)
			currentState = pr.Remaining
		}

		return Success(currentState, results)
	})
}

// ZeroOrMore matches zero or more of a given parser
func ZeroOrMore(parser Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		currentState := state
		results := []interface{}{}

		// log.Printf(`ZeroOrMore(0)`)
		// currentState.(*scanner).PrintErrorMessage(nil)

		pr, err := parser.Apply(currentState)
		if err != nil {
			return &ParserResult{results, currentState}, nil
		}

		for err == nil {
			results = append(results, pr.Result)
			currentState = pr.Remaining

			// log.Printf(`ZeroOrMore(+1)`)
			// currentState.(*scanner).PrintErrorMessage(nil)

			pr, err = parser.Apply(currentState)
		}

		return &ParserResult{results, currentState}, nil
	})
}

// Optional matches zero or one of a given parser
func Optional(parser Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		pr, err := parser.Apply(state)
		if err != nil {
			return Success(state, nil)
		}

		return Success(pr.Remaining, pr.Result)
	})
}

// Transform a parser result if successfull
func Transform(parser Parser, transform func(interface{}) interface{}) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		pr, err := parser.Apply(state)

		// log.Printf(`Transform()`)
		// pr.Remaining.(*scanner).PrintErrorMessage(nil)
		// log.Printf(`before: %+v`, pr.Result)

		if err != nil {
			return Fail(state, err)
		}

		result := transform(pr.Result)
		// log.Printf(`after: %+v`, result)

		return &ParserResult{result, pr.Remaining}, nil
	})
}
