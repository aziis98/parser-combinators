# Longer IDEAs

## Parser Combinators as Structs

To enhance performance a new Go implementation could be like

```go
type Parser interface {
    Apply(context ParseContext) (interface{}, error)
}

type ParseContext interface {
    
    // Operations *on* the cursor

    // pushes the current position
    Begin()
    // restores previous position
    Break()
    // pops the new position and sets its as new current
    End()

    // Operations *at* the cursor (~)

    // retrive current character at scanner position
    PeekRune() rune
    NextRune() rune
}

type Expect struct {
    Expected rune
}

func (p *Expect) Apply(context ParseContext) (interface{}, error) {
    next := context.NextRune()
    
    if p.Expected != next {
        return nil, next
    }

    return p.Expected, nil
}

type Seq struct {
    Parsers []Parser
}

func (p *Seq) Apply(context ParseContext) (interface{}, error) {
    results := []interface{}{}
    
    context.Begin()
    for _, parser := range p.Parsers {
        r, err := parser.Apply(context)
        if err != nil {
            context.Break()
            return nil, err
        }

        results = append(results, r)
    }
    context.End()

    return results, nil
}

type Any struct {
    Parsers []Parser
}

type Transform struct {
    Parser      Parser
    Transform   func(interface{}) interface{}
}

func (p *Any) Apply(context ParseContext) (interface{}, error) {
    errs := []error{}
    
    for _, parser := range p.Parsers {
        r, err := parser.Apply(context)
        if err != nil {
            errs = append(errs, err)
        }

        if err == nil {
            return r, nil
        }
    }

    return nil, NewCompoundError(errs)
}

type HTMLParser struct {
    // ...
}

type MarkdownParser struct {
    // ...
}

// ...

func Example1() {
    parser := &Any{ 
        &Seq{ &Expect{ 'a' }, &Expect{ 'a' } }, 
        &Seq{ &Expect{ 'b' }, &Expect{ 'b' } },
    }

    // r, err := Parse(parsers, reader)
    // reader ~~> "aaa"   =>    r = ['a', 'a']
    // reader ~~> "bbb"   =>    r = ['b', 'b']
    // reader ~~> "b"     =>    err = "Expected 'b'"
    // reader ~~> "a"     =>    err = "Expected 'a'"
    // reader ~~> ""      =>    err = "Expected 'a' or expected 'b'"
}

```

I think that definitions look similar enough to a parser combinator description and this has the up side that there aren't much allocations as the parse procedes (before the instances of "ParseState" were created for each return of a "Parser" function), all allocations are located in a single instance of the stacked-scanner that records cursor positions.

More, some grammars like `[a-z]+` could be parsed even without more than one level of the scanner stack.
