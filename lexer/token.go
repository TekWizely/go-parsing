package lexer

import "github.com/tekwizely/go-parsing/lexer/token"

const (
	// T_LEX_ERR represents a Lexer error
	//
	T_LEX_ERR token.Type = iota
	// T_UNKNOWN represents Unknown rune(s)
	//
	T_UNKNOWN
	// T_EOF represents end of file
	//
	T_EOF
	// T_START is a marker for user tokens ( use T_START + iota )
	//
	T_START
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

// eof returns true if the token.Type == T_EOF.
//
func (t *_token) eof() bool { return T_EOF == t.typ }
