package minimark

import (
	"strings"
	"testing"

	c "github.com/aziis98/parser-combinators"
	"github.com/aziis98/parser-combinators/examples/minimark/doc"
	"github.com/aziis98/parser-combinators/examples/minimark/parser"
	"github.com/stretchr/testify/assert"
)

var example1 = `
# Prova
## Prova
### Prova

Paragraph of some long text.
Paragraph of some long text.
Paragraph of some long text.

Paragraph of some long text.

 - First item of list
 - Second item of this list

`

func TestMinimark(t *testing.T) {
	{
		r, _ := c.ParseRuneReader(parser.Heading, strings.NewReader("### Prova"))
		assert.Equal(t, &doc.Heading{Level: 3, Text: "Prova"}, r)
	}
	{
		r, _ := c.ParseRuneReader(parser.Minimark, strings.NewReader(example1))

		// json, err := json.MarshalIndent(r, "", "  ")
		// if err != nil {
		// 	t.Fatal(err)
		// }

		// log.Printf(`%s`, json)

		assert.Equal(t,
			[]doc.MinimarkNode{
				&doc.Heading{Level: 1, Text: "Prova"},
				&doc.Heading{Level: 2, Text: "Prova"},
				&doc.Heading{Level: 3, Text: "Prova"},
				&doc.Paragraph{Text: "Paragraph of some long text.\nParagraph of some long text.\nParagraph of some long text."},
				&doc.Paragraph{Text: "Paragraph of some long text."},
				&doc.List{
					Items: []*doc.Item{
						{Depth: 0, Text: "First item of list"},
						{Depth: 0, Text: "Second item of this list"},
					},
				},
			},
			r,
		)
	}
}
