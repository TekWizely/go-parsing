package lexer

import (
	"testing"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// assertToken
//
func assertToken(t *testing.T, tok *_token, typ token.Type, value string, line int, column int, eof bool) {
	if tok.typ != typ {
		t.Errorf("token.typ expecting '%d', received '%d'", typ, tok.typ)
	}
	if tok.value != value {
		t.Errorf("token.value expecting '%s', received '%s'", value, tok.value)
	}
	if line >= 0 && tok.line != line {
		t.Errorf("token.line expecting '%d', received '%d'", line, tok.line)
	}
	if column >= 0 && tok.column != column {
		t.Errorf("token.column expecting '%d', received '%d'", column, tok.column)
	}
	if tok.eof() != eof {
		t.Errorf("token.EOF() expecting '%t'", eof)
	}
}

// TestTokenEnums
//
func TestTokenEnums(t *testing.T) {
	// TLexErr
	//
	if TLexErr != 0 { // iota
		t.Error("TLexErr != 0")
	}
	// TUnknown
	//
	if TUnknown != 1 {
		t.Error("TUnknown != 1")
	}
	// TEof
	//
	if TEof != 2 {
		t.Error("TEof != 2")
	}
	// TStart
	//
	if TStart != 3 {
		t.Error("TStart != 3")
	}
	// tEnd
	//
	if tEnd != 4 {
		t.Error("tEnd != 4, are there new tokens defined?")
	}
}

// TestNewToken
//
func TestNewToken(t *testing.T) {
	tok := newToken(TStart, "START", 10, 100)
	assertToken(t, tok, TStart, "START", 10, 100, false)
}

// TestNewTokenEmptyString
//
func TestNewTokenEmptyString(t *testing.T) {
	tok := newToken(TStart, "", 0, 0)
	assertToken(t, tok, TStart, "", 0, 0, false)
}

// TestNewTokenEOF
//
func TestNewTokenEOF(t *testing.T) {
	tok := newToken(TEof, "EOF", 0, 0)
	assertToken(t, tok, TEof, "EOF", 0, 0, true)
}

// TestNewTokenEOFEmptyString
//
func TestNewTokenEOFEmptyString(t *testing.T) {
	tok := newToken(TEof, "", 0, 0)
	assertToken(t, tok, TEof, "", 0, 0, true)
}
