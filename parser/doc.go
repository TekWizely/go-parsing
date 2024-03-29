/*
Package parser implements the base components of a token analyzer, enabling the
creation of hand-written parsers for generating Abstract Syntax Trees.

Some Features of this Parser:

 * Infinite Lookahead
 * Mark / Reset Functionality


Initiating a Parser

	// Parse initiates a parser against the input token stream.
	//
	func Parse(tokens token.Nexter, start parser.Fn) ASTNexter


Parser Functions

In addition to the `token.Nexter`, the Parse function also accepts a function which serves as the starting point for
your parser:

	// parser.Fn are user functions that scan tokens and emit ASTs.
	//
	type parser.Fn func(*Parser) parser.Fn

The main Parser process will call into this function to initiate parsing.


Simplified Parser.Fn Loop

You'll notice that the `Parser.Fn` return type is another `Parser.Fn`.

This is to allow for simplified flow control of your parser function.

Your parser function only needs to concern itself with matching the very next tokens(s) of input.

This alleviates the need to manage complex looping / restart logic.

Simply return from your method after (possibly) emitting an AST, and the Parser will manage the looping.


Switching Parser Context

Switching contexts is as easy as returning a reference to another `Parser.Fn`.


Shutting Down The Parser

You can shut down the main Parser loop from within your `Parser.Fn` by simply returning `nil`.

All previously emitted ASTs will still be available for pickup, but the parser will stop making any further `Parser.Fn`
calls.


Scanning Tokens

Your Parser function receives a `*Parser` when called and can use the following methods to inspect and match tokens:

	// CanPeek confirms if the requested number of tokens are available in the peek buffer.
	//
	func (p *Parser) CanPeek(n int) bool

	// Peek allows you to look ahead at tokens without consuming them.
	//
	func (p *Parser) Peek(n int) token.Token

	// Next matches and returns the next token in the input.
	//
	func (p *Parser) Next() token.Token


Emitting ASTs

Once you've processed one or more tokens, and built up an abstract syntax tree, you can emit it for further processing
(for example, by an interpreter):

	// Emit emits an AST.
	//
	func (p * Parser) Emit(ast interface{})


Discarding Matched Tokens

Sometimes, you may match a series of tokens that you simply wish to discard:

	// Clear discards all previously-matched tokens without emitting any ASTs.
	//
	func (p *Parser) Clear()


Creating Save Points

The Parser allows you to create save points and reset to them if you decide you want to re-try matching tokens in a
different context:

	// Marker returns a marker that you can use to reset the parser to a previous state.
	//
	func (p * Parser) Marker() *Marker

A marker is good up until the next `Emit()` or `Clear()` action.

Before using a marker, confirm it is still valid:

	// Valid confirms if the marker is still valid.
	//
	func (m *Marker) Valid) bool

Once you've confirmed a marker is still valid:

	// Apply resets the parser state to the marker position.
	// Returns the Parser.Fn that was stored at the time the marker was created.
	//
	func (m *Marker) Apply() parser.Fn

NOTE: Resetting a marker does not reset the parser function that was active when the marker was created.
Instead it simply returns the function reference.  If you want to return control to the function saved in the marker,
you can use this pattern:

	return marker.Apply(); // Resets the parser and returns control to the saved Parser.Fn


Retrieving Emitted ASTs

When called, the `Parse` function will return an `ASTNexter` which provides a means of retrieving ASTs emitted from the
parser:

	type ASTNexter interface {

		// Next tries to fetch the next available AST, returning an error if something goes wrong.
		// Will return io.EOF to indicate end-of-file.
		//
		Next() (interface{}, error)
	}


Example Programs

See the `examples` folder for programs that demonstrate the parser (and lexer) functionality.

*/
package parser
