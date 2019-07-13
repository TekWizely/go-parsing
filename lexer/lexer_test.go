package lexer

import (
	"io"
	"log"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// Define tokens used in various tests
//
const (
	TInt token.Type = TStart + iota
	TChar
	TString
)

// runeReaderErr
//
type runeReaderErr struct {
	err error
}

func (rr *runeReaderErr) ReadRune() (r rune, size int, err error) {
	return 0, 0, rr.err
}

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
	if l.CanPeek(peek) != match {
		t.Errorf("Lexer.CanPeek(%d) expecting '%t'", peek, match)
	}
}

// expectPeek
//
func expectPeek(t *testing.T, l *Lexer, peek int, match rune) {
	if r := l.Peek(peek); r != match {
		t.Errorf("Expecting Lexer.Peek(%d) to return '%c', but received '%c'", peek, match, r)
	}
}

// expectNext
//
func expectNext(t *testing.T, l *Lexer, match rune) {
	if r := l.Next(); r != match {
		t.Errorf("Lexer.Next() expecting rune '%c', received '%c'", match, r)
	}
}

// expectPeekToken
//
func expectPeekToken(t *testing.T, l *Lexer, match string) {
	s := l.PeekToken()
	if s != match {
		t.Errorf("Lexer.PeekToken() expecting '%s', received '%s'", match, s)
	}
}

// expectEOF
//
func expectEOF(t *testing.T, l *Lexer) {
	eof := l.eof && l.cache.Len() == l.matchLen
	if !eof {
		t.Error("Lexer expecting to be at EOF")
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
		expectNext(t, l, r[i])
	}
	expectPeekToken(t, l, match)
}

// expectMatchEmitString
//
func expectMatchEmitString(t *testing.T, l *Lexer, match string, emitType token.Type) {
	expectNextString(t, l, match)
	if t.Failed() == false {
		l.EmitToken(emitType)
	}
}

// TestNilFn
//
func TestNilFn(t *testing.T) {
	nexter := LexString("", nil)
	expectNexterEOF(t, nexter)
}

// TestLexerFnSkippedWhenNoCanPeek
//
func TestLexerFnSkippedWhenNoCanPeek(t *testing.T) {
	fn := func(l *Lexer) Fn {
		t.Error("Lexer should not call Lexer.Fn when CanPeek(1) == false")
		return nil
	}
	nexter := LexString("", fn)
	expectNexterEOF(t, nexter)
}

// TestEmittoken.Ttype
//
func TestEmitEmptyType(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitType(TStart)
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterNext(t, nexter, TStart, "")
	expectNexterEOF(t, nexter)
}

// TestEmitEmptyToken
//
func TestEmitEmptyToken(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitToken(TStart)
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterNext(t, nexter, TStart, "")
	expectNexterEOF(t, nexter)
}

// TestCanPeek
//
func TestCanPeek(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectCanPeek(t, l, 1, true)

		expectPeek(t, l, 1, '1')

		expectCanPeek(t, l, 2, true)

		expectPeek(t, l, 2, '2')

		expectCanPeek(t, l, 3, true)

		expectPeek(t, l, 3, '3')

		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestCanPeekPastEOF
//
func TestCanPeekPastEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectCanPeek(t, l, 4, false)

		expectCanPeek(t, l, 3, true)

		expectPeek(t, l, 3, '3')

		expectCanPeek(t, l, 2, true)

		expectPeek(t, l, 2, '2')

		expectCanPeek(t, l, 1, true)

		expectPeek(t, l, 1, '1')

		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestCanPeekRangeError
//
func TestCanPeekRangeError(t *testing.T) {
	fn := func(l *Lexer) Fn {
		assertPanic(t, func() {
			l.CanPeek(-1)
		}, "Lexer.CanPeek: range error")
		assertPanic(t, func() {
			l.CanPeek(0)
		}, "Lexer.CanPeek: range error")
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterEOF(t, nexter)
}

// TestPeek1
//
func TestPeek1(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectPeek(t, l, 1, '1')
		return nil
	}
	nexter := LexString("1", fn)
	expectNexterEOF(t, nexter)

}

// TestPeek11
//
func TestPeek11(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectPeek(t, l, 1, 'A')
		expectPeek(t, l, 1, 'A')
		return nil
	}
	nexter := LexString("AB", fn)
	expectNexterEOF(t, nexter)
}

// TestPeek12
//
func TestPeek12(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectPeek(t, l, 1, 'A')
		expectPeek(t, l, 2, 'B')
		return nil
	}
	nexter := LexString("AB", fn)
	expectNexterEOF(t, nexter)
}

// TestPeekEmpty
//
func TestPeekEmpty(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNext(t, l, '.')
		assertPanic(t, func() {
			expectPeek(t, l, 1, utf8.RuneError)
		}, "Lexer.Peek: No rune available")
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterEOF(t, nexter)
}

// TestPeekRangeError
//
func TestPeekRangeError(t *testing.T) {
	fn := func(l *Lexer) Fn {
		assertPanic(t, func() {
			l.Peek(-1)
		}, "Lexer.Peek: range error")
		assertPanic(t, func() {
			l.Peek(0)
		}, "Lexer.Peek: range error")
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterEOF(t, nexter)
}

// TestNext1
//
func TestNext1(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		return nil
	}
	nexter := LexString("AB", fn)
	expectNexterEOF(t, nexter)
}

// TestNext2
//
func TestNext2(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		expectPeek(t, l, 1, 'B')
		expectNext(t, l, 'B')
		return nil
	}
	nexter := LexString("AB", fn)
	expectNexterEOF(t, nexter)
}

// TestNextEmpty
//
func TestNextEmpty(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNext(t, l, '.')
		assertPanic(t, func() {
			expectNext(t, l, utf8.RuneError)
		}, "Lexer.Next: No rune available")
		return nil
	}
	nexter := LexString(".", fn)
	expectNexterEOF(t, nexter)
}

// TestNextEmit1
//
func TestNextEmit1(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		l.EmitToken(TChar)
		return nil
	}
	nexter := LexString("AB", fn)
	expectNexterNext(t, nexter, TChar, "A")
	expectNexterEOF(t, nexter)
}

// TestNextEmit2
//
func TestNextEmit2(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectPeek(t, l, 1, 'A')
		expectNext(t, l, 'A')
		l.EmitToken(TChar)
		expectPeek(t, l, 1, 'B')
		expectNext(t, l, 'B')
		l.EmitToken(TChar)
		return nil
	}
	nexter := LexString("AB", fn)
	expectNexterNext(t, nexter, TChar, "A")
	expectNexterNext(t, nexter, TChar, "B")
	expectNexterEOF(t, nexter)
}

// TestMatchInt
//
func TestMatchInt(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectMatchEmitString(t, l, "123", TInt)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterNext(t, nexter, TInt, "123")
	expectNexterEOF(t, nexter)
}

// TestMatchIntString
//
func TestMatchIntString(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectMatchEmitString(t, l, "123", TInt)
		expectMatchEmitString(t, l, "ABC", TString)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, TInt, "123")
	expectNexterNext(t, nexter, TString, "ABC")
	expectNexterEOF(t, nexter)
}

// TestMatchString
//
func TestMatchString(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectMatchEmitString(t, l, "123ABC", TString)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, TString, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMatchRunes
//
func TestMatchRunes(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectMatchEmitString(t, l, "123ABC", TString)
		return nil
	}
	nexter := LexRunes([]rune("123ABC"), fn)
	expectNexterNext(t, nexter, TString, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMatchBytes
//
func TestMatchBytes(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectMatchEmitString(t, l, "123ABC", TString)
		return nil
	}
	nexter := LexBytes([]byte("123ABC"), fn)
	expectNexterNext(t, nexter, TString, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestMatchReader
//
func TestMatchReader(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectMatchEmitString(t, l, "123ABC", TString)
		return nil
	}
	nexter := LexReader(strings.NewReader("123ABC"), fn)
	expectNexterNext(t, nexter, TString, "123ABC")
	expectNexterEOF(t, nexter)
}

// TestClear1
//
func TestClear1(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNextString(t, l, "123ABC")
		l.Clear()
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterEOF(t, nexter)
}

// TestClear2
//
func TestClear2(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectMatchEmitString(t, l, "123", TInt)
		expectNextString(t, l, "ABC")
		l.Clear()
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, TInt, "123")
	expectNexterEOF(t, nexter)
}

// TestClear3
//
func TestClear3(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNextString(t, l, "123")
		l.Clear()
		expectMatchEmitString(t, l, "ABC", TString)
		return nil
	}
	nexter := LexString("123ABC", fn)
	expectNexterNext(t, nexter, TString, "ABC")
	expectNexterEOF(t, nexter)
}

// TestEmitEOF1
//
func TestEmitEOF1(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitEOF()
		expectEOF(t, l)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestEmitEOF2
//
func TestEmitEOF2(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestEmitEOF3
//
func TestEmitEOF3(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitType(TEof)
		expectEOF(t, l)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestEmitEOF4
//
func TestEmitEOF4(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitToken(TEof)
		expectEOF(t, l)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestEmitEOF5
//
func TestEmitEOF5(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNextString(t, l, "123")
		l.EmitToken(TEof)
		expectEOF(t, l)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestEmitError
//
func TestEmitError(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitError("ERROR")
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterError(t, nexter, "ERROR")
	expectNexterEOF(t, nexter)
}

// TestEmitErrorf
//
func TestEmitErrorf(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitErrorf("ERROR: %s %d", "Error", 1)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterError(t, nexter, "ERROR: Error 1")
	expectNexterEOF(t, nexter)
}

// TestEmitAfterEOF
//
func TestEmitAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		l.EmitToken(TInt)
		return nil
	}
	assertPanic(t, func() {
		_, _ = LexString("123", fn).Next()
	}, "Lexer.EmitToken: No further emits allowed after EOF is emitted")
}

// TestEmitTypeAfterEOF
//
func TestEmitTypeAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitEOF()
		expectEOF(t, l)
		l.EmitType(TStart)
		return nil
	}
	assertPanic(t, func() {
		_, _ = LexString("123", fn).Next()
	}, "Lexer.EmitType: No further emits allowed after EOF is emitted")
}

// TestCanPeekAfterEOF
//
func TestCanPeekAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitEOF()
		expectEOF(t, l)
		expectCanPeek(t, l, 1, false)
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestPeekAfterEOF
//
func TestPeekAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitEOF()
		expectEOF(t, l)
		assertPanic(t, func() {
			l.Peek(1)
		}, "Lexer.Peek: No runes can be peeked after EOF is emitted")
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestNextAfterEOF
//
func TestNextAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitEOF()
		expectEOF(t, l)
		assertPanic(t, func() {
			l.Next()
		}, "Lexer.Next: No runes can be matched after EOF is emitted")
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestPeekTokenAfterEOF
//
func TestPeekTokenAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		assertPanic(t, func() {
			l.PeekToken()
		}, "Lexer.PeekToken: No token peeks allowed after EOF is emitted")
		return nil
	}
	nexter := LexString("123", fn)
	expectNexterEOF(t, nexter)
}

// TestEmitErrorAfterEOF
//
func TestEmitErrorAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		l.EmitEOF()
		expectEOF(t, l)
		l.EmitError("ERROR")
		return nil
	}
	assertPanic(t, func() {
		_, _ = LexString("123", fn).Next()
	}, "Lexer.EmitError: No further emits allowed after EOF is emitted")
}

// TestClearAfterEOF
//
func TestClearAfterEOF(t *testing.T) {
	fn := func(l *Lexer) Fn {
		expectNextString(t, l, "123")
		l.EmitEOF()
		expectEOF(t, l)
		l.Clear()
		return nil
	}
	assertPanic(t, func() {
		_, _ = LexString("123", fn).Next()
	}, "Lexer.Clear: No clears allowed after EOF is emitted")
}

// TestRuneReaderNonEOFError should log an error but otherwise behave as EOF
//
func TestRuneReaderNonEOFError(t *testing.T) {
	sb := &strings.Builder{}
	log.SetFlags(0)
	log.SetOutput(sb)
	fn := func(l *Lexer) Fn {
		return nil
	}
	nexter := LexRuneReader(&runeReaderErr{err: io.ErrUnexpectedEOF}, fn)
	expectNexterEOF(t, nexter)
	if log := sb.String(); log != "non-EOF error returned from rune reader, treating as EOF: unexpected EOF\n" {
		t.Errorf("Lexer.growPeek received wrong log message: '%s'", log)
	}
}
