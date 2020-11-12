package combinators

import (
	"fmt"
	"unicode"
)

// Any ...
var Any = ExpectPredicate(func(r rune) bool { return true }, `any`)

// EOF ...
var EOF = FuncParser(func(state ParserState) (*ParserResult, error) {
	if state.CurrentRune() == 0 {
		return Success(state.Remaining(), string(state.CurrentRune()))
	}

	return Fail(state, fmt.Errorf(`Expected end of stream`))
})

// InlineSpace ...
var InlineSpace = ExpectPredicate(func(r rune) bool {
	return r == ' ' || r == '\t'
}, `inline space`)

// Space ...
var Space = ExpectPredicate(unicode.IsSpace, `space`)

// Newline ...
var Newline = ExpectPredicate(func(r rune) bool {
	return r == '\n' || r == '\r'
}, `newline`)

// Digit ...
var Digit = ExpectPredicate(unicode.IsDigit, `digit`)

// Letter ...
var Letter = ExpectPredicate(unicode.IsLetter, `letter`)

// Alphanumeric ...
var Alphanumeric = AnyOf(Letter, Digit)

// Integer ...
var Integer = Transform(
	AnyOf(
		SeqOf(
			ExpectAny([]rune("123456789")),
			ZeroOrMore(Digit),
		),
		Expect('0'),
	),
	func(iany interface{}) interface{} {
		if zero, ok := iany.(string); ok {
			return zero
		}

		seq := iany.([]interface{})
		expectAny := seq[0].(string)
		iZeroOrMore := seq[1].([]interface{})

		rem := ""
		for _, digit := range iZeroOrMore {
			rem += digit.(string)
		}

		return expectAny + rem
	})

// Decimal ...
var Decimal = Transform(
	SeqOf(
		Transform(
			Optional(ExpectAny([]rune("+-"))),
			func(i interface{}) interface{} {
				if i == nil {
					return ""
				}

				return i.(string)
			},
		),
		Integer,
		Expect('.'),
		OneOrMore(Digit),
	),
	func(iSeq interface{}) interface{} {
		seq := iSeq.([]interface{})
		iOneOrMore := seq[3].([]interface{})
		str := ""
		for _, digit := range iOneOrMore {
			str += digit.(string)
		}
		return seq[0].(string) + seq[1].(string) + seq[2].(string) + str
	},
)

// Decimal2 ...
var Decimal2 = Transform(
	SeqOf(
		Optional(ExpectAny([]rune("+-"))),
		Integer,
		Expect('.'),
		OneOrMore(Digit),
	),
	func(i interface{}) interface{} {
		return StringifyInterfaces(i)
	},
)

// StringifyInterfaces ...
func StringifyInterfaces(i interface{}) string {
	str := ""

	switch v := i.(type) {
	case []interface{}:
		for _, vv := range v {
			str += StringifyInterfaces(vv)
		}
	case string:
		str += v
	default:
		str += fmt.Sprint(v)
	}

	return str
}
