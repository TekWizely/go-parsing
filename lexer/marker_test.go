package lexer

import "testing"

// expectCanReset
//
func expectCanReset(t *testing.T, l *Lexer, m *Marker, match bool) {
	if l.CanReset(m) != match {
		t.Errorf("Lexer.CanReset() expecting '%t'", match)
	}
}

// TestMarkerUnused
//
func TestMarkerUnused(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectCanReset(t, l, m, true)
		// Ignore marker
		//
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerCanReset
//
func TestMarkerCanReset(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectNextString(t, l, "123ABC")
		expectCanReset(t, l, m, true)
		l.EmitToken(T_STRING)
		expectCanReset(t, l, m, false)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerImmediateReset
//
func TestMarkerImmediateReset(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectCanReset(t, l, m, true)
		// Reset it immediately
		//
		l.Reset(m)
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		expectCanReset(t, l, m, false)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerReset
//
func TestMarkerReset(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectCanReset(t, l, m, true)
		expectNextString(t, l, "123ABC")
		expectCanReset(t, l, m, true)
		l.Reset(m)
		expectCanReset(t, l, m, true)
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		expectCanReset(t, l, m, false)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, T_STRING, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMarkerResetInvalid
//
func TestMarkerResetInvalid(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		m := l.Marker()
		expectCanReset(t, l, m, true)
		expectNextString(t, l, "123ABC")
		expectCanReset(t, l, m, true)
		l.Reset(m)
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		expectCanReset(t, l, m, false)
		// CanReset said no, but let's try anyway
		//
		l.Reset(m)
		return nil
	}
	assertPanic(t, func() {
		_, _ = LexString("123ABC", fn).Next()
	}, "Invalid marker")
}
