package parcomb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpectRune(t *testing.T) {
	var r interface{}
	var err error

	parser := Expect('a')

	r, _ = ParseRuneReader(parser, strings.NewReader("a"))
	assert.Equal(t, "a", r)
	r, _ = ParseRuneReader(parser, strings.NewReader("aaa"))
	assert.Equal(t, "a", r)

	_, err = ParseRuneReader(parser, strings.NewReader(""))
	assert.EqualError(t, err, `Expected "a"`)
	_, err = ParseRuneReader(parser, strings.NewReader("b"))
	assert.EqualError(t, err, `Expected "a"`)
	_, err = ParseRuneReader(parser, strings.NewReader("bbb"))
	assert.EqualError(t, err, `Expected "a"`)
}
