package lexer

// Tokens is returned by the various Lex* functions and provides methods to retrieve tokens emitted from the lexer.
//
// Tokens implements a basic iterator pattern with HasNext() and Next() methods.
//
type Tokens struct {
	lexer *Lexer
	next  *Token
	eof   bool
}

// Next Retrieves the next token from the lexer.
// See HasNext() to determine if any tokens are available.
// Panics if HasNext() returns false.
//
func (t *Tokens) Next() *Token {
	// We double check for saved next to maybe avoid the call
	//
	if t.next == nil && t.HasNext() == false {
		panic("Tokens.Next: No token available")
	}
	tok := t.next
	t.next = nil
	return tok
}

// HasNext confirms if there are tokens available.
// This method initiates calls to LexerFn functions and is the primary entry point for retrieving tokens from the lexer.
// If it returns true, you can safely call Next() to retrieve the next token.
// If it returns false, EOF has been reached and calling Next() will generate a panic.
//
func (t *Tokens) HasNext() bool {
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
	token := t.lexer.tokens.Remove(t.lexer.tokens.Front()).(*Token)
	// Is the token EOF?
	//
	if token.eof() {
		// Mark EOF, discarding the token
		//
		t.eof = true
		return false
	}
	// Store the token for pickup
	//
	t.next = token
	return true
}
