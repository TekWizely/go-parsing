package lexer

import "github.com/tekwizely/go-parsing/lexer/token"

const (
	// TLexErr represents a Lexer error
	//
	TLexErr token.Type = iota
	// TUnknown represents Unknown rune(s)
	//
	TUnknown
	// TEof represents end of file
	//
	TEof
	// TStart is a marker for user tokens ( use TStart + iota )
	//
	TStart
	// tEnd is an internal marker
	//
	tEnd
)

// token is the internal structure that backs the lexer's Token.
//
type _token struct {
	typ   token.Type
	value string
}

// newToken
//
func newToken(typ token.Type, value string) *_token {
	return &_token{typ: typ, value: value}
}

// Type implements Token.Type().
//
func (t *_token) Type() token.Type {
	return t.typ
}

// Value implements Token.Value().
//
func (t *_token) Value() string {
	return t.value
}

// eof returns true if the token.Type == TEof.
//
func (t *_token) eof() bool { return TEof == t.typ }
