package minimark

import (
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/aziis98/parcomb"
	"github.com/aziis98/parcomb/examples/minimark/parser"
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
		r, _ := parcomb.ParseRuneReader(parser.Heading, strings.NewReader("### Prova"))
		assert.Equal(t, []interface{}{3, "Prova"}, r)
	}
	{
		r, _ := parcomb.ParseRuneReader(parser.Minimark, strings.NewReader(example1))

		json, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			t.Fatal(err)
		}

		log.Printf(`%s`, json)

		assert.Equal(t, []interface{}([]interface{}{
			"\n",
			[]interface{}{1, "Prova"},
			"\n",
			[]interface{}{2, "Prova"},
			"\n",
			[]interface{}{3, "Prova"},
			"\n", "\n",
			"Paragraph of some long text.\nParagraph of some long text.\nParagraph of some long text.",
			"\n",
			"\n",
			"Paragraph of some long text.",
			"\n",
			"\n",
			[]interface{}{
				[]interface{}{" - ", "First item of list"},
				[]interface{}{" - ", "Second item of this list"},
			},
			"\n",
		},
		), r)
	}
}
