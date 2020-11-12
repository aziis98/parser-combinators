package minimark

import (
	"strings"
	"testing"

	"github.com/aziis98/parcomb"
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

 - Paragraph of some long text.
 - Paragraph of some long text.

`

func TestMinimark(t *testing.T) {
	{
		r, _ := parcomb.ParseRuneReader(Heading, strings.NewReader("### Prova"))
		assert.Equal(t, []interface{}{3, "Prova"}, r)
	}
	{
		r, _ := parcomb.ParseRuneReader(Minimark, strings.NewReader(example1))
		assert.Equal(t, "", r)
	}
}
