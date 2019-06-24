package lexer

import (
	"testing"
)

// expectMarkerValid
//
func expectMarkerValid(t *testing.T, m *Marker, match bool) {
	if m.Valid() != match {
		t.Errorf("Maker.Valid() expecting '%t'", match)
	}
}

// TestMarkerUnused
//
func TestMarkerUnused(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectMarkerValid(t, m, true)
		// Ignore marker
		//
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerValid
//
func TestMarkerValid(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectNextString(t, l, "123ABC")
		expectMarkerValid(t, m, true)
		l.EmitToken(T_STRING)
		expectMarkerValid(t, m, false)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerImmediateApply
//
func TestMarkerImmediateApply(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectMarkerValid(t, m, true)
		// Apply it immediately
		//
		m.Apply()
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		expectMarkerValid(t, m, false)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerApply
//
func TestMarkerApply(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectMarkerValid(t, m, true)
		expectNextString(t, l, "123ABC")
		expectMarkerValid(t, m, true)
		m.Apply()
		expectMarkerValid(t, m, true)
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		expectMarkerValid(t, m, false)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerApplyInvalid
//
func TestMarkerApplyInvalid(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectMarkerValid(t, m, true)
		expectNextString(t, l, "123ABC")
		expectMarkerValid(t, m, true)
		m.Apply()
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		expectMarkerValid(t, m, false)
		// Valid said no, but let's try anyway
		//
		m.Apply()
		return nil
	}
	assertPanic(t, func() {
		_, _ = LexString("123ABC", fn).Next()
	}, "Invalid marker")
}

func TestMarkerApplyNextFn(t *testing.T) {

	var marker *Marker
	var used = false

	fn1 := func(l *Lexer) LexerFn {
		if used {
			t.Error("Marker.Apply() expected to return function that marker was created in")
			return nil
		}
		used = true
		return marker.Apply()
	}

	fn2 := func(l *Lexer) LexerFn {
		if used {
			return nil
		}
		marker = l.Marker()
		return fn1
	}
	nexter := LexString(".", fn2)
	expectNexterEOF(t, nexter)
}
