package lexer

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// Fn are user functions that scan runes and emit tokens.
// Functions are allowed to emit multiple tokens within a single call-back.
// The lexer executes functions in a continuous loop until either the function returns nil or emits an EOF token.
// Functions should return nil after emitting EOF, as no further interactions are allowed afterwards.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
//
type Fn func(*Lexer) Fn

// LexString initiates a lexer against the input string.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input string in an io.RuneReader, then calling LexRuneReader().
//
func LexString(input string, start Fn) token.Nexter {
	return LexRuneReader(strings.NewReader(input), start)
}

// LexRuneReader initiates a lexer against the input io.RuneReader.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// LexRuneReader is the primary lexer entrypoint. All others are convenience methods that delegate to here.
//
func LexRuneReader(input io.RuneReader, start Fn) token.Nexter {
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
func LexReader(input io.Reader, start Fn) token.Nexter {
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
func LexRunes(input []rune, start Fn) token.Nexter {
	return LexRuneReader(strings.NewReader(string(input)), start)
}

// LexBytes initiates a lexer against the input []byte.
// The returned token.Nexter can be used to retrieve emitted tokens.
// Invalid runes in the input will be silently ignored and will not be available within the lexer.
// The lexer will auto-emit EOF before exiting if it has not already been emitted.
// This is a convenience method, wrapping the input []byte in an io.RuneReader, then calling LexRuneReader().
//
func LexBytes(input []byte, start Fn) token.Nexter {
	return LexRuneReader(bytes.NewReader(input), start)
}

// Lexer is passed into your Lexer.Fn functions and provides methods to inspect runes and match them to tokens.
// When your Lexer.Fn is called, the lexer guarantees that `CanPeek(1) == true`, ensuring there is at least one rune to
// review/match.
//
type Lexer struct {
	input     io.RuneReader // Source of runes
	cache     *list.List    // Cache of fetched runes, including matched & peeked
	matchTail *list.Element // Points to last matched element in the cache, nil if no runes matched yet
	matchLen  int           // Len of match buffer.  Makes growPeek faster when no growth needed
	nextFn    Fn            // the next lexing function to enter
	output    *list.List    // Cache of emitted tokens ready for pickup by a parser
	eof       bool          // Has EOF been reached on the input reader? NOTE Peek buffer may still have runes in it
	eofOut    bool          // Has EOF been emitted to the output buffer?
	markerID  int           // Incremented after each emit/clear - used to validate markers
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

// Next matches and returns the next rune in the input.
// See CanPeek(1) to confirm if a rune is available.
// See Peek(1) to review the rune before consuming it.
// Panics if no rune available.
// Panics if EOF already emitted.
//
func (l *Lexer) Next() rune {
	// Nothing can be returned after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.Next: No runes can be matched after EOF is emitted")
	}
	if !l.growPeek(1) {
		panic("Lexer.Next: No rune available")
	}
	// Element guaranteed to exist
	//
	e := l.peekHead()
	l.matchTail = e // Match next rune into token
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

// EmitToken emits a token of the specified type, along with all of the matched runes.
// It is safe to emit TEof via this method.
// If the type is TEof, then all previously-matched runes are discarded and this is treated as EmitEOF().
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

// EmitType emits a token of the specified type, discarding all previously-matched runes.
// The emitted token will have a Text() value of "".
// It is safe to emit TEof via this method.
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

// EmitError Emits a token of type TLexErr with the specified err string as the token text.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
//
func (l *Lexer) EmitError(err string) {
	// Nothing can be emitted after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.EmitError: No further emits allowed after EOF is emitted")
	}
	l.clear(false)
	// TODO This is a tad kludgie - Think of a better way to inject a string into the standard emit flow.
	l.output.PushBack(newToken(TLexErr, err))
}

// EmitErrorf Emits a token of type TLexErr with the formatted err string as the token text.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
// This is a convenience method that simply sends the formatted string to EmitError().
//
func (l *Lexer) EmitErrorf(format string, args ...interface{}) {
	l.EmitError(fmt.Sprintf(format, args...))
}

// EmitEOF emits a token of type TokenEOF, discarding all previously-matched runes.
// You will likely never need to call this directly, as Lex will auto-emit EOF (TEof) before exiting,
// if not already emitted.
// No more reads to the underlying RuneReader will happen once EOF is emitted.
// No more runes can be matched once EOF is emitted.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
// This is a convenience method that simply calls EmitType(TEof).
//
func (l *Lexer) EmitEOF() {
	l.EmitType(TEof)
}

// Clear discards all previously-matched runes without emitting any tokens.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
//
func (l *Lexer) Clear() {
	// Nothing can be cleared after EOF emitted
	//
	if l.eofOut {
		panic("Lexer.Clear: No clears allowed after EOF is emitted")
	}
	l.clear(false)
}

// newLexer
//
func newLexer(reader io.RuneReader, start Fn) *Lexer {
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
			//
			case io.EOF:
				l.eof = true

			// NON-EOF Error
			//
			default:
				// For lack of a better plan, treat as EOF for now
				// TODO Think about how to handle non-EOF errors.
				// TODO Expose upstream?
				//
				log.Printf("non-EOF error returned from rune reader, treating as EOF: %v", err)
				l.eof = true
			}
		}
	}
	return true
}

// peekHead computes the peek buffer head as a function of the matchTail.
//
func (l *Lexer) peekHead() *list.Element {
	// If any matched runes
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
// If token.Type is TEof, emitText is ignored and treated as false.
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
	if TEof == t {
		// Clear the peek buffer, discarding matched runes
		//
		l.matchTail = nil
		l.matchLen = 0
		l.cache.Init()
		// Invalidate outstanding markers manually,
		// avoiding otherwise redundant call to clear()
		//
		l.markerID++ // TODO If it ever takes 2+ commands to invalidate markers, then turn into separate method.
		// Mark EOF
		//
		l.eof = true
		l.eofOut = true
		// Emit EOF token
		//
		l.output.PushBack(newToken(TEof, ""))
	} else {
		s := l.clear(emitText)

		l.output.PushBack(newToken(t, s))
	}
}

// clear discards the previously-matched runes, optionally returning them as a
// string.
// All outstanding markers are invalidated after this call.
//
func (l *Lexer) clear(returnText bool) string {
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
