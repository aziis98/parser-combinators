package parser

import (
	c "github.com/aziis98/parser-combinators"
	"github.com/aziis98/parser-combinators/examples/minimark/doc"
)

// Heading ...
var Heading = c.Transform(
	c.SeqOf(
		c.Transform(
			c.OneOrMore(c.Expect('#')),
			func(i interface{}) interface{} {
				return len(i.([]interface{}))
			},
		),
		c.SeqIgnore(c.InlineSpace),
		c.StringifyResult(
			c.ZeroOrMore(
				c.ExpectPredicate(func(r rune) bool {
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
var Paragraph = c.Transform(
	c.StringifyResult(
		c.RepeatUntil(
			c.Any,
			c.AnyOf(
				c.SeqOf(
					c.Expect('\n'),
					c.AnyOf(
						c.Expect('\n'),
						c.EOF,
					),
				),
				c.EOF,
			),
		),
	),
	func(i interface{}) interface{} {
		text := i.(string)
		return &doc.Paragraph{Text: text}
	},
)

// List ...
var List = c.Transform(
	c.OneOrMore(
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
var Item = c.Transform(
	c.SeqOf(
		c.SeqIgnore(
			c.ExpectString([]rune(" - ")),
		),
		c.StringifyResult(
			c.ZeroOrMore(
				c.ExpectPredicate(func(r rune) bool {
					return r != '\n'
				}),
			),
		),
		c.SeqIgnore(
			c.Newline,
		),
	),
	func(i interface{}) interface{} {
		return i.([]interface{})[0]
	},
)

// Minimark ...
var Minimark = c.Transform(
	c.ZeroOrMore(
		c.AnyOf(
			c.Newline,
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
