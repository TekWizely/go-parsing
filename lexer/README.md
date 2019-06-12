# go-parsing / lexer
[![GoDoc](https://godoc.org/github.com/tekwizely/go-parsing/lexer?status.svg)](https://godoc.org/github.com/tekwizely/go-parsing/lexer)
[![MIT license](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/TekWizely/go-parsing/blob/master/LICENSE)

## Overview

Base components of a lexical analyzer, enabling the
creation of hand-written lexers for tokenizing textual content.

The tokenized data is suitable for processing with a parser. 

Some Features of this Lexer:

* Rune-Centric
* Infinite Lookahead
* Mark / Reset Functionality

## Using

#### Importing

```go
import "github.com/tekwizely/go-parsing/lexer"
```

#### Initiating a Lexer ( `lexer.Lex*` )

Lexing is initiated through various `Lex*` methods, each accepting a different type of input to lex from:

###### Input Type: `string`

```go
func LexString(input string, start LexerFn) *Tokens
```

###### Input Type: `io.RuneReader`

```go
func LexRuneReader(input io.RuneReader, start LexerFn) *Tokens
```

###### Input Type: `io.Reader`

```go
func LexReader(input io.Reader, start LexerFn) *Tokens
```

###### Input Type: `[]rune`

```go
func LexRunes(input []rune, start LexerFn) *Tokens
```

###### Input Type: `[]byte`

```go
func LexBytes(input []byte, start LexerFn) *Tokens
```

#### Lexer Functions ( `lexer.LexerFN` )

In addition to the `input` data, each Lex function also accepts a function which serves as the starting point for your lexer.

The main Lexer process will call into this `start` function to initiate lexing.

Lexer functions scan runes and emit tokens.

Lexer defines `LexerFN` with the following signature:

```go
type LexerFn func(*Lexer) LexerFn
```

#### Scanning Input Data ( `lexer.Lexer` )

When called, your lexer function will receive a `Lexer` object which provides methods to inspect runes and match them to tokens.

##### Peeking At Runes ( `CanPeek()` / `Peek()` )

###### Before Peeking, Ensure That You Can

A well-behaved lexer will first confirm if there are any runes to review before trying to peek at them.

For this, we have `CanPeek()`:

```go
// CanPeek confirms if the requested number of runes are available in the peek buffer.
// n is 1-based.
// If CanPeek returns true, you can safely Peek for values up to, and including, n.
//
func (l *Lexer) CanPeek(n int) bool
```

**NOTE:** When the Lexer calls your lexer function, it guarantees that `CanPeek(1) == true`, allowing you to review that first rune without having to confirm its availability.

###### Taking A Peek

Once you're sure you can safely peek ahead, `Peek()` will let you review the rune:

```go
// Peek allows you to look ahead at runes without consuming them.
// n is 1-based.
//
func (l *Lexer) Peek(n int) rune
```

##### Consuming Runes ( `HasNext()` / `Next()` )

###### Before Consuming, Ensure That You Can

A well-behaved lexer will first confirm if there are any runes available before trying to consume them.

For this, we have `HasNext()`:

```go
// HasNext confirms if a rune is available to consume.
// If HasNext returns true, you can safely call Next to consume and return the rune.
//
func (l *Lexer) HasNext() bool
```

**NOTE:** When the Lexer calls your lexer function, it guarantees that `HasNext() == true`, allowing you to consume that first rune without having to confirm its availability.

**NOTE:** `HasNext()` is functionally equivalent to `CanPeek(1)` - So if you've already confirmed `CanPeek(n >= 1)` then you can safely forgo the `HasNext()` check. 

###### Consume The Rune

Once you confirm its safe to do so, `Next()` will consume the next rune from the input, making it part of the current token.

```go
// Next consumes and returns the next rune in the input.
//
func (l *Lexer) Next() rune
```

##### Reviewing The Current Token String ( `PeekToken()` / `PeekTokenRunes()` )

Once you've built up a token by consuming 1 or more runes, you may want to review it in its entirety before deciding what type of token it represents.

For this we have the following PeekToken functions:
 
###### Return Type: `string`

```go
// PeekToken allows you to inspect the currently matched rune sequence.
// The value is returned as a string, same as EmitToken() would provide.
//
func (l *Lexer) PeekToken() string
```

###### Return Type: `[]rune`

```go
// PeekTokenRunes allows you to inspect the currently matched rune sequence as a rune array ( []rune )
//
func (l *Lexer) PeekTokenRunes() []rune
```

##### Emitting Tokens ( `EmitToken()` / `EmitType()` )

Once you've determined what the consumed rune(s) represent, you can emit a token for further processing (for example, by a parser).

###### Emitting Token With Consumed Runes

Along with the token text, we need to specify the token Type.

The general method for this is `EmitToken()`:

```go
// EmitToken emits a token of the specified type, along with all of the consumed runes.
//
func (l *Lexer) EmitToken(t TokenType)
```

**NOTE:** See the section of the document regarding `"Token Types"` for details on defining tokens for your lexer.

###### Emitting Token Type Only

For some token types, the text value of the token isn't needed, and the `TokenType` carries enough context to fully describe the token (ex. `'+' -> T_PLUS`).

For these scenarios, you can use `EmitType` to emit just the token type, discarding the consumed runes:

```go
// EmitType emits a token of the specified type, discarding consumed runes.
//
func (l *Lexer) EmitType(t TokenType)
```

##### Discarding Consumed Runes ( `DiscardToken()` )

Sometimes, you may match a series of runes that you simply wish to discard. For example, in certain contexts, whitespace characters may be ignorable.

To discard consumed runes without emitting any tokens, use the `DiscardToken()` method:

```go
// DiscardToken discards the consumed runes without emitting any tokens.
//
func (l *Lexer) DiscardToken()
```

##### Creating Save Points ( `Marker()` / `CanReset()` / `Reset()` )

The Lexer allows you to create save points and reset to them if you decide you want to re-try matching runes in a different context.

###### Marking Your Spot

To create a save point, use the `Marker()` function:

```go
// Marker returns a marker that you can use to reset the lexer to a previous state.
//
func (l *Lexer) Marker() *Marker
```

###### Before Using A Marker, Ensure That You Can

A marker is good up until the next `Emit()` or `Discard()` action.

A well-behaved lexer will first ensure that a marker is valid before trying to use it.

For this, we have `CanReset()`:

```go
// CanReset confirms if the marker is still valid.
// If CanReset returns true, you can safely reset the lexer state to the marker position. 
//
func (l *Lexer) CanReset(m *Marker) bool
```

###### Resetting Lexer State

Once you've confirmed a marker is still valid, `Reset()` will let you reset the lexer state.

```go
// Reset resets the lexer state to the marker position.
// Returns the LexerFn that was stored at the time the marker was created.
//
func (l *Lexer) Reset(m *Marker) LexerFn
```

**NOTE:** Resetting a marker does not reset the lexer function that was active when the marker was created.  Instead it returns the function reference, giving the current lexer function the choice to use it or not.

#### Returning From Lexer Function ( `return LexerFn` )

You'll notice that the `LexerFN` return type is another `LexerFN`

This is to allow for simplified flow control of your lexer function.

###### One Pass

Your lexer function only needs to concern itself with matching the very next rune(s) of input.

This alleviates the need to manage complex looping / restart logic.

Simply return from your method after (possibly) emitting a token, and the Lexer will manage the looping.

###### Context-Switching

Switching contexts is as easy as returning a reference to another `LexerFn`.

For example, if, within your main lexer function, you encounter a `"`, you can simply return a reference to your `quotedStringLexer` function and the Lexer will transfer control to it.

Once finished, your quoted string lexer can return control back to your main lexer by returning a reference to your `mainLexer` function.

###### Shutting Down The Lexer Loop

You can shut down the main Lexer loop from within your `LexerFn` by simply returning `nil`.

All previously emitted tokens will still be available for pickup, but the lexer will stop making any further `LexerFn` calls.

#### Token Types ( `lexer.TokenType` )

##### Built-Ins

Lexer defines the `TokenType` type and a few pre-defined values:

```go
// TokenType identifies the type of lex tokens.
//
type TokenType int

const (
    T_LEX_ERR TokenType = iota // Lexer error
    T_UNKNOWN                  // Unknown rune(s)
    T_EOF                      // EOF
    T_START                    // Marker for user tokens ( use T_START + iota )
)
```

##### Defining Your Lexer Tokens

You define your own token types starting from `T_START`:

```go
const (
    T_INT = lexer.T_START + iota
    T_CHAR    
)
```

#### Retrieving Emitted Tokens ( `lexer.Tokens` )

When called, the Lex* functions will return a `Tokens` object which provides methods to retrieve tokens emitted from the lexer.

##### Token Iterator ( `HasNext()` / `Next()` )

Tokens implements a basic iterator pattern.

###### Before Retrieving, Ensure That You Can

A well-behaved program will first ensure that a token is available before trying to retrieve it.

For this, we have `HasNext()` : 

```go
// HasNext confirms if there are tokens available.
// If it returns true, you can safely call Next() to retrieve the next token.
//
func (t *Tokens) HasNext() bool
```
    
###### Retrieving A Token

Once you confirm its safe to do so, `Next()` will retrieve the next Token from the lexer output.

```go
// Next Retrieves the next token from the lexer.
//
func (t *Tokens) Next() *Token
```

## Example (wordcount)

Here's an example program that utilizes the lexer to count the number of words, spaces, lines and characters in a file.

**NOTE:** The source for this example can be found in the examples folder under `examples/wordcount/wordcount.go`

```go
package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/tekwizely/go-parsing/lexer"
)

// Usage : wordcount <filename>
//
func usage() {
	fmt.Printf("usage: %s <filename>\n", os.Args[0])
}

// We define our lexer tokens starting from the pre-defined START token
//
const (
	T_SPACE = lexer.T_START + iota
	T_NEWLINE
	T_WORD
)

// We will attempt to match 3 newline styles: [ "\n", "\r", "\r\n" ]
//
const (
	runeNewLine = '\n'
	runeReturn  = '\r'
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	var (
		file *os.File
		err  error
	)

	//  Open the file, panic on error
	//
	if file, err = os.Open(os.Args[1]); err != nil {
		panic(err)
	}

	var (
		chars  int = 0
		words  int = 0
		spaces int = 0
		lines  int = 0
	)

	// To help us track last line in file (which may not have a newline)
	//
	var emptyLine bool = true

	tokens := lexer.LexReader(file, lexerFn)

	// Process lexer-emitted tokens
	//
	for tokens.HasNext() {
		t := tokens.Next()
		chars += len(t.String)

		switch t.Type {
		case T_WORD:
			words++
			emptyLine = false

		case T_NEWLINE:
			lines++
			spaces += len(t.String)
			emptyLine = true

		case T_SPACE:
			spaces += len(t.String)
			emptyLine = false

		default:
			panic("Unreachable")
		}
	}

	// If last line not empty, up line count
	//
	if !emptyLine {
		lines++
	}

	fmt.Printf("%d words, %d spaces, %d lines, %d chars\n", words, spaces, lines, chars)
}

func lexerFn(l *lexer.Lexer) lexer.LexerFn {

	// Can skip canPeek() check on first rune, per lexer rules
	//
	switch r := l.Peek(1); {

	// Newline '\n'
	// We check this before Space to avoid hit from unicode.IsSpace() check
	//
	case r == runeNewLine:
		l.Next()
		l.EmitToken(T_NEWLINE)

	// Return '\r', optionally followed by newLine '\n'
	// We check this before Space to avoid hit from unicode.IsSpace() check
	//
	case r == runeReturn:
		l.Next()
		if l.CanPeek(1) && l.Peek(1) == runeNewLine {
			l.Next()
		}
		l.EmitToken(T_NEWLINE)

	// Space or Word
	//
	default:
		isSpace := unicode.IsSpace(r)
		// Consume verified rune to avoid re-check
		//
		l.Next()
		// Consume further consecutive runes of same type
		//
		for l.CanPeek(1) && unicode.IsSpace(l.Peek(1)) == isSpace {
			l.Next()
		}
		// Emit token
		//
		if isSpace {
			l.EmitToken(T_SPACE)
		} else {
			l.EmitToken(T_WORD)
		}
	}

	return lexerFn // Let's do it again
}
```

## License

The `go-parsing` repo and all contained packages are released under the [MIT](https://opensource.org/licenses/MIT) License.  See `LICENSE` file.
