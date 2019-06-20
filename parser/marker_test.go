package parser

import (
	"testing"

	"github.com/tekwizely/go-parsing/lexer"
)

// expectCanReset
//
func expectCanReset(t *testing.T, p *Parser, m *Marker, match bool) {
	if p.CanReset(m) != match {
		t.Errorf("Parser.CanReset() expecting '%t'", match)
	}
}

// TestMarkerUnused
//
func TestMarkerUnused(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		m := p.Marker()
		expectCanReset(t, p, m, true)
		expectNext(t, p, lexer.T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(lexer.T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterHasNext(t, nexter, false)
}

// TestMarkerCanReset
//
func TestMarkerCanReset(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		m := p.Marker()
		expectCanReset(t, p, m, true)
		expectNext(t, p, lexer.T_START, "")
		p.Emit("T_START")
		expectCanReset(t, p, m, false)
		return nil
	}
	tokens := mockLexer(lexer.T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterHasNext(t, nexter, false)
}

// TestMarkerImmediateReset
//
func TestMarkerImmediateReset(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		m := p.Marker()
		expectCanReset(t, p, m, true)
		// Reset it immediately
		//
		p.Reset(m)
		expectCanReset(t, p, m, true)
		expectNext(t, p, lexer.T_START, "")
		p.Emit("T_START")
		expectCanReset(t, p, m, false)
		return nil
	}
	tokens := mockLexer(lexer.T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterHasNext(t, nexter, false)
}

// TestMarkerReset
//
func TestMarkerReset(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		m := p.Marker()
		expectCanReset(t, p, m, true)
		expectNext(t, p, lexer.T_START, "")
		expectCanReset(t, p, m, true)
		p.Reset(m)
		expectCanReset(t, p, m, true)
		expectNext(t, p, lexer.T_START, "")
		p.Emit("T_START")
		expectCanReset(t, p, m, false)
		return nil
	}
	tokens := mockLexer(lexer.T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterHasNext(t, nexter, false)
}

// TestMarkerResetInvalid
//
func TestMarkerResetInvalid(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		m := p.Marker()
		expectCanReset(t, p, m, true)
		expectNext(t, p, lexer.T_START, "")
		expectCanReset(t, p, m, true)
		p.Reset(m)
		expectCanReset(t, p, m, true)
		expectNext(t, p, lexer.T_START, "")
		p.Emit("T_START")
		expectCanReset(t, p, m, false)
		// CanReset said no, but let's try anyway
		//
		p.Reset(m)
		return nil
	}
	tokens := mockLexer(lexer.T_START)
	assertPanic(t, func() {
		Parse(tokens, fn).Next()
	}, "Invalid marker")
}
