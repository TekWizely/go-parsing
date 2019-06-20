package lexer

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// LexerFn are user functions that scan runes and emit tokens.
// Functions are allowed to emit multiple tokens within a single call-back.
// The lexer executes functions in a continuous loop until either the function returns nil or emits an EOF token.
// Functions should return nil after emitting EOF, as no further interactions are allowed afterwards.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
//
type LexerFn func(*Lexer) LexerFn

// LexString initiates a lexer against the input string.
// The returned TokenNexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input string in an io.RuneReader, then calling LexRuneReader().
//
func LexString(input string, start LexerFn) TokenNexter {
	return LexRuneReader(strings.NewReader(input), start)
}

// LexRuneReader initiates a lexer against the input io.RuneReader.
// The returned TokenNexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// LexRuneReader is the primary lexer entrypoint. All others are convenience methods that delegate to here.
//
func LexRuneReader(input io.RuneReader, start LexerFn) TokenNexter {
	l := newLexer(input, start)
	return &tokenNexter{lexer: l}
}

// LexReader initiates a lexer against the input io.Reader.
// The returned TokenNexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input io.Reader in an io.RuneReader, then calling LexRuneReader().
// If the provided reader already implements io.RuneReader, it is used without wrapping.
//
func LexReader(input io.Reader, start LexerFn) TokenNexter {
	var runeReader io.RuneReader
	if r, ok := input.(io.RuneReader); ok {
		runeReader = r
	} else {
		runeReader = bufio.NewReader(input)
	}
	return LexRuneReader(runeReader, start)
}

// LexRunes initiates a lexer against the input []rune.
// The returned TokenNexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input []rune in an io.RuneReader, then calling LexRuneReader().
//
func LexRunes(input []rune, start LexerFn) TokenNexter {
	return LexRuneReader(strings.NewReader(string(input)), start)
}

// LexBytes initiates a lexer against the input []byte.
// The returned TokenNexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input []byte in an io.RuneReader, then calling LexRuneReader().
//
func LexBytes(input []byte, start LexerFn) TokenNexter {
	return LexRuneReader(bytes.NewReader(input), start)
}

// Lexer is passed into your LexerFn functions and provides methods to inspect runes and match them to tokens.
// When your LexerFn is called, the lexer guarantees that `HasNext() == true` so your function can safely
// inspect/consume the next rune in the input.
//
type Lexer struct {
	reader    io.RuneReader // reader
	runes     *list.List    // working runes (token + look-ahead)
	tokenTail *list.Element // Points to last element of token, nil if token is empty
	tokenLen  int           // Len of peek buffer.  Makes growPeek faster when no growth needed
	nextFn    LexerFn       // the next lexing function to enter
	tokens    *list.List    // Cache of emitted tokens ready for pickup by a parser
	eof       bool          // Has EOF been reached on the input reader? NOTE Peek buffer may still have runes in it
	eofOut    bool          // Has EOF been emitted to the output buffer?
	markerId  int           // Incremented after each emit/discard - used to validate markers
}

// CanPeek confirms if the requested number of runes are available in the peek buffer.
// n is 1-based.
// If CanPeek returns true, you can safely Peek for values up to, and including, n.
// Returns false if EOF already emitted.
// Panics if n < 1.
//
func (l *Lexer) CanPeek(n int) bool {
	if n < 1 {
		panic("Lexer.CanPeek: range error")
	}
	// Nothing can be peeked after EOF emitted
	//
	if l.eofOut {
		return false
	}
	return l.growPeek(n)
}

// Peek allows you to look ahead at runes without consuming them.
// n is 1-based.
// See CanPeek to confirm a minimum number of runes are available in the peek buffer.
// Panics if n < 1.
// Panics if nth rune not available.
// Panics if EOF already emitted.
//
func (l *Lexer) Peek(n int) rune {
	if n < 1 {
		panic("Lexer.Peek: range error")
	}
	// Nothing can be peeked after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.Peek: No runes can be peeked after EOF is emitted")
	}
	if !l.growPeek(n) {
		panic("Lexer.Peek: No rune available")
	}
	// Elements guaranteed to exist
	//
	e := l.peekHead() // 1st element
	for ; n > 1; n-- {
		e = e.Next()
	}
	return e.Value.(rune)
}

// HasNext confirms if a rune is available to consume.
// If HasNext returns true, you can safely call Next to consume and return the rune.
// Returns false if EOF already emitted.
// This is a convenience method and simply calls CanPeek(1).
//
func (l *Lexer) HasNext() bool {
	return l.CanPeek(1)
}

// Next consumes and returns the next rune in the input.
// See CanPeek(1) or HasNext() to confirm if a rune is available.
// See Peek(1) to review the rune before consuming it.
// Panics if no rune available.
// Panics if EOF already emitted.
//
func (l *Lexer) Next() rune {
	// Nothing can be returned after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.Next: No runes can be consumed after EOF is emitted")
	}
	if !l.growPeek(1) {
		panic("Lexer.Next: No rune available")
	}
	// Element guaranteed to exist
	//
	e := l.peekHead()
	l.tokenTail = e // Consume next rune into token
	l.tokenLen++
	return e.Value.(rune)
}

// PeekToken allows you to inspect the currently matched rune sequence.
// The value is returned as a string, same as EmitToken() would provide.
// Panics if EOF already emitted.
//
func (l *Lexer) PeekToken() string {
	// Nothing can be peeked after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.PeekToken: No token peeks allowed after EOF is emitted")
	}
	b := &strings.Builder{}
	for n, e := 0, l.runes.Front(); n < l.tokenLen; n, e = n+1, e.Next() {
		b.WriteRune(e.Value.(rune))
	}
	return b.String()
}

// PeekTokenRunes allows you to inspect the currently matched rune sequence as a rune array ( []rune ).
// Panics if EOF already emitted.
// This is a convenience method and simply executes return []rune(l.PeekToken()).
//
func (l *Lexer) PeekTokenRunes() []rune {
	return []rune(l.PeekToken())
}

// EmitToken emits a token of the specified type, along with all of the consumed runes.
// It is safe to emit T_EOF via this method.
// If the type is T_EOF, then the consumed runes are discarded and this is treated as EmitEOF().
// All outstanding markers are invalidated after this call.
// See EmitEOF for more details on the effects of emitting EOF.
// Panics if EOF already emitted.
//
func (l *Lexer) EmitToken(t TokenType) {
	// Nothing can be emitted after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.EmitToken: No further emits allowed after EOF is emitted")
	}
	l.emit(t, true)
}

// EmitType emits a token of the specified type, discarding consumed runes.
// The emitted token will have a Text() value of "".
// It is safe to emit T_EOF via this method.
// All outstanding markers are invalidated after this call.
// See EmitEOF for more details on the effects of emitting EOF.
// Panics if EOF already emitted.
//
func (l *Lexer) EmitType(t TokenType) {
	// Nothing can be emitted after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.EmitType: No further emits allowed after EOF is emitted")
	}
	l.emit(t, false)
}

// EmitErrorf Emits a token of type T_LEX_ERR with the specified err string as the token text.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
//
func (l *Lexer) EmitError(err string) {
	// Nothing can be emitted after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.EmitError: No further emits allowed after EOF is emitted")
	}
	l.consume(false)
	// TODO This is a tad kludgie - Think of a better way to inject a string into the standard emit flow.
	l.tokens.PushBack(newToken(T_LEX_ERR, err))
}

// EmitErrorf Emits a token of type T_LEX_ERR with the formatted err string as the token text.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
// This is a convenience method that simply sends the formatted string to EmitError().
//
func (l *Lexer) EmitErrorf(format string, args ...interface{}) {
	l.EmitError(fmt.Sprintf(format, args...))
}

// EmitEOF emits a token of type TokenEOF, discarding consumed runes.
// You will likely never need to call this directly, as Lex will auto-emit EOF (T_EOF) before exiting,
// if not already emitted.
// No more reads to the underlying RuneReader will happen once EOF is emitted.
// No more runes can be consumed once EOF is emitted.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
// This is a convenience method that simply calls EmitType(T_EOF).
//
func (l *Lexer) EmitEOF() {
	l.EmitType(T_EOF)
}

// DiscardToken discards the consumed runes without emitting any tokens.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
//
func (l *Lexer) DiscardToken() {
	// Nothing can be discarded after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.Discard: No discards allowed after EOF is emitted")
	}
	l.consume(false)
}

// TODO Remove this after API settles down.
// DFARRELL - Not needed as Lexer ensures `HasNext()` before calling LexerFn and
//            `HasNext()` / `CanPeek()` are better tools to use from within LexerFn.
// // EOF returns true if the peek buffer is empty AND the input has reached EOF.
// // This will return false if there are any runes remaining in the peek buffer.
// // See CanPeek() / HasNext() to confirm if runes are available to peek / consume.
// //
// func (l *Lexer) EOF() bool {
// 	return l.eof && l.runes.Len() == l.tokenLen
// }

// newLexer
//
func newLexer(reader io.RuneReader, start LexerFn) *Lexer {
	l := &Lexer{
		reader:    reader,
		runes:     list.New(),
		tokenTail: nil,
		tokenLen:  0,
		nextFn:    start,
		tokens:    list.New(),
		eof:       false,
		eofOut:    false,
		markerId:  0,
	}
	return l
}

// growPeek tries to ensure the peek buffer has Len() >= n, growing if needed, returning success or failure.
// n is 1-based.
//
func (l *Lexer) growPeek(n int) bool {
	// Grow to n
	//
	peekLen := l.runes.Len() - l.tokenLen
	for peekLen < n {
		// Nothing to do if EOF reached already
		//
		if l.eof {
			return false
		}
		// Fetch next rune from input
		//
		r, size, err := l.reader.ReadRune()
		// Process any returned rune, regardless of err
		//
		if size > 0 {
			// Skip rune errors
			// TODO Log rune errors
			//
			if r != utf8.RuneError {
				// Add rune to peek buffer
				//
				l.runes.PushBack(r)
				peekLen++
			}
		}
		// If there was an error, process it now
		//
		if err != nil {
			switch err {
			// EOF Error
			// Treat NoProgress as EOF for now
			// TODO Decide if ErrNoProgress should be treated as Non-EOF error.
			//
			case io.EOF, io.ErrNoProgress:
				l.eof = true

			// NON-EOF Error
			//
			default:
				// For lack of a better plan, treat as EOF for now
				// TODO Think about how to handle non-EOF errors.
				// TODO Log error.
				// TODO Expose upstream?
				//
				l.eof = true
			}
		}
	}
	return true
}

// peekHead computes the peek buffer head as a function of the tokenTail.
//
func (l *Lexer) peekHead() *list.Element {
	// If any consumed runes
	//
	if l.tokenLen > 0 {
		// Peek buffer starts after token
		//
		// assert(l.tokenTail != nil)
		return l.tokenTail.Next()
	}
	// Its ALL the peek buffer
	//
	return l.runes.Front()
}

// emit Emits a Token, optionally including the matched text.
// If TokenType is T_EOF, emitExt is ignored and treated as false.
// Panics if EOF already emitted.
//
func (l *Lexer) emit(t TokenType, emitText bool) {
	// TODO Current tests show this will never be called. Maybe uncomment this once in awhile to confirm :)
	// // Nothing can be emitted after EOF
	// // NOTE: This check is a fail-safe and will likely never hit as all public methods check/panic explicitly.
	// //
	// if l.eofOut {
	// 	panic("Lexer: No further emits allowed after EOF is emitted")
	// }

	// If emitting EOF
	//
	if T_EOF == t {
		// Clear the peek buffer, discarding consumed runes
		//
		l.tokenTail = nil
		l.tokenLen = 0
		l.runes.Init()
		// Invalidate outstanding markers manually,
		// avoiding otherwise redundant call to consume()
		//
		l.markerId++ // TODO If it ever takes 2+ commands to invalidate markers, then turn into separate method.
		// Mark EOF
		//
		l.eof = true
		l.eofOut = true
		// Emit EOF token
		//
		l.tokens.PushBack(newToken(T_EOF, ""))
	} else {
		s := l.consume(emitText)

		l.tokens.PushBack(newToken(t, s))
	}
}

// consume consumes the matched token, optionally returning the token text.
// All outstanding markers are invalidated after this call.
//
func (l *Lexer) consume(returnText bool) string {
	var s string
	if returnText {
		// Build the token into a string
		//
		b := &strings.Builder{}
		for l.tokenLen > 0 {
			e := l.runes.Front()
			b.WriteRune(e.Value.(rune))
			l.runes.Remove(e)
			l.tokenLen--
		}
		s = b.String()
	} else {
		// Discard runes
		//
		for l.tokenLen > 0 {
			l.runes.Remove(l.runes.Front())
			l.tokenLen--
		}
		s = ""
	}
	l.markerId++ // Invalidate outstanding markers

	return s
}
