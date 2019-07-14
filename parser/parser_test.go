package parser

import (
	"errors"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// Define tokens used in various tests
//
const (
	TStart token.Type = iota
	TOne
	TTwo
	TThree
)

// mockToken creates a token.Token from a token.Type
//
type mockToken struct {
	typ token.Type
}

func (t *mockToken) Type() token.Type {
	return t.typ
}
func (t *mockToken) Value() string {
	return ""
}
func (t *mockToken) Line() int {
	return -1
}
func (t *mockToken) Column() int {
	return -1
}

// mockNexter creates a token.Nexter from a list of token.Type
//
type mockNexter struct {
	tokens []token.Type
	i      int
	err    error
}

func (n *mockNexter) Next() (token.Token, error) {
	if n.err != nil {
		return nil, n.err
	}
	if n.i >= cap(n.tokens) {
		return nil, io.EOF
	}
	t := n.tokens[n.i]
	n.i++
	return &mockToken{typ: t}, nil
}

// mockLexer
//
func mockLexer(tokens ...token.Type) token.Nexter {
	return &mockNexter{tokens: tokens}
}

// mockLexerErr
//
func mockLexerErr(err error) token.Nexter {
	return &mockNexter{err: err}
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
func expectPeekType(t *testing.T, p *Parser, peek int, match token.Type) {
	if typ := p.PeekType(peek); typ != match {
		t.Errorf("Parser.PeekType(%d) expecting Token.Type '%d', received '%d'", peek, match, typ)
	}
}

// expectPeek
//
func expectPeek(t *testing.T, p *Parser, peek int, typ token.Type, value string) {
	tok := p.Peek(peek)
	if tok.Type() != typ {
		t.Errorf("Parser.Peek(%d) expecting Token.Type '%d', received '%d'", peek, typ, tok.Type())
	}
	if tok.Value() != value {
		t.Errorf("Parser.Peek(%d) expecting Token.String '%s', received '%s'", peek, value, tok.Value())
	}
}

// expectNext
//
func expectNext(t *testing.T, p *Parser, typ token.Type, value string) {
	tok := p.Next()
	if tok.Type() != typ {
		t.Errorf("Parser.Next() expecting Token.Type '%d', received '%d'", typ, tok.Type())
	}
	if tok.Value() != value {
		t.Errorf("Parser.Next() expecting Token.String '%s', received '%s'", value, tok.Value())
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
	nexter := Parse(tokens, nil)
	expectNexterEOF(t, nexter)
}

func TestParserFnSkipedWhenNoCanPeek(t *testing.T) {
	fn := func(p *Parser) Fn {
		t.Error("Parser should not call Parser.Fn when CanPeek(1) == false")
		return nil
	}
	tokens := mockLexer()
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestEmit
//
func TestEmit(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectNext(t, p, TStart, "")
		p.Emit("TStart")
		return nil
	}
	tokens := mockLexer(TStart)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "TStart")
	expectNexterEOF(t, nexter)
}

// TestCanPeek
//
func TestCanPeek(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectCanPeek(t, p, 1, true)

		expectPeekType(t, p, 1, TOne)

		expectCanPeek(t, p, 2, true)

		expectPeekType(t, p, 2, TTwo)

		expectCanPeek(t, p, 3, true)

		expectPeekType(t, p, 3, TThree)

		return nil
	}
	tokens := mockLexer(TOne, TTwo, TThree)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestCanPeekPastEOF
//
func TestCanPeekPastEOF(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectCanPeek(t, p, 4, false)

		expectCanPeek(t, p, 3, true)

		expectPeekType(t, p, 3, TThree)

		expectCanPeek(t, p, 2, true)

		expectPeekType(t, p, 2, TTwo)

		expectCanPeek(t, p, 1, true)

		expectPeekType(t, p, 1, TOne)

		return nil
	}
	tokens := mockLexer(TOne, TTwo, TThree)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestCanPeekRangeError
//
func TestCanPeekRangeError(t *testing.T) {
	fn := func(p *Parser) Fn {
		assertPanic(t, func() {
			p.CanPeek(-1)
		}, "Parser.CanPeek: range error")
		assertPanic(t, func() {
			p.CanPeek(0)
		}, "Parser.CanPeek: range error")
		return nil
	}
	tokens := mockLexer()
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestPeek1
//
func TestPeek1(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeekType(t, p, 1, TOne)
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestPeek11
//
func TestPeek11(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeekType(t, p, 1, TOne)
		expectPeekType(t, p, 1, TOne)
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestPeek12
//
func TestPeek12(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeekType(t, p, 1, TOne)
		expectPeekType(t, p, 2, TTwo)
		return nil
	}
	tokens := mockLexer(TOne, TTwo)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestPeekEmpty
//
func TestPeekEmpty(t *testing.T) {
	fn := func(p *Parser) Fn {
		assertPanic(t, func() {
			p.Peek(1)
		}, "Parser.Peek: No AST available")
		return nil
	}
	tokens := mockLexer()
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestPeekRangeError
//
func TestPeekRangeError(t *testing.T) {
	fn := func(p *Parser) Fn {
		assertPanic(t, func() {
			p.Peek(-1)
		}, "Parser.Peek: range error")
		assertPanic(t, func() {
			p.Peek(0)
		}, "Parser.Peek: range error")
		return nil
	}
	tokens := mockLexer()
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestNext1
//
func TestNext1(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeek(t, p, 1, TOne, "")
		expectNext(t, p, TOne, "")
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestNext2
//
func TestNext2(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeek(t, p, 1, TOne, "")
		expectNext(t, p, TOne, "")
		expectPeek(t, p, 1, TTwo, "")
		expectNext(t, p, TTwo, "")
		return nil
	}
	tokens := mockLexer(TOne, TTwo)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestNextEmpty
//
func TestNextEmpty(t *testing.T) {
	fn := func(p *Parser) Fn {
		assertPanic(t, func() {
			p.Next()
		}, "Parser.Next: No AST available")
		return nil
	}
	tokens := mockLexer()
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestNextEmit1
//
func TestNextEmit1(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeek(t, p, 1, TOne, "")
		expectNext(t, p, TOne, "")
		p.Emit("TOne")
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "TOne")
	expectNexterEOF(t, nexter)
}

// TestNextEmit2
//
func TestNextEmit2(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeek(t, p, 1, TOne, "")
		expectNext(t, p, TOne, "")
		p.Emit("TOne")
		expectPeek(t, p, 1, TTwo, "")
		expectNext(t, p, TTwo, "")
		p.Emit("TTwo")
		return nil
	}
	tokens := mockLexer(TOne, TTwo)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "TOne")
	expectNexterNext(t, nexter, "TTwo")
	expectNexterEOF(t, nexter)
}

// TestClear1
//
func TestClear1(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeek(t, p, 1, TOne, "")
		expectNext(t, p, TOne, "")
		p.Clear()
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestClear2
//
func TestClear2(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeek(t, p, 1, TOne, "")
		expectNext(t, p, TOne, "")
		p.Emit("TOne")
		expectPeek(t, p, 1, TTwo, "")
		expectNext(t, p, TTwo, "")
		p.Clear()
		return nil
	}
	tokens := mockLexer(TOne, TTwo)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "TOne")
	expectNexterEOF(t, nexter)
}

// TestClear3
//
func TestClear3(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectPeek(t, p, 1, TOne, "")
		expectNext(t, p, TOne, "")
		p.Clear()
		expectPeek(t, p, 1, TTwo, "")
		expectNext(t, p, TTwo, "")
		p.Emit("TTwo")
		return nil
	}
	tokens := mockLexer(TOne, TTwo)
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "TTwo")
	expectNexterEOF(t, nexter)
}

// TestEmitEOF1
//
func TestEmitEOF1(t *testing.T) {
	fn := func(p *Parser) Fn {
		p.EmitEOF() // Emits EOF explicitly
		expectEOF(t, p)
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestEmitEOF2
//
func TestEmitEOF2(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectNext(t, p, TOne, "")
		p.EmitEOF()
		expectEOF(t, p)
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestEmitEOF3
//
func TestEmitEOF3(t *testing.T) {
	fn := func(p *Parser) Fn {
		p.Emit(nil) // Emits nil for EOF
		expectEOF(t, p)
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestEmitAfterEOF
//
func TestEmitAfterEOF(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectNext(t, p, TOne, "")
		p.EmitEOF()
		expectEOF(t, p)
		p.Emit("TOne")
		return nil
	}
	tokens := mockLexer(TOne)
	assertPanic(t, func() {
		_, _ = Parse(tokens, fn).Next()
	}, "Parser.Emit: No further emits allowed after EOF is emitted")
}

// TestCanPeekAfterEOF
//
func TestCanPeekAfterEOF(t *testing.T) {
	fn := func(p *Parser) Fn {
		p.EmitEOF()
		expectEOF(t, p)
		expectCanPeek(t, p, 1, false)
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestPeekAfterEOF
//
func TestPeekAfterEOF(t *testing.T) {
	fn := func(p *Parser) Fn {
		p.EmitEOF()
		expectEOF(t, p)
		assertPanic(t, func() {
			p.Peek(1)
		}, "Parser.Peek: No tokens can be peeked after EOF is emitted")
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestNextAfterEOF
//
func TestNextAfterEOF(t *testing.T) {
	fn := func(p *Parser) Fn {
		p.EmitEOF()
		expectEOF(t, p)
		assertPanic(t, func() {
			p.Next()
		}, "Parser.Next: No tokens can be matched after EOF is emitted")
		return nil
	}
	tokens := mockLexer(TOne)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestClearAfterEOF
//
func TestClearAfterEOF(t *testing.T) {
	fn := func(p *Parser) Fn {
		expectNext(t, p, TOne, "")
		p.EmitEOF()
		expectEOF(t, p)
		p.Clear()
		return nil
	}
	tokens := mockLexer(TOne)
	assertPanic(t, func() {
		_, _ = Parse(tokens, fn).Next()
	}, "Parser.Clear: No clears allowed after EOF is emitted")
}

// TestTokenNexterNonEOFError should log an error but otherwise behave as EOF
//
func TestTokenNexterNonEOFError(t *testing.T) {
	sb := &strings.Builder{}
	log.SetFlags(0)
	log.SetOutput(sb)
	fn := func(p *Parser) Fn {
		p.EmitEOF() // Emits EOF explicitly
		expectEOF(t, p)
		return nil
	}
	tokens := mockLexerErr(errors.New("test Error"))
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
	if log := sb.String(); log != "non-EOF error returned from lexer, treating as EOF: test Error\n" {
		t.Errorf("Parser.growPeek received wrong log message: '%s'", log)
	}
}
