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
	tok := newToken(TStart, "START")
	assertToken(t, tok, TStart, "START", false)
}

// TestNewTokenEmptyString
//
func TestNewTokenEmptyString(t *testing.T) {
	tok := newToken(TStart, "")
	assertToken(t, tok, TStart, "", false)
}

// TestNewTokenEOF
//
func TestNewTokenEOF(t *testing.T) {
	tok := newToken(TEof, "EOF")
	assertToken(t, tok, TEof, "EOF", true)
}

// TestNewTokenEOFEmptyString
//
func TestNewTokenEOFEmptyString(t *testing.T) {
	tok := newToken(TEof, "")
	assertToken(t, tok, TEof, "", true)
}
