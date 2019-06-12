package lexer

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// Define tokens used in various tests
//
const (
	T_INT TokenType = T_START + iota
	T_CHAR
	T_STRING
)

// assertPanic
//
func assertPanic(t *testing.T, f func(), msg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("assertPanic: did not generate panic()")
		} else if r != msg {
			t.Errorf("assertPanic: recover() recieved message '%s' instead of '%s'", r, msg)
		}
	}()
	f()
}

// expectCanPeek
//
func expectCanPeek(t *testing.T, l *Lexer, peek int, match bool) {
	if r := l.CanPeek(peek); r != match {
		t.Errorf("Expecting Lexer.CanPeek(%d) to return '%t', but received '%t'", peek, match, r)
	}
}

// expectPeek
//
func expectPeek(t *testing.T, l *Lexer, peek int, match rune) {
	if r := l.Peek(peek); r != match {
		t.Errorf("Expecting Lexer.Peek(%d) to return '%c', but received '%c'", peek, match, r)
	}
}

// expectHasNext
//
func expectHasNext(t *testing.T, l *Lexer, match bool) {
	if r := l.HasNext(); r != match {
		t.Errorf("Expecting Lexer.HasNext() to return '%t', but received '%t'", match, r)
	}
}

// expectNext
//
func expectNext(t *testing.T, l *Lexer, match rune) {
	if r := l.Next(); r != match {
		t.Errorf("Expecting Lexer.Next() to return '%c', but received '%c'", match, r)
	}
}

// expectPeekToken
//
func expectPeekToken(t *testing.T, l *Lexer, match string) {
	s := l.PeekToken()
	if s != match {
		t.Errorf("Expecting Lexer.PeekToken() to match '%s', received '%s'", match, s)
	}
}

// expectPeekTokenRunes
//
func expectPeekTokenRunes(t *testing.T, l *Lexer, match string) {
	r := l.PeekTokenRunes()
	s := string(r)
	if s != match {
		t.Errorf("Expecting Lexer.PeekTokenRunes() to match '%s', received '%s'", match, s)
	}
}

// expectEOF
//
func expectEOF(t *testing.T, l *Lexer) {
	eof := l.eof && l.runes.Len() == l.tokenLen
	if !eof {
		t.Error("Expecting to be at EOF")
	}
}

// expectPeekString
//
func expectPeekString(t *testing.T, l *Lexer, match string) {
	r := []rune(match)
	for i := 0; i < len(r); i++ {
		expectCanPeek(t, l, i+1, true) // 1-based
		expectPeek(t, l, i+1, r[i])    // 1-based
	}
}

// expectNextString
//
func expectNextString(t *testing.T, l *Lexer, match string) {
	expectPeekString(t, l, match)
	r := []rune(match)
	for i := 0; i < len(r); i++ {
		expectCanPeek(t, l, 1, true)
		expectPeek(t, l, 1, r[i])
		expectHasNext(t, l, true)
		expectNext(t, l, r[i])
	}
	expectPeekToken(t, l, match)
	expectPeekTokenRunes(t, l, match)
}

// expectMatchEmitString
//
func expectMatchEmitString(t *testing.T, l *Lexer, match string, emitType TokenType) {
	expectNextString(t, l, match)
	if t.Failed() == false {
		l.EmitToken(emitType)
	}
}

// TestNilFn
//
func TestNilFn(t *testing.T) {
	tokens := LexString("", nil)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitTokenType
//
func TestLexerFnSkippedWhenNoHasNext(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		t.Error("Lexer should not call LexerFn when HasNext() == false")
		return nil
	}
	tokens := LexString("", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitTokenType
//
func TestEmitEmptyType(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitType(T_START)
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, true)
	expectTokensNext(t, tokens, T_START, "")
	expectTokensHasNext(t, tokens, false)
}

// TestEmitEmptyToken
//
func TestEmitEmptyToken(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitToken(T_START)
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, true)
	expectTokensNext(t, tokens, T_START, "")
	expectTokensHasNext(t, tokens, false)
}

// TestCanPeekPastEOF
//
func TestCanPeekPastEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectCanPeek(t, l, 4, false)

		expectCanPeek(t, l, 3, true)

		expectPeek(t, l, 3, '3')

		expectCanPeek(t, l, 2, true)

		expectPeek(t, l, 2, '2')

		expectCanPeek(t, l, 1, true)

		expectPeek(t, l, 1, '1')

		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestCanPeekRangeError
//
func TestCanPeekRangeError(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		assertPanic(t, func() {
			l.CanPeek(-1)
		}, "Lexer.CanPeek: range error")
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestPeek1
//
func TestPeek1(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectPeek(t, l, 1, '1')
		return nil
	}
	tokens := LexString("1", fn)
	expectTokensHasNext(t, tokens, false)

}

// TestPeek11
//
func TestPeek11(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectPeek(t, l, 1, 'A')
		expectPeek(t, l, 1, 'A')
		return nil
	}
	tokens := LexString("AB", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestPeek12
//
func TestPeek12(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectPeek(t, l, 1, 'A')
		expectPeek(t, l, 2, 'B')
		return nil
	}
	tokens := LexString("AB", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestPeekEmpty
//
func TestPeekEmpty(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNext(t, l, '.')
		assertPanic(t, func() {
			expectPeek(t, l, 1, utf8.RuneError)
		}, "Lexer.Peek: No rune available")
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestPeekRangeError
//
func TestPeekRangeError(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		assertPanic(t, func() {
			l.Peek(-1)
		}, "Lexer.Peek: range error")
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestHasNextTrue
//
func TestHasNextTrue(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectHasNext(t, l, true)
		return nil
	}
	tokens := LexString("1", fn)
	expectTokensHasNext(t, tokens, false)

}

// TestHasNextFalse
//
func TestHasNextFalse(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectHasNext(t, l, true)
		expectNext(t, l, '.')
		expectHasNext(t, l, false)
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, false)

}

// TestNext1
//
func TestNext1(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		return nil
	}
	tokens := LexString("AB", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestNext2
//
func TestNext2(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		expectPeek(t, l, 1, 'B')
		expectNext(t, l, 'B')
		return nil
	}
	tokens := LexString("AB", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestNextEmpty
//
func TestNextEmpty(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNext(t, l, '.')
		assertPanic(t, func() {
			expectNext(t, l, utf8.RuneError)
		}, "Lexer.Next: No rune available")
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestNextEmit1
//
func TestNextEmit1(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		l.EmitToken(T_CHAR)
		return nil
	}
	tokens := LexString("AB", fn)
	expectTokensNext(t, tokens, T_CHAR, "A")
	expectTokensHasNext(t, tokens, false)
}

// TestNextEmit2
//
func TestNextEmit2(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		l.EmitToken(T_CHAR)
		expectPeek(t, l, 1, 'B')
		expectNext(t, l, 'B')
		l.EmitToken(T_CHAR)
		return nil
	}
	tokens := LexString("AB", fn)
	expectTokensNext(t, tokens, T_CHAR, "A")
	expectTokensNext(t, tokens, T_CHAR, "B")
	expectTokensHasNext(t, tokens, false)
}

// TestMatchInt
//
func TestMatchInt(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectMatchEmitString(t, l, "123", T_INT)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensNext(t, tokens, T_INT, "123")
	expectTokensHasNext(t, tokens, false)
}

// TestMatchIntString
//
func TestMatchIntString(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectMatchEmitString(t, l, "123", T_INT)
		expectMatchEmitString(t, l, "ABC", T_STRING)
		return nil
	}
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_INT, "123")
	expectTokensNext(t, tokens, T_STRING, "ABC")
	expectTokensHasNext(t, tokens, false)
}

// TestMatchString
//
func TestMatchString(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		return nil
	}
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
}

// TestMatchRunes
//
func TestMatchRunes(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		return nil
	}
	tokens := LexRunes([]rune("123ABC"), fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
}

// TestMatchBytes
//
func TestMatchBytes(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		return nil
	}
	tokens := LexBytes([]byte("123ABC"), fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
}

// TestMatchReader
//
func TestMatchReader(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectMatchEmitString(t, l, "123ABC", T_STRING)
		return nil
	}
	tokens := LexReader(strings.NewReader("123ABC"), fn)
	expectTokensNext(t, tokens, T_STRING, "123ABC")
	expectTokensHasNext(t, tokens, false)
}

// TestDiscardToken1
//
func TestDiscardToken1(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNextString(t, l, "123ABC")
		l.DiscardToken()
		return nil
	}
	tokens := LexString("123ABC", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestDiscardToken2
//
func TestDiscardToken2(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectMatchEmitString(t, l, "123", T_INT)
		expectNextString(t, l, "ABC")
		l.DiscardToken()
		return nil
	}
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_INT, "123")
	expectTokensHasNext(t, tokens, false)
}

// TestDiscardToken3
//
func TestDiscardToken3(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNextString(t, l, "123")
		l.DiscardToken()
		expectMatchEmitString(t, l, "ABC", T_STRING)
		return nil
	}
	tokens := LexString("123ABC", fn)
	expectTokensNext(t, tokens, T_STRING, "ABC")
	expectTokensHasNext(t, tokens, false)
}

// TestEmitEOF1
//
func TestEmitEOF1(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitEOF()
		expectEOF(t, l)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitEOF2
//
func TestEmitEOF2(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitEOF3
//
func TestEmitEOF3(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitType(T_EOF)
		expectEOF(t, l)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitEOF4
//
func TestEmitEOF4(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitToken(T_EOF)
		expectEOF(t, l)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitEOF5
//
func TestEmitEOF5(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNextString(t, l, "123")
		l.EmitToken(T_EOF)
		expectEOF(t, l)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitError
//
func TestEmitError(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitError("ERROR")
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensNext(t, tokens, T_LEX_ERR, "ERROR")
	expectTokensHasNext(t, tokens, false)
}

// TestEmitErrorf
//
func TestEmitErrorf(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitErrorf("ERROR: %s %d", "Error", 1)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensNext(t, tokens, T_LEX_ERR, "ERROR: Error 1")
	expectTokensHasNext(t, tokens, false)
}

// TestEmitAfterEOF
//
func TestEmitAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		l.EmitToken(T_INT)
		return nil
	}
	assertPanic(t, func() {
		LexString("123", fn).Next()
	}, "Lexer.EmitToken: No further emits allowed after EOF is emitted")
}

// TestEmitTypeAfterEOF
//
func TestEmitTypeAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitEOF()
		expectEOF(t, l)
		l.EmitType(T_START)
		return nil
	}
	assertPanic(t, func() {
		LexString("123", fn).Next()
	}, "Lexer.EmitType: No further emits allowed after EOF is emitted")
}

// TestCanPeekAfterEOF
//
func TestCanPeekAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitEOF()
		expectEOF(t, l)
		expectCanPeek(t, l, 1, false)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestPeekAfterEOF
//
func TestPeekAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitEOF()
		expectEOF(t, l)
		assertPanic(t, func() {
			l.Peek(1)
		}, "Lexer.Peek: No runes can be peeked after EOF is emitted")
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestHasNextAfterEOF
//
func TestHasNextAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitEOF()
		expectEOF(t, l)
		expectHasNext(t, l, false)
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestNextAfterEOF
//
func TestNextAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitEOF()
		expectEOF(t, l)
		assertPanic(t, func() {
			l.Next()
		}, "Lexer.Next: No runes can be consumed after EOF is emitted")
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestPeekTokenAfterEOF
//
func TestPeekTokenAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		assertPanic(t, func() {
			l.PeekToken()
		}, "Lexer.PeekToken: No token peeks allowed after EOF is emitted")
		return nil
	}
	tokens := LexString("123", fn)
	expectTokensHasNext(t, tokens, false)
}

// TestEmitErrorAfterEOF
//
func TestEmitErrorAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitEOF()
		expectEOF(t, l)
		l.EmitError("ERROR")
		return nil
	}
	assertPanic(t, func() {
		LexString("123", fn).Next()
	}, "Lexer.EmitError: No further emits allowed after EOF is emitted")
}

// TestDiscardAfterEOF
//
func TestDiscardAfterEOF(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		l.DiscardToken()
		return nil
	}
	assertPanic(t, func() {
		LexString("123", fn).Next()
	}, "Lexer.Discard: No discards allowed after EOF is emitted")
}
