package parser

import (
	"io"
	"testing"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// Define tokens used in various tests
//
const (
	T_START token.Type = iota
	T_ONE
	T_TWO
	T_THREE
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

// mockNexter creates a token.Nexter from a list of token.Type
//
type mockNexter struct {
	tokens []token.Type
	i      int
}

func (n *mockNexter) Next() (token.Token, error) {
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
	fn := func(p *Parser) ParserFn {
		t.Error("Parser should not call ParserFn when CanPeek(1) == false")
		return nil
	}
	tokens := mockLexer()
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_START")
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestPeek1
//
func TestPeek1(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		expectPeekType(t, p, 1, T_ONE)
		return nil
	}
	tokens := mockLexer(T_ONE)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_ONE")
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_ONE")
	expectNexterNext(t, nexter, "T_TWO")
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_ONE")
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterNext(t, nexter, "T_TWO")
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
		_, _ = Parse(tokens, fn).Next()
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
}

// TestNextAfterEOF
//
func TestNextAfterEOF(t *testing.T) {
	fn := func(p *Parser) ParserFn {
		p.EmitEOF()
		expectEOF(t, p)
		assertPanic(t, func() {
			p.Next()
		}, "Parser.Next: No tokens can be matched after EOF is emitted")
		return nil
	}
	tokens := mockLexer(T_ONE)
	nexter := Parse(tokens, fn)
	expectNexterEOF(t, nexter)
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
		_, _ = Parse(tokens, fn).Next()
	}, "Parser.Discard: No discards allowed after EOF is emitted")
}
