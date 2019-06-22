package lexer

import (
	"errors"
	"io"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// tokenNexter is the internal structure that backs the lexer's token.Nexter.
//
type tokenNexter struct {
	lexer *Lexer
	next  token.Token
	eof   bool
}

// Next implements token.Nexter.Next().
// We build on the previous HasNext/Next impl to keep changes minimal.
//
func (t *tokenNexter) Next() (token.Token, error) {
	if !t.hasNext() {
		return nil, io.EOF
	}
	tok := t.next
	t.next = nil
	// Error?
	//
	if tok.Type() == T_LEX_ERR {
		return nil, errors.New(tok.Value())
	}
	return tok, nil
}

// hasNext Initiates calls to LexerFn functions and is the primary entry point for retrieving tokens from the lexer.
//
func (t *tokenNexter) hasNext() bool {
	// If token previously fetched, return now
	//
	if t.next != nil {
		return true
	}
	// Nothing to do once EOF reached
	//
	if t.eof {
		return false
	}
	// If no tokens available, try to fetch some.
	//
	for t.lexer.tokens.Len() == 0 {
		// Anyone to call?
		// Anything to scan?
		//
		if t.lexer.nextFn != nil && t.lexer.HasNext() {
			t.lexer.nextFn = t.lexer.nextFn(t.lexer)
		} else {
			// Lexer Terminated or input at EOF, let's clean up.
			// If EOF was never emitted, then emit it now.
			//
			if t.lexer.eofOut == false {
				t.lexer.EmitEOF()
			}
		}
	}
	// Consume the token.
	// We'll either cache it or discard it.
	//
	tok := t.lexer.tokens.Remove(t.lexer.tokens.Front()).(*_token)
	// Is the token EOF?
	//
	if tok.eof() {
		// Mark EOF, discarding the token
		//
		t.eof = true
		return false
	}
	// Store the token for pickup
	//
	t.next = tok
	return true
}
