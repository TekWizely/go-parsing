/*
Package parsing is a multi-package Go repo focused on text parsing, with lexers, parsers, and related utils.

Goal

This repo aspires to be the best toolset for creating hand-written lexers and parsers in Golang.


Multi-Module Repo

The modules within this repo are intended to work together, but are allowed to evolve separately.


Exported Modules

The following packages are currently exported:

 * github.com/tekwizely/go-parsing/lexer
 * github.com/tekwizely/go-parsing/parser


Lexer

Base components of a lexical analyzer, enabling the
creation of hand-written lexers for tokenizing textual content.

The tokenized data is suitable for processing with a parser.

Some Features of this Lexer:

 * Rune-Centric
 * Infinite Lookahead
 * Mark / Reset Functionality


Parser

Base components of a token analyzer, enabling the
creation of hand-written parsers for generating Abstract Syntax Trees.

Some Features of this Parser:

 * Infinite Lookahead
 * Mark / Reset Functionality


Links

You can learn more online:

  * GitHub https://github.com/TekWizely/go-parsing
  * GoDoc  https://godoc.org/github.com/tekwizely/go-parsing


NOTE

Although useful in its own right, this file (doc.go) mostly exists to prevent pre-commit hooks from generating
"no file" errors against the root folder.  See:

      https://github.com/TekWizely/go-parsing/issues/3


License

The go-parsing repo and all contained packages are released under the MIT License. See LICENSE file.

*/
package parsing
