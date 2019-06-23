package lexer

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// LexerFn are user functions that scan runes and emit tokens.
// Functions are allowed to emit multiple tokens within a single call-back.
// The lexer executes functions in a continuous loop until either the function returns nil or emits an EOF token.
// Functions should return nil after emitting EOF, as no further interactions are allowed afterwards.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
//
type LexerFn func(*Lexer) LexerFn

// LexString initiates a lexer against the input string.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input string in an io.RuneReader, then calling LexRuneReader().
//
func LexString(input string, start LexerFn) token.Nexter {
	return LexRuneReader(strings.NewReader(input), start)
}

// LexRuneReader initiates a lexer against the input io.RuneReader.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// LexRuneReader is the primary lexer entrypoint. All others are convenience methods that delegate to here.
//
func LexRuneReader(input io.RuneReader, start LexerFn) token.Nexter {
	l := newLexer(input, start)
	return &tokenNexter{lexer: l}
}

// LexReader initiates a lexer against the input io.Reader.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input io.Reader in an io.RuneReader, then calling LexRuneReader().
// If the provided reader already implements io.RuneReader, it is used without wrapping.
//
func LexReader(input io.Reader, start LexerFn) token.Nexter {
	var runeReader io.RuneReader
	if r, ok := input.(io.RuneReader); ok {
		runeReader = r
	} else {
		runeReader = bufio.NewReader(input)
	}
	return LexRuneReader(runeReader, start)
}

// LexRunes initiates a lexer against the input []rune.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input []rune in an io.RuneReader, then calling LexRuneReader().
//
func LexRunes(input []rune, start LexerFn) token.Nexter {
	return LexRuneReader(strings.NewReader(string(input)), start)
}

// LexBytes initiates a lexer against the input []byte.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input []byte in an io.RuneReader, then calling LexRuneReader().
//
func LexBytes(input []byte, start LexerFn) token.Nexter {
	return LexRuneReader(bytes.NewReader(input), start)
}

// Lexer is passed into your LexerFn functions and provides methods to inspect runes and match them to tokens.
// When your LexerFn is called, the lexer guarantees that `CanPeek(1) == true` so your function can safely
// inspect/consume the next rune in the input.
//
type Lexer struct {
	input     io.RuneReader // Source of runes
	cache     *list.List    // Cache of fetched runes, including matched & peeked
	matchTail *list.Element // Points to last matched element in the cache, nil if no runes matched yet
	matchLen  int           // Len of match buffer.  Makes growPeek faster when no growth needed
	nextFn    LexerFn       // the next lexing function to enter
	output    *list.List    // Cache of emitted tokens ready for pickup by a parser
	eof       bool          // Has EOF been reached on the input reader? NOTE Peek buffer may still have runes in it
	eofOut    bool          // Has EOF been emitted to the output buffer?
	markerID  int           // Incremented after each emit/discard - used to validate markers
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

// Next consumes and returns the next rune in the input.
// See CanPeek(1) to confirm if a rune is available.
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
	l.matchTail = e // Consume next rune into token
	l.matchLen++
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
	for n, e := 0, l.cache.Front(); n < l.matchLen; n, e = n+1, e.Next() {
		b.WriteRune(e.Value.(rune))
	}
	return b.String()
}

// EmitToken emits a token of the specified type, along with all of the consumed runes.
// It is safe to emit T_EOF via this method.
// If the type is T_EOF, then the consumed runes are discarded and this is treated as EmitEOF().
// All outstanding markers are invalidated after this call.
// See EmitEOF for more details on the effects of emitting EOF.
// Panics if EOF already emitted.
//
func (l *Lexer) EmitToken(t token.Type) {
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
func (l *Lexer) EmitType(t token.Type) {
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
	l.output.PushBack(newToken(T_LEX_ERR, err))
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

// newLexer
//
func newLexer(reader io.RuneReader, start LexerFn) *Lexer {
	l := &Lexer{
		input:     reader,
		cache:     list.New(),
		matchTail: nil,
		matchLen:  0,
		nextFn:    start,
		output:    list.New(),
		eof:       false,
		eofOut:    false,
		markerID:  0,
	}
	return l
}

// growPeek tries to ensure the peek buffer has Len() >= n, growing if needed, returning success or failure.
// n is 1-based.
//
func (l *Lexer) growPeek(n int) bool {
	// Grow to n
	//
	peekLen := l.cache.Len() - l.matchLen
	for peekLen < n {
		// Nothing to do if EOF reached already
		//
		if l.eof {
			return false
		}
		// Fetch next rune from input
		//
		r, size, err := l.input.ReadRune()
		// Process any returned rune, regardless of err
		//
		if size > 0 {
			// Skip rune errors
			// TODO Log rune errors
			//
			if r != utf8.RuneError {
				// Add rune to peek buffer
				//
				l.cache.PushBack(r)
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

// peekHead computes the peek buffer head as a function of the matchTail.
//
func (l *Lexer) peekHead() *list.Element {
	// If any consumed runes
	//
	if l.matchLen > 0 {
		// Peek buffer starts after token
		//
		// assert(l.matchTail != nil)
		return l.matchTail.Next()
	}
	// Its ALL the peek buffer
	//
	return l.cache.Front()
}

// emit Emits a Token, optionally including the matched text.
// If token.Type is T_EOF, emitExt is ignored and treated as false.
// Panics if EOF already emitted.
//
func (l *Lexer) emit(t token.Type, emitText bool) {
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
		l.matchTail = nil
		l.matchLen = 0
		l.cache.Init()
		// Invalidate outstanding markers manually,
		// avoiding otherwise redundant call to consume()
		//
		l.markerID++ // TODO If it ever takes 2+ commands to invalidate markers, then turn into separate method.
		// Mark EOF
		//
		l.eof = true
		l.eofOut = true
		// Emit EOF token
		//
		l.output.PushBack(newToken(T_EOF, ""))
	} else {
		s := l.consume(emitText)

		l.output.PushBack(newToken(t, s))
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
		for l.matchLen > 0 {
			e := l.cache.Front()
			b.WriteRune(e.Value.(rune))
			l.cache.Remove(e)
			l.matchLen--
		}
		s = b.String()
	} else {
		// Discard runes
		//
		for l.matchLen > 0 {
			l.cache.Remove(l.cache.Front())
			l.matchLen--
		}
		s = ""
	}
	l.markerID++ // Invalidate outstanding markers

	return s
}
