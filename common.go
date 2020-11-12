package parcomb

import (
	"unicode"
)

// Any ...
var Any = ExpectPredicate(func(r rune) bool { return true }, `any`)

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
		Integer,
		Expect('.'),
		OneOrMore(Digit),
	),
	func(iseq interface{}) interface{} {
		seq := iseq.([]interface{})
		iOneOrMore := seq[2].([]interface{})
		str := ""
		for _, digit := range iOneOrMore {
			str += digit.(string)
		}
		return seq[0].(string) + seq[1].(string) + str
	},
)
