package parser

import (
	"testing"
)

// expectMarkerValid
//
func expectMarkerValid(t *testing.T, m *Marker, match bool) {
	if m.Valid() != match {
		t.Errorf("Marker.Valid() expecting '%t'", match)
	}
}

// TestMarkerUnused
//
func TestMarkerUnused(t *testing.T) {
	fn := func(p *Parser) Fn {
		m := p.Marker()
		expectMarkerValid(t, m, true)
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterEOF(t, nexter)
}

// TestMarkerValid
//
func TestMarkerValid(t *testing.T) {
	fn := func(p *Parser) Fn {
		m := p.Marker()
		expectMarkerValid(t, m, true)
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		expectMarkerValid(t, m, false)
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterEOF(t, nexter)
}

// TestMarkerImmediateApply
//
func TestMarkerImmediateApply(t *testing.T) {
	fn := func(p *Parser) Fn {
		m := p.Marker()
		expectMarkerValid(t, m, true)
		// Apply it immediately
		//
		m.Apply()
		expectMarkerValid(t, m, true)
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		expectMarkerValid(t, m, false)
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterEOF(t, nexter)
}

// TestMarkerApply
//
func TestMarkerApply(t *testing.T) {
	fn := func(p *Parser) Fn {
		m := p.Marker()
		expectMarkerValid(t, m, true)
		expectNext(t, p, T_START, "")
		expectMarkerValid(t, m, true)
		m.Apply()
		expectMarkerValid(t, m, true)
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		expectMarkerValid(t, m, false)
		return nil
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterEOF(t, nexter)
}

// TestMarkerApplyInvalid
//
func TestMarkerApplyInvalid(t *testing.T) {
	fn := func(p *Parser) Fn {
		m := p.Marker()
		expectMarkerValid(t, m, true)
		expectNext(t, p, T_START, "")
		expectMarkerValid(t, m, true)
		m.Apply()
		expectMarkerValid(t, m, true)
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		expectMarkerValid(t, m, false)
		// Valid said no, but let's try anyway
		//
		m.Apply()
		return nil
	}
	tokens := mockLexer(T_START)
	assertPanic(t, func() {
		_, _ = Parse(tokens, fn).Next()
	}, "Invalid marker")
}

func TestMarkerApplyNextFn(t *testing.T) {

	var marker *Marker
	var used = false

	fn1 := func(p *Parser) Fn {
		if used {
			t.Error("Marker.Apply() expected to return function that marker was created in")
			return nil
		}
		used = true
		return marker.Apply()
	}

	fn2 := func(p *Parser) Fn {
		if used {
			return nil
		}
		marker = p.Marker()
		return fn1
	}
	tokens := mockLexer(T_START)
	nexter := Parse(tokens, fn2)
	expectNexterEOF(t, nexter)
}
