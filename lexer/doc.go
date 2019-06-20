/*
Package lexer implements the base components of a lexical analyzer, enabling the
creation of hand-written lexers for tokenizing textual content.

Some Features of this Lexer:

 * Rune-Centric
 * Infinite Lookahead
 * Mark / Reset Functionality


Initiating a Lexer

Lexing is initiated through various Lex* methods, each accepting a different type of input to lex from:

	// Input Type: string
	//
	func LexString(input string, start LexerFn) TokenNexter

	// Input Type: io.RuneReader
	//
	func LexRuneReader(input io.RuneReader, start LexerFn) TokenNexter

	// Input Type: io.Reader
	//
	func LexReader(input io.Reader, start LexerFn) TokenNexter

	// Input Type: []rune
	//
	func LexRunes(input []rune, start LexerFn) TokenNexter

	// Input Type: []byte
	//
	func LexBytes(input []byte, start LexerFn) TokenNexter


Lexer Functions

In addition to the input data, each Lex function also accepts a function which serves as the starting point for your
lexer:

	// LexerFn are user functions that scan runes and emit tokens.
	//
	type LexerFn func(*Lexer) LexerFn

The main Lexer process will call into this function to initiate lexing.


Simplified LexerFn Loop

You'll notice that the `LexerFN` return type is another `LexerFn`.

This is to allow for simplified flow control of your lexer function.

Your lexer function only needs to concern itself with matching the very next rune(s) of input.

This alleviates the need to manage complex looping / restart logic.

Simply return from your method after (possibly) emitting a token, and the Lexer will manage the looping.


Switching Lexer Context

Switching contexts is as easy as returning a reference to another LexerFn.

For example, if, within your main lexer function, you encounter a `"`, you can simply return a reference to your
`quotedStringLexer` function and the Lexer will transfer control to it.

Once finished, your quoted string lexer can return control back to your main lexer by returning a reference to your
`mainLexer` function.


Shutting Down The Lexer

You can shut down the main Lexer loop from within your `LexerFn` by simply returning `nil`.

All previously emitted tokens will still be available for pickup, but the lexer will stop making any further `LexerFn`
calls.


Scanning Runes

Your Lexer function receives a `*Lexer` when called and can use the following methods to inspect and consume runes:

	// CanPeek confirms if the requested number of runes are available in the peek buffer.
	//
	func (l *Lexer) CanPeek(n int) bool

	// Peek allows you to look ahead at runes without consuming them.
	//
	func (l *Lexer) Peek(n int) rune

	// HasNext confirms if a rune is available to consume.
	//
	func (l *Lexer) HasNext() bool

	// Next consumes and returns the next rune in the input.
	//
	func (l *Lexer) Next() rune

	// PeekToken allows you to inspect the currently matched rune sequence.
	//
	func (l *Lexer) PeekToken() string


Emitting Tokens

Once you've determined what the consumed rune(s) represent, you can emit a token for further processing
(for example, by a parser):

	// EmitToken emits a token of the specified type, along with all of the consumed runes.
	//
	func (l *Lexer) EmitToken(t TokenType)

	// EmitType emits a token of the specified type, discarding consumed runes.
	//
	func (l *Lexer) EmitType(t TokenType)

NOTE: See the section of the document regarding "Token Types" for details on defining tokens for your lexer.


Discarding Matched Runes

Sometimes, you may match a series of runes that you simply wish to discard:

	// DiscardToken discards the consumed runes without emitting any tokens.
	//
	func (l *Lexer) DiscardToken()


Creating Save Points

The Lexer allows you to create save points and reset to them if you decide you want to re-try matching runes in a
different context:

	// Marker returns a marker that you can use to reset the lexer to a previous state.
	//
	func (l *Lexer) Marker() *Marker

A marker is good up until the next `Emit()` or `Discard()` action.

Before using a marker, confirm it is still valid:

	// CanReset confirms if the marker is still valid.
	//
	func (l *Lexer) CanReset(m *Marker) bool

Once you've confirmed a marker is still valid:

	// Reset resets the lexer state to the marker position.
	// Returns the LexerFn that was stored at the time the marker was created.
	//
	func (l *Lexer) Reset(m *Marker) LexerFn

NOTE: Resetting a marker does not reset the lexer function that was active when the marker was created.
Instead it simply returns the function reference.  If you want to return control to the function saved in the marker,
you can use this pattern:

	return lexer.Reset(marker); // Resets the lexer and returns control to the saved LexerFn


Token Types

Lexer defines the TokenType type and a few pre-defined values:

	// TokenType identifies the type of lex tokens.
	//
	type TokenType int

	const (
		T_LEX_ERR TokenType = iota // Lexer error
		T_UNKNOWN                  // Unknown rune(s)
		T_EOF                      // EOF
		T_START                    // Marker for user tokens ( use T_START + iota )
	)

You define your own token types starting from T_START:

	const (
		T_INT = lexer.T_START + iota
		T_CHAR
	)


Retrieving Emitted Tokens

When called, the `Lex*` functions will return a `TokenNexter` which provides methods to retrieve tokens emitted from the
lexer.

`TokenNexter` implements a basic iterator pattern:

	type TokenNexter interface {

		// HasNext confirms if there are tokens available.
		//
		HasNext() bool

		// Next Retrieves the next token from the lexer.
		//
		Next() Token
	}


Example Programs

See the `examples` folder for programs that demonstrate the lexer functionality.

*/
package lexer
