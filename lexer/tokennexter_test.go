package lexer

import (
	"io"
	"testing"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// expectNexterEOF confirms Next() == (nil, io.EOF)
//
func expectNexterEOF(t *testing.T, nexter token.Nexter) {
	tok, err := nexter.Next()
	if err == nil {
		if tok == nil {
			t.Errorf("Nexter.Next() expecting (nil, EOF), received (nil, nil)")
		} else {
			t.Errorf("Nexter.Next() expecting (nil, EOF), received ({%d, '%s'}, nil)", tok.Type(), tok.Value())
		}
	} else if tok != nil {
		t.Errorf("Nexter.Next() expecting (nil, EOF), received ({%d, '%s'}, '%s')'", tok.Type(), tok.Value(), err.Error())
	} else if err != io.EOF {
		t.Errorf("Nexter.Next() expecting (nil, EOF), received (nil, '%s')", err.Error())
	}
}

// expectNexterNext confirms Next() == (Token{type, value}, nil)
//
func expectNexterNext(t *testing.T, nexter token.Nexter, typ token.Type, value string) {
	tok, err := nexter.Next()
	if tok == nil {
		if err == nil {
			t.Errorf("Nexter.Next() expecting ({%d, '%s'}, nil), received (nil, nil)'", typ, value)
		} else {
			t.Errorf("Nexter.Next() expecting ({%d, '%s'}, nil), received (nil, '%s')'", typ, value, err.Error())
		}
	} else if err != nil {
		t.Errorf("Nexter.Next() expecting ({%d, '%s'}, nil), received ({%d, '%s'}, '%s')'", typ, value, tok.Type(), tok.Value(), err.Error())
	} else if tok.Type() != typ || tok.Value() != value {
		t.Errorf("Nexter.Next() expecting ({%d, '%s'}, nil), received ({%d, '%s'}, nil)'", typ, value, tok.Type(), tok.Value())
	}
}

// expectNexterError confirms Next() == (nil, "$errMsg")
//
func expectNexterError(t *testing.T, nexter token.Nexter, errMsg string) {
	tok, err := nexter.Next()
	if err == nil {
		if tok == nil {
			t.Errorf("Nexter.Next() expecting (nil, '%s'), received (nil, nil)", errMsg)
		} else {
			t.Errorf("Nexter.Next() expecting (nil, '%s'), received ({%d, '%s'}, nil)", errMsg, tok.Type(), tok.Value())
		}
	} else if tok != nil {
		t.Errorf("Nexter.Next() expecting (nil, '%s'), received ({%d, '%s'}, '%s')", errMsg, tok.Type(), tok.Value(), err.Error())
	} else if err.Error() != errMsg {
		t.Errorf("Nexter.Next() expecting (nil, '%s'), received (nil, '%s')", errMsg, err.Error())
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
	expectNexterNext(t, nexter, T_START, "")
	expectNexterEOF(t, nexter)
}

// TestTokensHasNext2
//
func TestTokensHasNext2(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitType(T_START)
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterNext(t, nexter, T_START, "")
	expectNexterEOF(t, nexter)
}

// TestTokensEOF
//
func TestTokensEOF(t *testing.T) {
	nexter := LexString(".", nil)
	expectNexterEOF(t, nexter)
}

// TestTokensNextAfterEOF
//
func TestTokensNextAfterEOF(t *testing.T) {
	nexter := LexString(".", nil)
	expectNexterEOF(t, nexter)
	// Call again, should continue to return EOF
	//
	expectNexterEOF(t, nexter)
}
