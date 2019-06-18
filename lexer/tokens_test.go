package lexer

import "testing"

// expectTokensNext
//
func expectTokensNext(t *testing.T, tokens *Tokens, typ TokenType, str string) {
	tok := tokens.Next()
	if tok.Type != typ {
		t.Errorf("Tokens.Next() Expecting Token.Type '%s', received '%s'", typ, tok.Type)
	}
	if tok.String != str {
		t.Errorf("Tokens.Next() Expecting Token.String '%s', received '%s'", str, tok.String)
	}
}

// expectTokensHasNext
//
func expectTokensHasNext(t *testing.T, tokens *Tokens, match bool) {
	if tokens.HasNext() != match {
		t.Errorf("Tokens.HasNext() expected to return %t", match)
	}
}

// // expectTokensEOF
// //
// func expectTokensEOF(t *testing.T, tokens *Tokens) {
// 	if tokens.HasNext() == true {
// 		t.Errorf("Tokens.HasNext() expected to return false")
// 	}
// }

// TestTokensHasNext1
//
func TestTokensHasNext1(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitType(T_START)
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, true)
	expectTokensNext(t, tokens, T_START, "")
	expectTokensHasNext(t, tokens, false)
}

// TestTokensHasNext2
//
func TestTokensHasNext2(t *testing.T) {
	fn := func(l *Lexer) LexerFn {
		l.EmitType(T_START)
		return nil
	}
	tokens := LexString(".", fn)
	expectTokensHasNext(t, tokens, true)
	expectTokensHasNext(t, tokens, true) // Call again, should hit cached 'next' value
	expectTokensNext(t, tokens, T_START, "")
	expectTokensHasNext(t, tokens, false)
}

// TestTokenEOF
//
func TestTokensEOF(t *testing.T) {
	tokens := LexString(".", nil)
	expectTokensHasNext(t, tokens, false)
}

// TestTokensNextAfterEOF
//
func TestTokensNextAfterEOF(t *testing.T) {
	tokens := LexString(".", nil)
	expectTokensHasNext(t, tokens, false)
	assertPanic(t, func() {
		tokens.Next()
	}, "Tokens.Next: No token available")
}