package parser

import (
	"testing"

	"github.com/tekwizely/go-parsing/lexer"
)

// Define tokens used in various tests
//
const (
	T_START lexer.TokenType = lexer.T_START + iota
	T_ONE
	T_TWO
	T_THREE
)

// mockLexer
//
func mockLexer(tokens ...lexer.TokenType) lexer.TokenNexter {
	i := 0
	var fn lexer.LexerFn
	fn = func(l *lexer.Lexer) lexer.LexerFn {
		if i >= len(tokens) {
			return nil
		}
		l.EmitType(tokens[i])
		i++
		return fn
	}
	return lexer.LexString(".", fn) // Input can't be empty or LexerFn will not be called
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
func expectCanPeek(t *testing.T, p *Parser, peek int, match bool) {
	if p.CanPeek(peek) != match {
		t.Errorf("Parser.CanPeek(%d) expecting '%t'", peek, match)
	}
}

// expectPeekType
//
func expectPeekType(t *testing.T, p *Parser, peek int, match lexer.TokenType) {
	if typ := p.PeekType(peek); typ != match {
		t.Errorf("Parser.PeekType(%d) expecting Token.Type '%s', received '%s'", peek, match, typ)
	}
}

// expectPeek
//
func expectPeek(t *testing.T, p *Parser, peek int, typ lexer.TokenType, str string) {
	tok := p.Peek(peek)
	if tok.Type != typ {
		t.Errorf("Parser.Peek(%d) expecting Token.Type '%s', received '%s'", peek, typ, tok.Type)
	}
	if tok.String != str {
		t.Errorf("Parser.Peek(%d) expecting Token.String '%s', received '%s'", peek, str, tok.String)
	}
}

// expectHasNext
//
func expectHasNext(t *testing.T, p *Parser, match bool) {
	if p.HasNext() != match {
		t.Errorf("Parser.HasNext() expecting '%t'", match)
	}
}

// expectNext
//
func expectNext(t *testing.T, p *Parser, typ lexer.TokenType, str string) {
	tok := p.Next()
	if tok.Type != typ {
		t.Errorf("Parser.Next() expecting Token.Type '%s', received '%s'", typ, tok.Type)
	}
	if tok.String != str {
		t.Errorf("Parser.Next() expecting Token.String '%s', received '%s'", str, tok.String)
	}
}

// expectEOF
//
func expectEOF(t *testing.T, p *Parser) {
	eof := p.eof && p.cache.Len() == p.matchLen
	if !eof {
		t.Error("Parser expecting to be at EOF")
	}
}

// TestNilFn
//
func TestNilFn(t *testing.T) {
	tokens := mockLexer()
	emits := Parse(tokens, nil)
	expectEmitsHasNext(t, emits, false)
}

func TestParserFnSkipedWhenNoHasNext(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		t.Error("Parser should not call ParserFn when HasNext() == false")
		return nil
	}
	tokens := mockLexer()
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)

}

// TestEmit
//
func TestEmit(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, T_START, "")
		p.Emit("T_START")
		return nil
	}
	tokens := mockLexer(T_START)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, true)
	expectEmitsNext(t, emits, "T_START")
	expectEmitsHasNext(t, emits, false)
}

// TestCanPeek
//
func TestCanPeek(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectCanPeek(t, p, 1, true)

		expectPeekType(t, p, 1, T_ONE)

		expectCanPeek(t, p, 2, true)

		expectPeekType(t, p, 2, T_TWO)

		expectCanPeek(t, p, 3, true)

		expectPeekType(t, p, 3, T_THREE)

		return nil
	}
	tokens := mockLexer(T_ONE, T_TWO, T_THREE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestCanPeekPastEOF
//
func TestCanPeekPastEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectCanPeek(t, p, 4, false)

		expectCanPeek(t, p, 3, true)

		expectPeekType(t, p, 3, T_THREE)

		expectCanPeek(t, p, 2, true)

		expectPeekType(t, p, 2, T_TWO)

		expectCanPeek(t, p, 1, true)

		expectPeekType(t, p, 1, T_ONE)

		return nil
	}
	tokens := mockLexer(T_ONE, T_TWO, T_THREE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestCanPeekRangeError
//
func TestCanPeekRangeError(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		assertPanic(t, func() {
			p.CanPeek(-1)
		}, "Parser.CanPeek: range error")
		assertPanic(t, func() {
			p.CanPeek(0)
		}, "Parser.CanPeek: range error")
		return nil
	}
	tokens := mockLexer()
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestPeek1
//
func TestPeek1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeekType(t, p, 1, T_ONE)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestPeek11
//
func TestPeek11(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeekType(t, p, 1, T_ONE)
		expectPeekType(t, p, 1, T_ONE)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestPeek12
//
func TestPeek12(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeekType(t, p, 1, T_ONE)
		expectPeekType(t, p, 2, T_TWO)
		return nil
	}
	tokens := mockLexer(T_ONE, T_TWO)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestPeekEmpty
//
func TestPeekEmpty(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		assertPanic(t, func() {
			p.Peek(1)
		}, "Parser.Peek: No AST available")
		return nil
	}
	tokens := mockLexer()
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestPeekRangeError
//
func TestPeekRangeError(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		assertPanic(t, func() {
			p.Peek(-1)
		}, "Parser.Peek: range error")
		assertPanic(t, func() {
			p.Peek(0)
		}, "Parser.Peek: range error")
		return nil
	}
	tokens := mockLexer()
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestHasNextTrue
//
func TestHasNextTrue(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectHasNext(t, p, true)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestHasNextFalse
//
func TestHasNextFalse(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectHasNext(t, p, false)
		return nil
	}
	tokens := mockLexer()
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestNext1
//
func TestNext1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeek(t, p, 1, T_ONE, "")
		expectNext(t, p, T_ONE, "")
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestNext2
//
func TestNext2(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeek(t, p, 1, T_ONE, "")
		expectNext(t, p, T_ONE, "")
		expectPeek(t, p, 1, T_TWO, "")
		expectNext(t, p, T_TWO, "")
		return nil
	}
	tokens := mockLexer(T_ONE, T_TWO)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestNextEmpty
//
func TestNextEmpty(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		assertPanic(t, func() {
			p.Next()
		}, "Parser.Next: No AST available")
		return nil
	}
	tokens := mockLexer()
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestNextEmit1
//
func TestNextEmit1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeek(t, p, 1, T_ONE, "")
		expectNext(t, p, T_ONE, "")
		p.Emit("T_ONE")
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsNext(t, emits, "T_ONE")
	expectEmitsHasNext(t, emits, false)
}

// TestNextEmit2
//
func TestNextEmit2(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeek(t, p, 1, T_ONE, "")
		expectNext(t, p, T_ONE, "")
		p.Emit("T_ONE")
		expectPeek(t, p, 1, T_TWO, "")
		expectNext(t, p, T_TWO, "")
		p.Emit("T_TWO")
		return nil
	}
	tokens := mockLexer(T_ONE, T_TWO)
	emits := Parse(tokens, fn)
	expectEmitsNext(t, emits, "T_ONE")
	expectEmitsNext(t, emits, "T_TWO")
	expectEmitsHasNext(t, emits, false)
}

// TestDiscard1
//
func TestDiscard1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeek(t, p, 1, T_ONE, "")
		expectNext(t, p, T_ONE, "")
		p.Discard()
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestDiscard2
//
func TestDiscard2(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeek(t, p, 1, T_ONE, "")
		expectNext(t, p, T_ONE, "")
		p.Emit("T_ONE")
		expectPeek(t, p, 1, T_TWO, "")
		expectNext(t, p, T_TWO, "")
		p.Discard()
		return nil
	}
	tokens := mockLexer(T_ONE, T_TWO)
	emits := Parse(tokens, fn)
	expectEmitsNext(t, emits, "T_ONE")
	expectEmitsHasNext(t, emits, false)
}

// TestDiscard3
//
func TestDiscard3(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeek(t, p, 1, T_ONE, "")
		expectNext(t, p, T_ONE, "")
		p.Discard()
		expectPeek(t, p, 1, T_TWO, "")
		expectNext(t, p, T_TWO, "")
		p.Emit("T_TWO")
		return nil
	}
	tokens := mockLexer(T_ONE, T_TWO)
	emits := Parse(tokens, fn)
	expectEmitsNext(t, emits, "T_TWO")
	expectEmitsHasNext(t, emits, false)
}

// TestEmitEOF1
//
func TestEmitEOF1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		p.EmitEOF()
		expectEOF(t, p)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestEmitEOF2
//
func TestEmitEOF2(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, T_ONE, "")
		p.EmitEOF()
		expectEOF(t, p)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestEmitEOF3
//
func TestEmitEOF3(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		p.Emit(nil)
		expectEOF(t, p)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestEmitAfterEOF
//
func TestEmitAfterEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, T_ONE, "")
		p.EmitEOF()
		expectEOF(t, p)
		p.Emit("T_ONE")
		return nil
	}
	tokens := mockLexer(T_ONE)
	assertPanic(t, func() {
		Parse(tokens, fn).Next()
	}, "Parser.Emit: No further emits allowed after EOF is emitted")
}

// TestCanPeekAfterEOF
//
func TestCanPeekAfterEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		p.EmitEOF()
		expectEOF(t, p)
		expectCanPeek(t, p, 1, false)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestPeekAfterEOF
//
func TestPeekAfterEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		p.EmitEOF()
		expectEOF(t, p)
		assertPanic(t, func() {
			p.Peek(1)
		}, "Parser.Peek: No tokens can be peeked after EOF is emitted")
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestHesNextAfterEOF
//
func TestHasNextAfterEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		p.EmitEOF()
		expectEOF(t, p)
		expectHasNext(t, p, false)
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestNextAfterEOF
//
func TestNextAfterEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		p.EmitEOF()
		expectEOF(t, p)
		assertPanic(t, func() {
			p.Next()
		}, "Parser.Next: No tokens can be consumed after EOF is emitted")
		return nil
	}
	tokens := mockLexer(T_ONE)
	emits := Parse(tokens, fn)
	expectEmitsHasNext(t, emits, false)
}

// TestDiscardAfterEOF
//
func TestDiscardAfterEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectNext(t, p, T_ONE, "")
		p.EmitEOF()
		expectEOF(t, p)
		p.Discard()
		return nil
	}
	tokens := mockLexer(T_ONE)
	assertPanic(t, func() {
		Parse(tokens, fn).Next()
	}, "Parser.Discard: No discards allowed after EOF is emitted")
}
