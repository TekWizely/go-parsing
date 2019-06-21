package parser

import (
	"testing"
)

// expectNexterHasNext
//
func expectNexterHasNext(t *testing.T, nexter ASTNexter, match bool) {
	if nexter.HasNext() != match {
		t.Errorf("ASTNexter.HasNext() expecting '%t'", match)
	}
}

// expectNexterNext
//
func expectNexterNext(t *testing.T, nexter ASTNexter, match string) {
	str := nexter.Next().(string)
	if str != match {
		t.Errorf("ASTNexter.Next() expecting '%s', received '%s'", match, str)
	}
}

// TestNexterHasNext1
//
func TestNexterHasNext1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterHasNext(t, nexter, true)
	expectNexterNext(t, nexter, "T_START")
	expectNexterHasNext(t, nexter, false)
}

// TestNexterHasNext2
//
func TestNexterHasNext2(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterHasNext(t, nexter, true)
	expectNexterHasNext(t, nexter, true) // Call again, should hit cached 'next' value
	expectNexterNext(t, nexter, "T_START")
	expectNexterHasNext(t, nexter, false)
}

// TestEmitEOF
//
func TestNexterEOF(t *testing.T) {
	tokens := mockLexer()
	nexter := Parse(tokens, nil)
	expectNexterHasNext(t, nexter, false)
}

// TestNexterNextAfterEOF
//
func TestNexterNextAfterEOF(t *testing.T) {
	tokens := mockLexer()
	nexter := Parse(tokens, nil)
	expectNexterHasNext(t, nexter, false)
	assertPanic(t, func() {
		nexter.Next()
	}, "ASTNexter.Next: No AST available")
}
