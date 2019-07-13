/*
Package token isolates the token-related types and interfaces used between the lexer and the parser.

*/
package token

// Token captures the type code, text string (optional), and positional
// information (optional) of tokens emitted from the lexer.
//
type Token interface {

	// Type returns the type code of the token.
	//
	Type() Type

	// Value returns the matched rune(s) that represent the token value.
	// Can be the empty string.
	//
	Value() string

	// Line returns the line number, within the source input, that the token originated on.
	// The definition of a 'line' is implementation-specific.
	// The field is optional.  If set by the token creator, it should be 1-based.
	// The value may only represent a best-guess.
	// A value less < 1 should be interpreted as not set for the token.
	//
	Line() int

	// Column returns the column number, relative to line(), that the token originated on.
	// The definition of a 'line' is implementation-specific.
	// Some implementations may track column from the beginning of the input (i.e file offset).
	// The value is generally expected to represent rune count (vs bytes).
	// The field is optional.  If set by the token creator, it should be 1-based.
	// The value may only represent a best-guess.
	// A value less < 1 should be interpreted as not set for the token.
	//
	Column() int
}

// Type identifies the type code of tokens emitted from the lexer.
//
type Type int

// Nexter provides a means of retrieving tokens (and errors) emitted from the lexer.
//
type Nexter interface {

	// Next tries to fetch the next available token, returning an error if something goes wrong.
	// Will return io.EOF to indicate end-of-file.
	// An error other than io.EOF may be recoverable and does not necessarily indicate end-of-file.
	// Even when an error is present, the returned token may still be valid and should be checked.
	// Once io.EOF is returned, any further calls will continue to return io.EOF.
	//
	Next() (Token, error)
}
