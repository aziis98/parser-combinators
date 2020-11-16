# go `parser-combinators`

Personal parser combinator library for Golang.

## Examples

TODO

## Longer Examples

- [Minimark Syntax](/examples/minimark)

## TODO

 - **[Work in progress]** Recoverable parsing
 - **[Idea]** Add tab indented parser combinator inside the Minimark example syntax.
 - **[Idea]** LibConfig Syntax

## Commands

#### Benchmarks

```bash
go test -bench=. ./... > "notes/benchmark-$(git rev-parse --short HEAD).txt"
```
	