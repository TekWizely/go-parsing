package lexer

import (
	"testing"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// expectNexterHasNext
//
func expectNexterHasNext(t *testing.T, nexter token.Nexter, match bool) {
	if nexter.HasNext() != match {
		t.Errorf("Nexter.HasNext() expecting '%t'", match)
	}
}

// expectNexterNext
//
func expectNexterNext(t *testing.T, nexter token.Nexter, typ token.Type, value string) {
	tok := nexter.Next()
	if tok.Type() != typ {
		t.Errorf("Nexter.Next() expecting Token.Type '%d', received '%d'", typ, tok.Type())
	}
	if tok.Value() != value {
		t.Errorf("Nexter.Next() expecting Token.String '%s', received '%s'", value, tok.Value())
	}
}

// TestTokensHasNext1
//
func TestTokensHasNext1(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitType(T_START)
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterHasNext(t, nexter, true)
	expectNexterNext(t, nexter, T_START, "")
	expectNexterHasNext(t, nexter, false)
}

// TestTokensHasNext2
//
func TestTokensHasNext2(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitType(T_START)
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterHasNext(t, nexter, true)
	expectNexterHasNext(t, nexter, true) // Call again, should hit cached 'next' value
	expectNexterNext(t, nexter, T_START, "")
	expectNexterHasNext(t, nexter, false)
}

// TestTokensEOF
//
func TestTokensEOF(t *testing.T) {
	nexter := LexString(".", nil)
	expectNexterHasNext(t, nexter, false)
}

// TestTokensNextAfterEOF
//
func TestTokensNextAfterEOF(t *testing.T) {
	nexter := LexString(".", nil)
	expectNexterHasNext(t, nexter, false)
	assertPanic(t, func() {
		nexter.Next()
	}, "Nexter.Next: No token available")
}
