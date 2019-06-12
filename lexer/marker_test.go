package lexer

import "testing"

// expectCanReset
//
func expectCanReset(t *testing.T, l *Lexer, m *Marker, match bool) {
	if r := l.CanReset(m); r != match {
		t.Errorf("Expecting Lexer.CanReset() to return %t, but received %t", match, r)
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
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
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
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
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
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
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
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
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
		LexString("123ABC", fn).Next()
	}, "Invalid marker")
}
