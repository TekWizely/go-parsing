/*
Package token isolates the token-related types and interfaces used between the lexer and the parser.

*/
package token

// Token captures the type code + optional text string emitted from the lexer.
//
type Token interface {

	// Type returns the type code of the token.
	//
	Type() Type

	// Value returns the matched rune(s) that represent the token value.
	// Can be the empty string.
	//
	Value() string
}

// Type identifies the type code of tokens emitted from the lexer.
//
type Type int

// Nexter provides methods to retrieve tokens emitted from the lexer.
// Implements a basic iterator pattern with HasNext() and Next() methods.
//
type Nexter interface {

	// HasNext confirms if there are tokens available.
	// If it returns true, you can safely call Next() to retrieve the next token.
	// If it returns false, EOF has been reached and calling Next() will generate a panic.
	//
	HasNext() bool

	// Next Retrieves the next token from the lexer.
	// See HasNext() to determine if any tokens are available.
	// Panics if HasNext() returns false.
	//
	Next() Token
}
