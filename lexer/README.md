# go-parsing / lexer [![GoDoc](https://godoc.org/github.com/tekwizely/go-parsing/lexer?status.svg)](https://godoc.org/github.com/tekwizely/go-parsing/lexer) [![MIT license](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/tekwizely/go-parsing/blob/master/LICENSE)

## Overview

Base components of a lexical analyzer, enabling the creation of hand-written lexers for tokenizing textual content.

The tokenized data is suitable for processing with a parser.

Some Features of this Lexer:

 * Rune-Centric
 * Infinite Lookahead
 * Mark / Reset Functionality
 * Line / Column Tracking


## Using

#### Importing

```go
import "github.com/tekwizely/go-parsing/lexer"
```

-----------------------
#### Initiating a Lexer ( `lexer.Lex*` )

Lexing is initiated through various `Lex*` methods, each accepting a different type of input to lex from:

###### Input Type: `string`

```go
func LexString(input string, start lexer.Fn) token.Nexter
```

###### Input Type: `io.RuneReader`

```go
func LexRuneReader(input io.RuneReader, start lexer.Fn) token.Nexter
```

###### Input Type: `io.Reader`

```go
func LexReader(input io.Reader, start lexer.Fn) token.Nexter
```

###### Input Type: `[]rune`

```go
func LexRunes(input []rune, start lexer.Fn) token.Nexter
```

###### Input Type: `[]byte`

```go
func LexBytes(input []byte, start lexer.Fn) token.Nexter
```

--------------------
#### Lexer Functions ( `lexer.Fn` )

In addition to the `input` data, each Lex function also accepts a function which serves as the starting point for your lexer.

The main Lexer process will call into this `start` function to initiate lexing.

Lexer functions scan runes and emit tokens.

Lexer defines `Lexer.Fn` with the following signature:

```go
// lexer.Fn are user functions that scan runes and emit tokens.
//
type lexer.Fn func(*Lexer) lexer.Fn
```

-------------------
#### Scanning Runes ( `lexer.Lexer` )

When called, your lexer function will receive a `Lexer` object which provides methods to inspect runes and match them to tokens.

##### Peeking At Runes ( `CanPeek()` / `Peek()` )

###### Before Peeking, Ensure That You Can

A well-behaved lexer will first confirm if there are any runes to review before trying to peek at or match them.

For this, we have `CanPeek()`:

```go
// CanPeek confirms if the requested number of runes are available in the peek buffer.
// n is 1-based.
// If CanPeek returns true, you can safely Peek for values up to, and including, n.
//
func (l *Lexer) CanPeek(n int) bool
```

**NOTE:** When the Lexer calls your lexer function, it guarantees that `CanPeek(1) == true`, ensuring there is at least one rune to review/match.

###### Taking A Peek

Once you're sure you can safely peek ahead, `Peek()` will let you review the rune:

```go
// Peek allows you to look ahead at runes without consuming them.
// n is 1-based.
//
func (l *Lexer) Peek(n int) rune
```

---------------------
##### Consuming Runes ( `Next()` )

Once you confirm its safe to do so (see `CanPeek()` / `Peek()`), `Next()` will match the next rune from the input, making it part of the current token:

```go
// Next matches and returns the next rune in the input.
//
func (l *Lexer) Next() rune
```

**NOTE:** When the Lexer calls your lexer function, it guarantees that `CanPeek(1) == true`, ensuring there is at least one rune to review/match.

----------------------------------------
##### Reviewing The Current Token String ( `PeekToken()` )

Once you've built up a token by consuming 1 or more runes, you may want to review it in its entirety before deciding what type of token it represents.

For this we have `PeekToken()`:

```go
// PeekToken allows you to inspect the currently matched rune sequence.
// The value is returned as a string, same as EmitToken() would provide.
//
func (l *Lexer) PeekToken() string
```

---------------------
##### Emitting Tokens ( `EmitToken()` / `EmitType()` )

Once you've determined what the matched rune(s) represent, you can emit a token for further processing (for example, by a parser).

###### Emitting Token With Matched Runes

Along with the token text, we need to specify the token Type.

The general method for this is `EmitToken()`:

```go
// EmitToken emits a token of the specified type, along with all of the matched runes.
//
func (l *Lexer) EmitToken(t token.Type)
```

**NOTE:** See the section of the document regarding `"Token Types"` for details on defining tokens for your lexer.

###### Emitting Token Type Only

For some token types, the text value of the token isn't needed, and the `token.Type` carries enough context to fully describe the token (ex. `'+' -> TPlus`).

For these scenarios, you can use `EmitType` to emit just the token type, discarding the previously-matched runes:

```go
// EmitType emits a token of the specified type, discarding all previously-matched runes.
//
func (l *Lexer) EmitType(t token.Type)
```

------------------------------
##### Discarding Matched Runes ( `Clear()` )

Sometimes, you may match a series of runes that you simply wish to discard. For example, in certain contexts, whitespace characters may be ignorable.

To discard previously-matched runes without emitting any tokens, use the `Clear()` method:

```go
// Clear discards all previously-matched runes without emitting any tokens.
//
func (l *Lexer) Clear()
```

--------------------------
##### Creating Save Points ( `Marker()` / `Valid()` / `Apply()` )

The Lexer allows you to create save points and reset to them if you decide you want to re-try matching runes in a different context.

###### Marking Your Spot

To create a save point, use the `Marker()` function:

```go
// Marker returns a marker that you can use to reset the lexer to a previous state.
//
func (l *Lexer) Marker() *Marker
```

###### Before Using A Marker, Ensure That You Can

A marker is good up until the next `Emit()` or `Clear()` action.

A well-behaved lexer will first ensure that a marker is valid before trying to use it.

For this, we have `Marker.Valid()`:

```go
// Valid confirms if the marker is still valid.
// If Valid returns true, you can safely reset the lexer state to the marker position.
//
func (m *Marker) Valid() bool
```

###### Resetting Lexer State

Once you've confirmed a marker is still valid, `Marker.Apply()` will let you reset the lexer state.

```go
// Apply resets the lexer state to the marker position.
// Returns the Lexer.Fn that was stored at the time the marker was created.
//
func (m *Marker) Apply() lexer.Fn
```

**NOTE:** Resetting a marker does not reset the lexer function that was active when the marker was created.  Instead it returns the function reference, giving the current lexer function the choice to use it or not.

----------------------------------
#### Returning From Lexer Function ( `return lexer.Fn` )

You'll notice that the `Lexer.Fn` return type is another `Lexer.Fn`

This is to allow for simplified flow control of your lexer function.

###### One Pass

Your lexer function only needs to concern itself with matching the very next rune(s) of input.

This alleviates the need to manage complex looping / restart logic.

Simply return from your method after (possibly) emitting a token, and the Lexer will manage the looping.

###### Context-Switching

Switching contexts is as easy as returning a reference to another `Lexer.Fn`.

For example, if, within your main lexer function, you encounter a `"`, you can simply return a reference to your `quotedStringLexer` function and the Lexer will transfer control to it.

Once finished, your quoted string lexer can return control back to your main lexer by returning a reference to your `mainLexer` function.

###### Shutting Down The Lexer Loop

You can shut down the main Lexer loop from within your `Lexer.Fn` by simply returning `nil`.

All previously emitted tokens will still be available for pickup, but the lexer will stop making any further `Lexer.Fn` calls.

----------------
#### Token Types ( `token.Type` )

##### Built-Ins

Lexer defines a few pre-defined token values:

```go

const (
    TLexErr token.Type = iota // Lexer error
    TUnknown                  // Unknown rune(s)
    TEof                      // EOF
    TStart                    // Marker for user tokens ( use TStart + iota )
)
```

##### Defining Your Lexer Tokens

You define your own token types starting from `TStart`:

```go
const (
    TInt = lexer.TStart + iota
    TChar
)
```

------------------------------
#### Retrieving Emitted Tokens ( `token.Nexter` )

When called, the Lex* functions will return a `token.Nexter` which provides a means of retrieving tokens (and errors) emitted from the lexer:

```go
type Nexter interface {

	// Next tries to fetch the next available token, returning an error if something goes wrong.
	// Will return io.EOF to indicate end-of-file.
	//
	Next() (Token, error)
}
```

-------------------------------
#### Tracking Lines and Columns ( `Token.Line()` / `Token.Column()` )

Lexer tracks lines and columns as runes are consumed, and exposes them in the emitted Tokens.

Lexer uses `'\n'` as the newline separator when tracking line counts.

**NOTE:** Error messages with line/column information may reference the start of an attempted token match and not the position of the rune(s) that generated the error.

----------
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
	TSpace = lexer.TStart + iota
	TNewline
	TWord
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

	// int inits to 0
	//
	var (
		chars  int
		words  int
		spaces int
		lines  int
	)

	// To help us track last line in file (which may not have a newline)
	//
	var emptyLine = true

	tokens := lexer.LexReader(file, lexerFn)

	// Process lexer-emitted tokens
	//
	for t, lexErr := tokens.Next(); lexErr == nil; t, lexErr = tokens.Next() { // We only emit EOF so !nil should do it
		chars += len(t.Value())

		switch t.Type() {
		case TWord:
			words++
			emptyLine = false

		case TNewline:
			lines++
			spaces += len(t.Value())
			emptyLine = true

		case TSpace:
			spaces += len(t.Value())
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

func lexerFn(l *lexer.Lexer) lexer.Fn {

	// Can skip canPeek() check on first rune, per lexer rules
	//
	switch r := l.Peek(1); {

	// Newline '\n'
	// We check this before Space to avoid hit from unicode.IsSpace() check
	//
	case r == runeNewLine:
		l.Next()
		l.EmitToken(TNewline)

	// Return '\r', optionally followed by newLine '\n'
	// We check this before Space to avoid hit from unicode.IsSpace() check
	//
	case r == runeReturn:
		l.Next()
		if l.CanPeek(1) && l.Peek(1) == runeNewLine {
			l.Next()
		}
		l.EmitToken(TNewline)

	// Space or Word
	//
	default:
		isSpace := unicode.IsSpace(r)
		// Match verified rune to avoid re-check
		//
		l.Next()
		// Match further consecutive runes of same type
		//
		for l.CanPeek(1) && unicode.IsSpace(l.Peek(1)) == isSpace {
			l.Next()
		}
		// Emit token
		//
		if isSpace {
			l.EmitToken(TSpace)
		} else {
			l.EmitToken(TWord)
		}
	}

	return lexerFn // Let's do it again
}
```

----------
## License

The `tekwizely/go-parsing` repo and all contained packages are released under the [MIT](https://opensource.org/licenses/MIT) License.  See `LICENSE` file.
