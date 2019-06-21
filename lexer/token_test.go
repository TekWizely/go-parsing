package lexer

import (
	"testing"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// assertToken
//
func assertToken(t *testing.T, tok *_token, typ token.Type, value string, eof bool) {
	if tok.typ != typ {
		t.Errorf("token.typ expecting '%d', received '%d'", typ, tok.typ)
	}
	if tok.value != value {
		t.Errorf("token.value expecting '%s', received '%s'", value, tok.value)
	}
	if tok.eof() != eof {
		t.Errorf("token.EOF() expecting '%t'", eof)
	}
}

// TestTokenEnums
//
func TestTokenEnums(t *testing.T) {
	// T_LEX_ERR
	//
	if T_LEX_ERR != 0 { // iota
		t.Error("T_LEX_ERR != 0")
	}
	// T_UNKNOWN
	//
	if T_UNKNOWN != 1 {
		t.Error("T_UNKNOWN != 1")
	}
	// T_EOF
	//
	if T_EOF != 2 {
		t.Error("T_EOF != 2")
	}
	// T_START
	//
	if T_START != 3 {
		t.Error("T_START != 3")
	}
	// t_end
	//
	if t_end != 4 {
		t.Error("t_end != 4, are there new tokens defined?")
	}
}

// TestNewToken
//
func TestNewToken(t *testing.T) {
	tok := newToken(T_START, "START")
	assertToken(t, tok, T_START, "START", false)
}

// TestNewTokenEmptyString
//
func TestNewTokenEmptyString(t *testing.T) {
	tok := newToken(T_START, "")
	assertToken(t, tok, T_START, "", false)
}

// TestNewTokenEOF
//
func TestNewTokenEOF(t *testing.T) {
	tok := newToken(T_EOF, "EOF")
	assertToken(t, tok, T_EOF, "EOF", true)
}

// TestNewTokenEOFEmptyString
//
func TestNewTokenEOFEmptyString(t *testing.T) {
	tok := newToken(T_EOF, "")
	assertToken(t, tok, T_EOF, "", true)
}
