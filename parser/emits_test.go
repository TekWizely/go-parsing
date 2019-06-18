package parser

import (
	"testing"

	"github.com/tekwizely/go-parsing/lexer"
)

// expectEmitsHasNext
//
func expectEmitsHasNext(t *testing.T, emits *Emits, match bool) {
	if emits.HasNext() != match {
		t.Errorf("Emits.HasNext() expecting '%t'", match)
	}
}

// expectEmitsNext
//
func expectEmitsNext(t *testing.T, emits *Emits, match string) {
	str := emits.Next().(string)
	if str != match {
		t.Errorf("Emits.Next() expecting '%s', received '%s'", match, str)
	}
}

// TestEmitsHasNext1
//
func TestEmitsHasNext1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, lexer.T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(lexer.T_START)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, true)
	expectEmitsNext(t, emits, "T_START")
	expectEmitsHasNext(t, emits, false)
}

// TestEmitsHasNext2
//
func TestEmitsHasNext2(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, lexer.T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(lexer.T_START)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, true)
	expectEmitsHasNext(t, emits, true) // Call again, should hit cached 'next' value
	expectEmitsNext(t, emits, "T_START")
	expectEmitsHasNext(t, emits, false)
}

// TestEmitEOF
//
func TestEmitsEOF(t *testing.T) {
	tokens := mockLexer()
	emits := Parse(tokens, nil)
	expectEmitsHasNext(t, emits, false)
}

// TestEmitsNextAfterEOF
//
func TestEmitsNextAfterEOF(t *testing.T) {
	tokens := mockLexer()
	emits := Parse(tokens, nil)
	expectEmitsHasNext(t, emits, false)
	assertPanic(t, func() {
		emits.Next()
	}, "Emits.Next: No AST available")
}
