# go-parsing / parser
[![GoDoc](https://godoc.org/github.com/tekwizely/go-parsing/parser?status.svg)](https://godoc.org/github.com/tekwizely/go-parsing/parser)
[![MIT license](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/tekwizely/go-parsing/blob/master/LICENSE)

## Overview

Base components of a token analyzer, enabling the creation of hand-written parsers for generating Abstract Syntax Trees.

Some Features of this Parser:

* Infinite Lookahead
* Mark / Reset Functionality

## Using

#### Importing

```go
import "github.com/tekwizely/go-parsing/parser"
```

#### Initiating a Parser ( `parser.Parse` )

Parsing is initiated through the `Parse` method:

```go
// Parse initiates a parser against the input token stream.
//
func Parse(tokens token.Nexter, start ParserFn) ASTNexter
```

#### Parser Functions ( `parser.ParserFN` )

In addition to the `token.Nexter`, the Parse function also accepts a function which serves as the starting point for your parser.

The main Parser process will call into this `start` function to initiate parsing.

Parser functions scan tokens and emit Abstract Syntax Trees (ASTs).

Parser defines `ParserFn` with the following signature:

```go
// ParserFn are user functions that scan tokens and emit ASTs.
//
type ParserFn func(*Parser) ParserFn
```

#### Scanning Tokens ( `parser.Parser` )

When called, your parser function will receive a `Parser` object which provides methods to inspect tokens.

##### Peeking At Tokens ( `CanPeek()` / `Peek()` )

###### Before Peeking, Ensure That You Can

A well-behaved parser will first confirm if there are any tokens to review before trying to peek at or consume them.

For this, we have `CanPeek()`:

```go
// CanPeek confirms if the requested number of tokens are available in the peek buffer.
// n is 1-based.
// If CanPeek returns true, you can safely Peek for values up to, and including, n.
//
func (p *Parser) CanPeek(n int) bool
```

**NOTE:** When the Parser calls your parser function, it guarantees that `CanPeek(1) == true`, allowing you to review that first token without having to confirm its availability.

###### Taking A Peek

Once you're sure you can safely peek ahead, `Peek()` will let you review the token:

```go
// Peek allows you to look ahead at tokens without consuming them.
// n is 1-based.
//
func (p *Parser) Peek(n int) token.Token
```

##### Consuming Tokens ( `Next()` )

Once you confirm its safe to do so (see `CanPeek()` / `Peek()`), `Next()` will consume the next token from the input:

```go
// Next consumes and returns the next token in the input.
//
func (p *Parser) Next() token.Token
```

**NOTE:** When the Parser calls your parser function, it guarantees that `CanPeek(1) == true`, allowing you to consume that first token without having to confirm its availability.

##### Emitting ASTs ( `Emit()` )

Once you've processed one or more tokens, and built up an abstract syntax tree, you can emit it for further processing (for example, by an interpreter).

For this, we have `Emit()`:

```go
// Emit emits an AST.
// Consumed tokens are discarded.
//
func (p * Parser) Emit(ast interface{})
```

##### Discarding Consumed Tokens ( `Discard()` )

Sometimes, you may match a series of tokens that you simply wish to discard.

To discard consumed tokens without emitting an AST, use the `Discard()` method:

```go
// Discard discards the consumed tokens without emitting any ASTs.
//
func (p *Parser) Discard()
```

##### Creating Save Points ( `Marker()` / `CanReset()` / `Reset()` )

The Parser allows you to create save points and reset to them if you decide you want to re-try matching tokens in a different context.

###### Marking Your Spot

To create a save point, use the `Marker()` function:

```go
// Marker returns a marker that you can use to reset the parser to a previous state.
//
func (p * Parser) Marker() *Marker
```

###### Before Using A Marker, Ensure That You Can

A marker is good up until the next `Emit()` or `Discard()` action.

A well-behaved parser will first ensure that a marker is valid before trying to use it.

For this, we have `CanReset()`:

```go
// CanReset confirms if the marker is still valid.
// If CanReset returns true, you can safely reset the parser state to the marker position.
//
func (p * Parser) CanReset(m *Marker) bool
```

###### Resetting Parser State

Once you've confirmed a marker is still valid, `Reset()` will let you reset the parser state.

```go
// Reset resets the parser state to the marker position.
// Returns the ParserFn that was stored at the time the marker was created.
//
func (p * Parser) Reset(m *Marker) ParserFn
```

**NOTE:** Resetting a marker does not reset the parser function that was active when the marker was created.  Instead it returns the function reference, giving the current parser function the choice to use it or not.

#### Returning From Parser Function ( `return ParserFn` )

You'll notice that the `ParserFn` return type is another `ParserFn`

This is to allow for simplified flow control of your parser function.

###### One Pass

Your parser function only needs to concern itself with matching the very next tokens(s) of input.

This alleviates the need to manage complex looping / restart logic.

Simply return from your method after (possibly) emitting an AST, and the Parser will manage the looping.

###### Context-Switching

Switching contexts is as easy as returning a reference to another `ParserFn`.

###### Shutting Down The Parser Loop

You can shut down the main Parser loop from within your `ParserFn` by simply returning `nil`.

All previously emitted ASTs will still be available for pickup, but the parser will stop making any further `ParserFn` calls.

#### Retrieving Emitted ASTs ( `parser.ASTNexter` )

When called, the Parse function will return an `ASTNexter` which provides a means of retrieving ASTs emitted from the parser:

```go
type ASTNexter interface {

	// Next tries to fetch the next available AST, returning an error if something goes wrong.
	// Will return io.EOF to indicate end-of-file.
	//
	Next() (interface{}, error)
}
```

## Example (calculator)

Here's an example program that utilizes the parser (and lexer) to provide a simple calculator with support for variables.

**NOTE:** The source for this example can be found in the examples folder under `examples/calc/calc.go`

```go
package main

//
//	Input is read from STDIN
//
//	The input expression is matched against the following pattern:
//
//	input_exp:
//	( id '=' )? general_exp
//	general_exp:
//		operand ( operator operand )?
//	operand:
//		number | id | '(' general_exp ')'
//	operator:
//		'+' | '-' | '*' | '/'
//	number:
//		digit+ ( '.' digit+ )?
//	digit:
//		['0'..'9']
//	id:
//		alpha ( alpha | digit )*
//	alpha:
//		['a'..'z'] | ['A'..'Z']
//
//	Precedence is as expected, with '*' and '/' have higher precedence
//	than '+' and '-', as follows:
//
//	1 + 2 * 3 - 4 / 5  ==  1 + (2 * 3) - (4 / 5)
//

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/tekwizely/go-parsing/lexer"
	"github.com/tekwizely/go-parsing/lexer/token"
	"github.com/tekwizely/go-parsing/parser"
)

// We define our lexer tokens starting from the pre-defined EOF token
//
const (
	T_ID token.Type = lexer.T_START + iota
	T_NUMBER
	T_PLUS
	T_MINUS
	T_MULTIPLY
	T_DIVIDE
	T_EQUALS
	T_OPEN_PAREN
	T_CLOSE_PAREN
)

// To store variables
//
var vars = map[string]float64{}

// Single-character tokens
//
var singleChars = []byte{'+', '-', '*', '/', '=', '(', ')'}

var singleTokens = []token.Type{T_PLUS, T_MINUS, T_MULTIPLY, T_DIVIDE, T_EQUALS, T_OPEN_PAREN, T_CLOSE_PAREN}

// Whitespace
//
var bytesWhitespace = []byte{' ', '\t'}

// main
//
func main() {
	// Create a buffered reader from STDIN
	//
	stdin := bufio.NewReader(os.Stdin)

	// Read each line of input
	//
	for input, _, err := stdin.ReadLine(); err == nil; input, _, err = stdin.ReadLine() {
		// Anything to process?
		//
		if len(input) > 0 {
			// Create a new lexer to turn the input text into tokens
			//
			tokens := lexer.LexBytes(input, lex)

			// Create a new parser that feeds off the lexer and generates expression values
			//
			values := parser.Parse(tokens, parse)

			// Loop over parser emits
			//
			for value, parseErr := values.Next(); parseErr == nil; value, parseErr = values.Next() {
				fmt.Printf("%v\n", value)
			}
		}
	}
}

// lex is the starting (and only) StateFn for lexing the input into tokens
//
func lex(l *lexer.Lexer) lexer.LexerFn {

	// Single-char token?
	//
	if i := bytes.IndexRune(singleChars, l.Peek(1)); i >= 0 {
		l.Next()                    // Consuming the character
		l.EmitType(singleTokens[i]) // Emit just the type, discarding the consumed character
		return lex
	}

	switch {

	// Skip whitespace
	//
	case tryMatchWhitespace(l):
		l.DiscardToken()

	// Number
	//
	case tryMatchNumber(l):
		l.EmitToken(T_NUMBER)

	// ID
	//
	case tryMatchId(l):
		l.EmitToken(T_ID)

	// Unknown
	//
	default:
		r := l.Next()
		l.DiscardToken()
		fmt.Printf("Unknown Character: '%c'\n", r)
	}

	// See you again soon!
	return lex
}

// tryMatchWhitespace
//
func tryMatchWhitespace(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); r == ' ' || r == '\t' {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchRune
//
func tryMatchRune(l *lexer.Lexer, r rune) bool {
	if l.CanPeek(1) {
		if p := l.Peek(1); r == p {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchDigit
//
func tryMatchDigit(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); r >= '0' && r <= '9' {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchAlpha
//
func tryMatchAlpha(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchAlphaNum
//
func tryMatchAlphaNum(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchNumber [0-9]+ ( . [0-9]+ )?
//
func tryMatchNumber(l *lexer.Lexer) bool {
	if tryMatchDigit(l) {
		for tryMatchDigit(l) {
			// Nothing to do, rune already matched
		}
		if m := l.Marker(); tryMatchRune(l, '.') && tryMatchDigit(l) {
			for tryMatchDigit(l) {

			}
		} else {
			l.Reset(m)
		}
		return true
	}
	return false
}

// tryMatchId [a-zA-Z] [0-9a-zA-Z]*
//
func tryMatchId(l *lexer.Lexer) bool {
	if tryMatchAlpha(l) {
		for tryMatchAlphaNum(l) {
			// Nothing to do
		}
		return true
	}
	return false
}

// parse tries to parse an expression from the lexed tokens.
// Delegates to either parseEvaluation or parseAssignment.
//
func parse(p *parser.Parser) parser.ParserFn {

	switch {

	// Assignment
	//
	case p.CanPeek(3) && p.PeekType(1) == T_ID && p.PeekType(2) == T_EQUALS:
		return parseAssignment

	// Evaluation
	//
	default:
		return parseEvaluation
	}
}

// parseAssignment evaluates an expression and stores the results in the specified variable.
// The assignment will be in the form [ ID '=' expression ].
// Assumes "ID '='" has been peek-matched by root parser.
//
func parseAssignment(p *parser.Parser) parser.ParserFn {
	tId := p.Next()
	p.Next() // Skip '='
	if value, err := parseGeneralExpression(p); err == nil {
		// Should be at end of input
		//
		if !p.CanPeek(1) {
			vars[tId.Value()] = value
		} else {
			fmt.Println("Expecting Operator")
		}
	} else {
		fmt.Println(err.Error())
	}
	return nil // One pass
}

// parseEvaluation parses a general experssion and emits the computed result.
//
func parseEvaluation(p *parser.Parser) parser.ParserFn {
	if value, err := parseGeneralExpression(p); err == nil {
		// Should be at end of input
		//
		if !p.CanPeek(1) {
			p.Emit(value)
		} else {
			fmt.Println("Expecting Operator")
		}
	} else {
		fmt.Println(err.Error())
	}
	return nil // One pass
}

// parseGeneralExpression is the starting point for parsing a General Expression.
// It is basically a pass-through to parseAdditiveExpression, but it feels cleaner.
//
func parseGeneralExpression(p *parser.Parser) (f float64, err error) {
	return parseAdditiveExpression(p)
}

// parseAdditiveExpression parses [ expression ( ( '+' | '-' ) expression )? ].
//
func parseAdditiveExpression(p *parser.Parser) (f float64, err error) {

	var a float64
	if f, err = parseMultiplicitiveExpression(p); err == nil && p.CanPeek(1) {

		switch p.PeekType(1) {

		// Add (+)
		//
		case T_PLUS:
			p.Next() // Skip '+'
			if a, err = parseAdditiveExpression(p); err == nil {
				f += a
			}

		// Subtract (-)
		//
		case T_MINUS:
			p.Next() // Skip '-'
			if a, err = parseAdditiveExpression(p); err == nil {
				f -= a
			}
		}
	}

	return
}

// parseMultiplicitiveExpression parses [ expression ( ( '*' | '/' ) expression )? ].
//
func parseMultiplicitiveExpression(p *parser.Parser) (f float64, err error) {

	var m float64
	if f, err = parseOperand(p); err == nil && p.CanPeek(1) {

		switch p.PeekType(1) {

		// Multiply (*)
		//
		case T_MULTIPLY:
			p.Next() // Skip '*'
			if m, err = parseMultiplicitiveExpression(p); err == nil {
				f *= m
			}

		// Divide (/)
		//
		case T_DIVIDE:
			p.Next() // Skip '/'
			if m, err = parseMultiplicitiveExpression(p); err == nil {
				f /= m
			}
		}
	}

	return
}

// parseOperand parses [ id | number | '(' expression ')' ].
//
func parseOperand(p *parser.Parser) (f float64, err error) {

	// EOF
	//
	if !p.CanPeek(1) {
		return 0, errors.New("Unexpected EOF - Expecting operand")
	}

	m := p.Marker()

	switch p.PeekType(1) {

	// ID
	//
	case T_ID:
		var id = p.Next().Value()
		var ok bool
		if f, ok = vars[id]; !ok {
			err = errors.New(fmt.Sprintf("id '%s' not defined", id))
		}

	// Number
	//
	case T_NUMBER:
		n := p.Next().Value()
		if f, err = strconv.ParseFloat(n, 64); err != nil {
			fmt.Printf("Error parsing number '%s': %s", n, err.Error())
		}

	// '(' Expresson ')'
	//
	case T_OPEN_PAREN:
		p.Next() // Skip '('
		if f, err = parseGeneralExpression(p); err == nil {
			if p.CanPeek(1) && p.PeekType(1) == T_CLOSE_PAREN {
				p.Next() // Skip ')'
			} else {
				err = errors.New("Unbalanced Paren")
			}
		}

	// Unknown
	//
	default:
		err = errors.New("Expecting operand")
	}

	if err != nil {
		p.Reset(m)
	}

	return
}
```

## License

The `go-parsing` repo and all contained packages are released under the [MIT](https://opensource.org/licenses/MIT) License.  See `LICENSE` file.
