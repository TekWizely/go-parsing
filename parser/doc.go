/*
Package parser implements the base components of a token analyzer, enabling the
creation of hand-written parsers for generating Abstract Syntax Trees.

Some Features of this Lexer:

 * Infinite Lookahead
 * Mark / Reset Functionality


Initiating a Parser

	// Parse initiates a parser against the input token stream.
	//
	func Parse(tokens *lexer.Tokens, start ParserFn) *Emits


Parser Functions

In addition to the `Tokens` supplier, the Parse function also accepts a function which serves as the starting point for
your parser:

	// ParserFn are user functions that scan tokens and emit ASTs.
	//
	type ParserFn func(*Parser) ParserFn

The main Parser process will call into this function to initiate parsing.


Simplified ParserFn Loop

You'll notice that the `ParserFn` return type is another `ParserFn`.

This is to allow for simplified flow control of your parser function.

Your parser function only needs to concern itself with matching the very next tokens(s) of input.

This alleviates the need to manage complex looping / restart logic.

Simply return from your method after (possibly) emitting an AST, and the Parser will manage the looping.


Switching Parser Context

Switching contexts is as easy as returning a reference to another `ParserFn`.


Shutting Down The Parser

You can shut down the main Parser loop from within your `ParserFn` by simply returning `nil`.

All previously emitted ASTs will still be available for pickup, but the parser will stop making any further `ParserFn`
calls.


Scanning Tokens

Your Parser function receives a `*Parser` when called and can use the following methods to inspect and consume tokens:

	// CanPeek confirms if the requested number of tokens are available in the peek buffer.
	//
	func (p *Parser) CanPeek(n int) bool

	// Peek allows you to look ahead at tokens without consuming them.
	//
	func (p *Parser) Peek(n int) *lexer.Token

	// HasNext confirms if a token is available to consume.
	//
	func (p *Parser) HasNext() bool

	// Next consumes and returns the next token in the input.
	//
	func (p *Parser) Next() *lexer.Token


Emitting ASTs

Once you've processed one or more tokens, and built up an abstract syntax tree, you can emit it for further processing
(for example, by an interpreter):

	// Emit emits an AST.
	//
	func (p * Parser) Emit(ast interface{})


Discarding Matched Tokens

Sometimes, you may match a series of tokens that you simply wish to discard:

	// Discard discards the consumed tokens without emitting any ASTs.
	//
	func (p *Parser) Discard()


Creating Save Points

The Parser allows you to create save points and reset to them if you decide you want to re-try matching tokens in a
different context:

	// Marker returns a marker that you can use to reset the parser to a previous state.
	//
	func (p * Parser) Marker() *Marker

A marker is good up until the next `Emit()` or `Discard()` action.

Before using a marker, confirm it is still valid:

	// CanReset confirms if the marker is still valid.
	//
	func (p * Parser) CanReset(m *Marker) bool

Once you've confirmed a marker is still valid:

	// Reset resets the parser state to the marker position.
	// Returns the ParserFn that was stored at the time the marker was created.
	//
	func (p * Parser) Reset(m *Marker) ParserFn

NOTE: Resetting a marker does not reset the parser function that was active when the marker was created.
Instead it simply returns the function reference.  If you want to return control to the function saved in the marker,
you can use this pattern:

	return parser.Reset(marker); // Resets the parser and returns control to the saved ParserFn


Retrieving Emitted ASTs

When called, the `Parse` function will return an `Emits` object which provides methods to retrieve ASTs emitted from the
parser.

`Emits` implements a basic iterator pattern:

	// HasNext confirms if there are ASTs available.
	//
	func (e *Emits) HasNext() bool

	// Next Retrieves the next AST from the parser.
	//
	func (e *Emits) Next() interface{}


Example Programs

See the `examples` folder for programs that demonstrate the parser (and lexer) functionality.

*/
package parser