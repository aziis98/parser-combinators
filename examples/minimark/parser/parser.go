package parser

import (
	"github.com/aziis98/parcomb"
	"github.com/aziis98/parcomb/examples/minimark/doc"
)

// Heading ...
var Heading = parcomb.Transform(
	parcomb.SeqOf(
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
	),
	func(iSeq interface{}) interface{} {
		seq := iSeq.([]interface{})
		level := seq[0].(int)
		text := seq[1].(string)
		return &doc.Heading{Level: level, Text: text}
	},
)

// Paragraph ...
var Paragraph = parcomb.Transform(
	parcomb.StringifyResult(
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
	),
	func(i interface{}) interface{} {
		text := i.(string)
		return &doc.Paragraph{Text: text}
	},
)

// List ...
var List = parcomb.Transform(
	parcomb.OneOrMore(
		Item,
	),
	func(i interface{}) interface{} {
		items := []*doc.Item{}
		for _, iItem := range i.([]interface{}) {
			items = append(items, &doc.Item{Depth: 0, Text: iItem.(string)})
		}
		return &doc.List{Items: items}
	},
)

// Item ...
var Item = parcomb.Transform(
	parcomb.SeqOf(
		parcomb.SeqIgnore(
			parcomb.ExpectString([]rune(" - ")),
		),
		parcomb.StringifyResult(
			parcomb.ZeroOrMore(
				parcomb.ExpectPredicate(func(r rune) bool {
					return r != '\n'
				}),
			),
		),
		parcomb.SeqIgnore(
			parcomb.Newline,
		),
	),
	func(i interface{}) interface{} {
		return i.([]interface{})[0]
	},
)

// Minimark ...
var Minimark = parcomb.Transform(
	parcomb.ZeroOrMore(
		parcomb.AnyOf(
			parcomb.Newline,
			Heading,
			List,
			Paragraph,
		),
	),
	func(i interface{}) interface{} {
		nodes := i.([]interface{})
		result := []doc.MinimarkNode{}

		for _, node := range nodes {
			if _, ok := node.(string); !ok {
				result = append(result, node.(doc.MinimarkNode))
			}
		}

		return result
	},
)
