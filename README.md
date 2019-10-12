# go-parsing [![GoDoc](https://godoc.org/github.com/tekwizely/go-parsing?status.svg)](https://godoc.org/github.com/tekwizely/go-parsing) [![MIT license](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/tekwizely/go-parsing/blob/master/LICENSE)

A Multi-Package Go Repo Focused on Text Parsing, with Lexers, Parsers, and Related Utils

## Goal

This repo aspires to be a useful toolset for creating hand-written lexers and parsers in Golang.

## Multi-Module Repo

The modules within this repo are intended to work together, but are allowed to evolve separately.

## Exported Modules

The following packages are currently exported:

- [github.com/tekwizely/go-parsing/lexer](https://godoc.org/github.com/tekwizely/go-parsing/lexer)
- [github.com/tekwizely/go-parsing/lexer/token](https://godoc.org/github.com/tekwizely/go-parsing/lexer/token)
- [github.com/tekwizely/go-parsing/parser](https://godoc.org/github.com/tekwizely/go-parsing/parser)

---------
### lexer ([github](https://github.com/TekWizely/go-parsing/tree/master/lexer) | [godoc](https://godoc.org/github.com/tekwizely/go-parsing/lexer))

Base components of a lexical analyzer, enabling the
creation of hand-written lexers for tokenizing textual content.

The tokenized data is suitable for processing with a parser.

Some Features of this Lexer:

* Rune-Centric
* Infinite Lookahead
* Mark / Reset Functionality

#### Example

See [go-parsing/lexer/examples/wordcount](https://github.com/TekWizely/go-parsing/tree/master/lexer/examples/wordcount) for an example program that utilizes the lexer.

-----------------
### lexer / token ( [github](https://github.com/TekWizely/go-parsing/tree/master/lexer/token) | [godoc](https://godoc.org/github.com/tekwizely/go-parsing/lexer))

Token-related types and interfaces used between the lexer and the parser.

----------
### parser ([github](https://github.com/TekWizely/go-parsing/tree/master/parser) | [godoc](https://godoc.org/github.com/tekwizely/go-parsing/parser))

Base components of a token analyzer, enabling the
creation of hand-written parsers for generating Abstract Syntax Trees.

Some Features of this Parser:

 * Infinite Lookahead
 * Mark / Reset Functionality

#### Example

See [go-parsing/parser/examples/calc](https://github.com/TekWizely/go-parsing/tree/master/parser/examples/calc) for an example program that utilizes the parser (and lexer).

----------
## License

The `tekwizely/go-parsing` repo and all contained packages are released under the [MIT](https://opensource.org/licenses/MIT) License.  See `LICENSE` file.
