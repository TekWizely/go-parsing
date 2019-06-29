package parser

import (
	"io"
	"testing"
)

// expectNexterEOF confirms Next() == (nil, io.EOF)
//
func expectNexterEOF(t *testing.T, nexter ASTNexter) {
	ast, err := nexter.Next()
	if err == nil {
		if ast == nil {
			t.Errorf("Nexter.Next() expecting (nil, EOF), received (nil, nil)")
		} else {
			t.Errorf("Nexter.Next() expecting (nil, EOF), received ('%v', nil)", ast)
		}
	} else if ast != nil {
		t.Errorf("Nexter.Next() expecting (nil, EOF), received ('%v', '%s')'", ast, err.Error())
	} else if err != io.EOF {
		t.Errorf("Nexter.Next() expecting (nil, EOF), received (nil, '%s')", err.Error())
	}
}

// expectNexterNext confirms Next() == ("$match", nil)
//
func expectNexterNext(t *testing.T, nexter ASTNexter, match string) {
	ast, err := nexter.Next() // Assume ast, when non-nil, is of type string
	if ast == nil {
		if err == nil {
			t.Errorf("Nexter.Next() expecting ('%s', nil), received (nil, nil)'", match)
		} else {
			t.Errorf("Nexter.Next() expecting ('%s', nil), received (nil, '%s')'", match, err.Error())
		}
	} else if err != nil {
		t.Errorf("Nexter.Next() expecting ('%s', nil), received ('%v', '%s')'", match, ast, err.Error())
	} else if ast.(string) != match {
		t.Errorf("Nexter.Next() expecting ('%s', nil), received ('%v', nil)'", match, ast)
	}
}

// expectNexterError confirms Next() == (nil, "$errMsg")
//
func expectNexterError(t *testing.T, nexter ASTNexter, errMsg string) {
	ast, err := nexter.Next()
	if err == nil {
		if ast == nil {
			t.Errorf("Nexter.Next() expecting (nil, '%s'), received (nil, nil)", errMsg)
		} else {
			t.Errorf("Nexter.Next() expecting (nil, '%s'), received ('%v', nil)", errMsg, ast)
		}
	} else if ast != nil {
		t.Errorf("Nexter.Next() expecting (nil, '%s'), received ('%v', '%s')", errMsg, ast, err.Error())
	} else if err.Error() != errMsg {
		t.Errorf("Nexter.Next() expecting (nil, '%s'), received (nil, '%s')", errMsg, err.Error())
	}
}

// TestNexterHasNext1
//
func TestNexterHasNext1(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterEOF(t, nexter)
}

// TestNexterHasNext2
//
func TestNexterHasNext2(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterEOF(t, nexter)
}

// TestEmitEOF
//
func TestNexterEOF(t *testing.T) {
	tokens := mockLexer()
	nexter := Parse(tokens, nil)
	expectNexterEOF(t, nexter)
}

// TestNexterNextAfterEOF
//
func TestNexterNextAfterEOF(t *testing.T) {
	tokens := mockLexer()
	nexter := Parse(tokens, nil)
	expectNexterEOF(t, nexter)
	// Call again, should continue to return EOF
	//
	expectNexterEOF(t, nexter)
}
