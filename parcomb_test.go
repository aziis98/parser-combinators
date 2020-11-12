package combinators

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestExpectRune(t *testing.T) {
	var r interface{}
	var err error

	parser := Expect('a')

	r, _ = ParseRuneReader(parser, strings.NewReader("a"))
	assert.Equal(t, "a", r)
	r, _ = ParseRuneReader(parser, strings.NewReader("aaa"))
	assert.Equal(t, "a", r)

	_, err = ParseRuneReader(parser, strings.NewReader(""))
	assert.EqualError(t, err, `Stream ended, expected "a"`)
	_, err = ParseRuneReader(parser, strings.NewReader("b"))
	assert.EqualError(t, err, `Expected "a"`)
	_, err = ParseRuneReader(parser, strings.NewReader("bbb"))
	assert.EqualError(t, err, `Expected "a"`)
}

func TestExpectAny(t *testing.T) {
	var r interface{}
	var err error

	parser := ExpectAny([]rune("abc"))

	r, _ = ParseRuneReader(parser, strings.NewReader("a"))
	assert.Equal(t, "a", r)
	r, _ = ParseRuneReader(parser, strings.NewReader("b"))
	assert.Equal(t, "b", r)
	r, _ = ParseRuneReader(parser, strings.NewReader("c"))
	assert.Equal(t, "c", r)

	_, err = ParseRuneReader(parser, strings.NewReader(""))
	assert.EqualError(t, err, `Stream ended, expected one of a, b, c`)
	_, err = ParseRuneReader(parser, strings.NewReader("d"))
	assert.EqualError(t, err, `Expected one of a, b, c`)
}

func TestSeq(t *testing.T) {
	parser1 := SeqOf(Expect('a'), Expect('a'))
	{
		r, _ := ParseRuneReader(parser1, strings.NewReader("aa"))
		assert.EqualValues(t, []interface{}{"a", "a"}, r)
	}
	{
		r, _ := ParseRuneReader(parser1, strings.NewReader("aaa"))
		assert.EqualValues(t, []interface{}{"a", "a"}, r)
	}
	{
		_, err := ParseRuneReader(parser1, strings.NewReader("a"))
		assert.EqualError(t, err, `Stream ended, expected "a"`)
	}
	{
		_, err := ParseRuneReader(parser1, strings.NewReader("ababab"))
		assert.EqualError(t, err, `Expected "a"`)
	}
}

func TestAny(t *testing.T) {
	parser1 := AnyOf(Expect('a'), Expect('b'))
	{
		r, _ := ParseRuneReader(parser1, strings.NewReader("aa"))
		assert.Equal(t, "a", r)
	}
	{
		r, _ := ParseRuneReader(parser1, strings.NewReader("bb"))
		assert.Equal(t, "b", r)
	}
	{
		_, err := ParseRuneReader(parser1, strings.NewReader(""))
		assert.EqualError(t, err, "All cases failed:\n - Stream ended, expected \"a\"\n - Stream ended, expected \"b\"")
	}
	{
		_, err := ParseRuneReader(parser1, strings.NewReader("ccc"))
		assert.EqualError(t, err, "All cases failed:\n - Expected \"a\"\n - Expected \"b\"")
	}
}

func TestDigit(t *testing.T) {
	{
		digit := AnyOf(Expect('0'), Expect('1'), Expect('2'), Expect('3'), Expect('4'), Expect('5'), Expect('6'), Expect('7'), Expect('8'), Expect('9'))
		r, _ := ParseRuneReader(digit, strings.NewReader("5"))
		assert.Equal(t, "5", r)
	}
}

func TestDigits(t *testing.T) {
	digit := ExpectAny([]rune("0123456789"))
	digits := OneOrMore(digit)

	{
		r, _ := ParseRuneReader(digits, strings.NewReader("1"))
		assert.EqualValues(t, []interface{}{"1"}, r)
	}
	{
		r, _ := ParseRuneReader(digits, strings.NewReader("15"))
		assert.EqualValues(t, []interface{}{"1", "5"}, r)
	}
	{
		r, _ := ParseRuneReader(digits, strings.NewReader("15Kb"))
		assert.EqualValues(t, []interface{}{"1", "5"}, r)
	}
	{
		r, _ := ParseRuneReader(digits, strings.NewReader("123456789000kg"))
		assert.EqualValues(t, []interface{}{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "0", "0"}, r)
	}
	{
		_, err := ParseRuneReader(digits, strings.NewReader("abcabc"))
		assert.EqualError(t, err, "Expected one of 0, 1, 2, 3, 4, 5, 6, 7, 8, 9")
	}
}

func TestStrings(t *testing.T) {
	word := ExpectString([]rune("symbol"))

	{
		r, _ := ParseRuneReader(word, strings.NewReader("symbol"))
		assert.Equal(t, "symbol", r)
	}
	{
		r, _ := ParseRuneReader(word, strings.NewReader("symbol1111111"))
		assert.Equal(t, "symbol", r)
	}
	{
		_, err := ParseRuneReader(word, strings.NewReader("symb0l"))
		assert.EqualError(t, err, `Expected "o"`)
	}
}

func TestMultipleStrings(t *testing.T) {
	word1 := ExpectString([]rune("foo"))
	word2 := ExpectString([]rune("symbol"))
	word3 := ExpectString([]rune("word"))
	parser := AnyOf(word1, word2, word3)

	{
		r, _ := ParseRuneReader(parser, strings.NewReader("foo"))
		assert.Equal(t, "foo", r)
	}
	{
		r, _ := ParseRuneReader(parser, strings.NewReader("symbol"))
		assert.Equal(t, "symbol", r)
	}
	{
		r, _ := ParseRuneReader(parser, strings.NewReader("word"))
		assert.Equal(t, "word", r)
	}
}

func TestMultipleStrings2(t *testing.T) {
	// Before, scanner.buffer was just a []rune and not a pointer, this caused
	//  backtracking to not work and for example this test matched "f" instead
	//  of  "foo". Now all the  "scanner"  instances  have the  buffer  all in
	//  common so this is no longer a  problem  (and might also be more memory
	//	efficient)

	word1 := ExpectString([]rune("fooer"))
	word2 := ExpectString([]rune("foo"))
	word3 := ExpectString([]rune("f"))
	parser := AnyOf(word1, word2, word3)

	{
		r, _ := ParseRuneReader(parser, strings.NewReader("foooer"))
		assert.Equal(t, "foo", r)
	}
}

func TestCommon(t *testing.T) {
	{
		r, _ := ParseRuneReader(Integer, strings.NewReader("7"))
		assert.Equal(t, "7", r)
	}
	{
		r, _ := ParseRuneReader(Integer, strings.NewReader("123"))
		assert.Equal(t, "123", r)
	}
	{
		r, _ := ParseRuneReader(Integer, strings.NewReader("0123"))
		assert.Equal(t, "0", r)
	}
	{
		r, _ := ParseRuneReader(Decimal, strings.NewReader("3.14"))
		assert.Equal(t, "3.14", r)
	}
	{
		r, _ := ParseRuneReader(Decimal, strings.NewReader("0.681"))
		assert.Equal(t, "0.681", r)
	}
	{
		r, _ := ParseRuneReader(Decimal, strings.NewReader("123.456"))
		assert.Equal(t, "123.456", r)
	}
	{
		r, _ := ParseRuneReader(Decimal, strings.NewReader("+123.456"))
		assert.Equal(t, "+123.456", r)
	}
	{
		r, _ := ParseRuneReader(Decimal, strings.NewReader("-123.456"))
		assert.Equal(t, "-123.456", r)
	}
	{
		r, _ := ParseRuneReader(Decimal2, strings.NewReader("-123.456"))
		assert.Equal(t, "-123.456", r)
	}
}

func BenchmarkDecimal0(b *testing.B) {
	var r interface{}

	for n := 0; n < b.N; n++ {
		r, _ = ParseRuneReader(Decimal, strings.NewReader("1.0"))
	}

	b.Log(r)
}

func BenchmarkDecimal1(b *testing.B) {
	var r interface{}

	for n := 0; n < b.N; n++ {
		r, _ = ParseRuneReader(Decimal, strings.NewReader("-123.456"))
	}

	b.Log(r)
}

func BenchmarkDecimal2(b *testing.B) {
	var r interface{}

	for n := 0; n < b.N; n++ {
		r, _ = ParseRuneReader(Decimal2, strings.NewReader("-123.456"))
	}

	b.Log(r)
}

// func ExampleRestarableOneOrMore() {
// 	parser := RestarableOneOrMore(
// 		StringifyResult(
// 			SeqOf(
// 				OneOrMore(ExpectAny([]rune("ab"))),
// 				SeqIgnore(AnyOf(Expect('_'), EOF)),
// 			),
// 		),
// 		Expect('_'),
// 	)

// 	r, err := ParseRuneReader(parser, strings.NewReader("aaaa_aaaaa_aaacccc_aaaaa_bbbbb"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// json, err := json.Marshal(r)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	log.Printf("%+v", r)
// 	// Output: ["aaaa", "aaaaa", &Partial{"aaa"}, "aaaaa", "aa"]
// }
