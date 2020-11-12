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
func Fail(state ParserState, err error) (*ParserResult, error) {
	return &ParserResult{nil, state}, err
}

// Expect creates a Parser that expects a single given character and if successfull returns a string as Result
func Expect(expected rune) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		if state.CurrentRune() == 0 {
			return Fail(state, fmt.Errorf(`Stream ended, expected "%c"`, expected))
		}

		if state.CurrentRune() != expected {
			return Fail(state, fmt.Errorf(`Expected "%c"`, expected))
		}

		return Success(state.Remaining(), string(expected))
	})
}

// ExpectPredicate creates a Parser for a rune based on given predicate function
func ExpectPredicate(predicate func(rune) bool, descriptions ...string) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		if state.CurrentRune() == 0 {
			return Fail(state, fmt.Errorf(`Stream ended, expected "%v"`, descriptions))
		}

		// log.Printf(`predicate: "%c" is %v`, state.CurrentRune(), descriptions)

		if predicate(state.CurrentRune()) {
			return Success(state.Remaining(), string(state.CurrentRune()))
		}

		return Fail(state, fmt.Errorf(`Expected "%+v"`, descriptions))
	})
}

// ExpectAny creates a Parser that expects any rune from a given list
func ExpectAny(expectedList []rune) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		for _, expected := range expectedList {
			if state.CurrentRune() == 0 {
				return Fail(state, fmt.Errorf(`Stream ended, expected one of %v`, strings.Join(strings.Split(string(expectedList), ""), ", ")))
			}

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

		for i, expected := range expectedList {
			if state.CurrentRune() == 0 {
				return Fail(state, fmt.Errorf(`Stream ended, expected "%s"`, string(expectedList[i:])))
			}

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

			pr, err := parser.Apply(currentState)

			if err != nil {
				return Fail(currentState, err)
			}

			if _, ok := parser.(*seqIgnore); !ok {
				results = append(results, pr.Result)
			}
			currentState = pr.Remaining
		}

		return Success(currentState, results)
	})
}

type seqIgnore struct {
	Parser Parser
}

func (p seqIgnore) Apply(state ParserState) (*ParserResult, error) {
	return p.Parser.Apply(state)
}

// SeqIgnore wrapps an existing parser and ignores its result when used in "SeqOf"
func SeqIgnore(parser Parser) Parser {
	return &seqIgnore{parser}
}

// AnyOf must match one of the given parsers
func AnyOf(parsers ...Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		errors := []string{}

		for _, parser := range parsers {

			pr, err := parser.Apply(state)

			if err == nil {
				return Success(pr.Remaining, pr.Result)
			}

			errors = append(errors, fmt.Sprintf(" - %v", err))
		}

		return Fail(state, fmt.Errorf("All cases failed:\n%s", strings.Join(errors, "\n")))
	})
}

// RepeatUntil ...
func RepeatUntil(parser Parser, terminator Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		currentState := state
		results := []interface{}{}

		var err error
		_, err = terminator.Apply(currentState)

		for err != nil {
			pr, err2 := parser.Apply(currentState)
			if err2 != nil {
				return Fail(state, err2)
			}

			results = append(results, pr.Result)
			currentState = pr.Remaining

			_, err = terminator.Apply(currentState)
			// log.Printf(`until: %+v`, err)
		}

		// log.Printf(`stopped`)
		return Success(currentState, results)
	})
}

type Partial interface{}

// RestarableOneOrMore restarts the parsing from a given safepoint matched by
// another "restart" parser. For example (see examples for the precise
// definition of this parser, special modifiers omitted for brevity):
//
//  parser := RestarableOneOrMore(SeqOf(OneOrMore(Expect('a')), AnyOf(Newline, EOF)), Newline)
//  ParseRuneReader(parser, strings.NewReader("aaaa\naaaaa\naaabbbb\naaaaa\naa"))
//
// and the result is ["aaaa", "aaaaa", Partial("aaa"), "aaaaa", "aa"]
func RestarableOneOrMore(parser Parser, restart Parser) Parser {
	return FuncParser(func(state ParserState) (*ParserResult, error) {
		// currentState := state
		// results := []interface{}{}

		// currentState.(*scanner).PrintErrorMessage(nil)
		// pr, err := parser.Apply(currentState)
		// if pr.Remaining.CurrentRune() == 0 {
		// 	return Success(pr.Remaining, results)
		// }
		// if err != nil {
		// 	results = append(results, pr.Result)

		// 	currentState = pr.Remaining
		// 	pr, err = restart.Apply(currentState)
		// 	for err != nil {
		// 		log.Printf(`Recovering... "%c"`, currentState.CurrentRune())
		// 		if currentState.CurrentRune() == 0 {
		// 			return Fail(currentState, fmt.Errorf(`(2) Stream ended, %v`, err))
		// 		}
		// 		pr, err = restart.Apply(currentState)
		// 		currentState = currentState.Remaining()
		// 	}
		// 	log.Printf(`Parser recovered!`)
		// 	pr.Remaining.(*scanner).PrintErrorMessage(nil)

		// 	currentState = pr.Remaining
		// }

		// currentState = pr.Remaining

		// for err == nil {
		// 	results = append(results, pr.Result)

		// 	currentState.(*scanner).PrintErrorMessage(nil)
		// 	pr, err = parser.Apply(currentState)
		// 	pr.Remaining.(*scanner).PrintErrorMessage(nil)
		// 	if pr.Remaining.CurrentRune() == 0 {
		// 		return Success(pr.Remaining, results)
		// 	}
		// 	if err != nil {
		// 		results = append(results, pr.Result)

		// 		currentState = pr.Remaining
		// 		pr, err = restart.Apply(currentState)
		// 		for err != nil {
		// 			log.Printf(`Recovering... "%c"`, currentState.CurrentRune())
		// 			if currentState.CurrentRune() == 0 {
		// 				return Fail(currentState, fmt.Errorf(`(4) Stream ended, %v`, err))
		// 			}
		// 			pr, err = restart.Apply(currentState)
		// 			currentState = currentState.Remaining()
		// 		}
		// 		log.Printf(`Parser recovered!`)
		// 		pr.Remaining.(*scanner).PrintErrorMessage(nil)

		// 		currentState = pr.Remaining
		// 	}

		// 	currentState = pr.Remaining
		// }

		// return Success(currentState, results)
		return parser.Apply(state)
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

		pr, err := parser.Apply(currentState)
		if err != nil {
			return Success(currentState, results)
		}

		for currentState.CurrentRune() != 0 && err == nil {
			results = append(results, pr.Result)
			currentState = pr.Remaining

			pr, err = parser.Apply(currentState)
		}

		return Success(currentState, results)
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

		if err != nil {
			return Fail(state, err)
		}

		result := transform(pr.Result)

		return Success(pr.Remaining, result)
	})
}

// StringifyResult ...
func StringifyResult(parser Parser) Parser {
	return Transform(parser, func(i interface{}) interface{} {
		return StringifyInterfaces(i)
	})
}
