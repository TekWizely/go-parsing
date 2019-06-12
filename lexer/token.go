package lexer

import "strconv"

// TokenType identifies the type of lex tokens.
//
type TokenType int

const (
	T_LEX_ERR TokenType = iota // Lexer error
	T_UNKNOWN                  // Unknown rune(s)
	T_EOF                      // EOF
	T_START                    // Marker for user tokens ( use T_START + iota )
	t_end                      // internal marker
)

// String
//
func (t TokenType) String() string {
	switch t {
	case T_LEX_ERR:
		return "T_LEX_ERR"
	case T_UNKNOWN:
		return "T_UNKNOWN"
	case T_EOF:
		return "T_EOF"
	default:
		return "token(" + strconv.Itoa(int(t)) + ")"
	}
}

// Token represents a token (with optional text string) returned from the lexer.
//
type Token struct {
	Type   TokenType
	String string
}

// newToken
//
func newToken(typ TokenType, str string) *Token {
	return &Token{Type: typ, String: str}
}

// eof returns true if the TokenType == T_EOF
//
func (t *Token) eof() bool { return T_EOF == t.Type }
