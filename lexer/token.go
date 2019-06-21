package lexer

import "github.com/tekwizely/go-parsing/lexer/token"

const (
	T_LEX_ERR token.Type = iota // Lexer error
	T_UNKNOWN                   // Unknown rune(s)
	T_EOF                       // EOF
	T_START                     // Marker for user tokens ( use T_START + iota )
	t_end                       // internal marker
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
