package parser

import "github.com/aziis98/parcomb"

// Heading ...
var Heading = parcomb.SeqOf(
	parcomb.Transform(
		parcomb.OneOrMore(parcomb.Expect('#')),
		func(i interface{}) interface{} {
			return len(i.([]interface{}))
		},
	),
	parcomb.SeqIgnore(parcomb.InlineSpace),
	parcomb.StringifyResult(
		parcomb.ZeroOrMore(
			parcomb.ExpectPredicate(func(r rune) bool {
				return r != '\n'
			}),
		),
	),
)

// Paragraph ...
var Paragraph = parcomb.StringifyResult(
	parcomb.RepeatUntil(
		parcomb.Any,
		parcomb.AnyOf(
			parcomb.SeqOf(
				parcomb.Expect('\n'),
				parcomb.AnyOf(
					parcomb.Expect('\n'),
					parcomb.EOF,
				),
			),
			parcomb.EOF,
		),
	),
)

// Itemize ...
var Itemize = parcomb.OneOrMore(
	parcomb.SeqOf(
		parcomb.ExpectString([]rune(" - ")),
		parcomb.StringifyResult(
			parcomb.ZeroOrMore(
				parcomb.ExpectPredicate(func(r rune) bool {
					return r != '\n'
				}),
			),
		),
		parcomb.SeqIgnore(parcomb.Newline),
	),
)

// Minimark ...
var Minimark = parcomb.ZeroOrMore(
	parcomb.AnyOf(
		parcomb.Newline,
		Heading,
		Itemize,
		Paragraph,
	),
)
