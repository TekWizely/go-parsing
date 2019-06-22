package parser

import (
	"container/list"
	"io"

	"github.com/tekwizely/go-parsing/lexer/token"
)

// ParserFn are user functions that scan tokens and emit ASTs.
// Functions are allowed to emit multiple ASTs within a single call-back.
// The parser executes functions in a continuous loop until either the function returns nil or emits an EOF value.
// Functions should return nil after emitting EOF, as no further interactions are allowed afterwards.
// The parser will auto-emit EOF before exiting if it has not already been emitted.
//
type ParserFn func(*Parser) ParserFn

// Parse initiates a parser against the input token stream.
// The returned ASTNexter can be used to retrieve emitted ASTs.
// The parser will auto-emit EOF before exiting it if has not already been emitted.
//
func Parse(tokens token.Nexter, start ParserFn) ASTNexter {
	p := newParser(tokens, start)
	return &astNexter{parser: p}
}

// Parser is passed into your ParserFn functions and provides methods to inspect tokens and emit ASTs.
// When your ParserFn is called, the parser guarantees that 'HasNext == true` so your function can safely
// inspect/consume the next token in the input.
//
type Parser struct {
	tokens    token.Nexter  // Source for lexer tokens
	cache     *list.List    // Cache of fetched lexer tokens ready for pickup by the parser
	matchTail *list.Element // Points to last element of consumed tokens, nil if no tokens consumed yet
	matchLen  int           // Len of peek buffer.  Makes growPeek faster when no growth needed
	nextFn    ParserFn      // the next parsing function to enter
	emits     *list.List    // Cache of emitted ASTs ready for pickup
	eof       bool          // Has EOF been reached on the input tokens? NOTE Peek buffer may still have tokens in it
	eofOut    bool          // Has EOF been emitted to the output buffer?
	markerId  int           // Incremented after each emit/discard - used to validate markers
}

// CanPeek confirms if the requested number of tokens are available in the peek buffer.
// n is 1-based.
// If CanPeek returns true, you can safely Peek for values up to, and including, n.
// Returns false if EOF already emitted.
// Panics if n < 1.
//
func (p *Parser) CanPeek(n int) bool {
	if n < 1 {
		panic("Parser.CanPeek: range error")
	}
	// Nothing can be peeked after EOF emitted
	//
	if p.eofOut {
		return false
	}
	return p.growPeek(n)
}

// Peek allows you to look ahead at tokens without consuming them.
// n is 1-based.
// See CanPeek to confirm a minimum number of tokens are available in the peek buffer.
// Panics if n < 1.
// Panics if nth token not available.
// Panics if EOF already emitted.
//
func (p *Parser) Peek(n int) token.Token {
	if n < 1 {
		panic("Parser.Peek: range error")
	}
	// Nothing can be peeked after EOF
	//
	if p.eofOut {
		panic("Parser.Peek: No tokens can be peeked after EOF is emitted")
	}
	if !p.growPeek(n) {
		panic("Parser.Peek: No token available")
	}
	// Elements guaranteed to exist
	//
	e := p.peekHead() // 1st element
	for ; n > 1; n-- {
		e = e.Next()
	}
	return e.Value.(token.Token)
}

// PeekType allows you to look ahead at token types without consuming them.
// n is 1-based.
// See CanPeek to confirm a minimum number of tokens are available in the peek buffer.
// Panics if n < 1.
// Panics if nth token not available.
// Panics if EOF already emitted.
// This is mostly a convenience method that calls Peek(n), returning the token type.
//
func (p *Parser) PeekType(n int) token.Type {
	return p.Peek(n).Type()
}

// HasNext confirms if a token is available to consume.
// If HasNext returns true, you can safely call Next to consume and return the token.
// Returns false if EOF already emitted.
// This is a convenience method and simply calls CanPeek(1).
//
func (p *Parser) HasNext() bool {
	return p.CanPeek(1)
}

// Next consumes and returns the next token in the input.
// See CanPeek(1) to confirm if a token is available.
// See Peek(1) and PeekType(1) to review the token before consuming it.
// Panics if no token available.
// Panics if EOF already emitted.
//
func (p *Parser) Next() token.Token {
	// Nothing can be peeked after EOF
	//
	if p.eofOut {
		panic("Parser.Next: No tokens can be consumed after EOF is emitted")
	}
	if !p.growPeek(1) { // Cache next emit. 1-based
		panic("Parser.Next: No token available")
	}
	// Element guaranteed to exist
	//
	e := p.peekHead()
	p.matchTail = e // Consume peek into token
	p.matchLen++
	return e.Value.(token.Token)
}

// Emit emits an AST.
// Consumed tokens are discarded.
// It is safe to emit nil via this method.
// If the emit value is nil, then this is treated as EmitEOF().
// All outstanding markers are invalidated after this call.
// See EmitEOF for more details on the effects of emitting EOF.
// Panics if EOF already emitted.
//
func (p *Parser) Emit(ast interface{}) {
	// Nothing can be emitted after EOF emitted
	//
	if p.eofOut {
		panic("Parser.Emit: No further emits allowed after EOF is emitted")
	}
	p.emit(ast)
}

// EmitEOF emits a nil, discarding consumed tokens.
// You will likely never need to call this directly, as Parse will auto-emit EOF (nil) before exiting,
// if not already emitted.
// No more reads to the underlying Lexer will happen once EOF is emitted.
// No more tokens can be consumed once EOF is emitted.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
// This is a convenience method that simply calls Emit(nil).
//
func (p *Parser) EmitEOF() {
	p.Emit(nil)
}

// Discard discards the consumed tokens without emitting any ASTs.
// All outstanding markers are invalidated after this call.
// Panics if EOF already emitted.
//
func (p *Parser) Discard() {
	// Nothing can be discarded after EOF emitted
	//
	if p.eofOut {
		panic("Parser.Discard: No discards allowed after EOF is emitted")
	}
	p.consume()
}

// TODO Remove this after API settles down
// DFARRELL - Not needed as Parser ensure `HasNext()` before calling ParserFn and
//            `HasNext()` / `CanPeek()` are better tools to use from withing ParserFn.
// // EOF returns true if the peek buffer is empty AND the input has reached EOF.
// // This will return false if there are any tokens remaining in the peek buffer.
// // See CanPeek() / HasNext() to confirm if tokens are available to peek / consume.
// //
// func (p *Parser) EOF() bool {
// 	return p.eof && p.cache.Len() == p.matchLen
// }

// newParser
//
func newParser(tokens token.Nexter, start ParserFn) *Parser {
	return &Parser{
		tokens:    tokens,
		cache:     list.New(),
		matchTail: nil,
		matchLen:  0,
		nextFn:    start,
		emits:     list.New(),
		eof:       false,
		eofOut:    false,
		markerId:  0,
	}
}

// growPeek tries to ensure the peek buffer has Len() >= n, growing if needed, returning success or failure.
// n is 1-based.
//
func (p *Parser) growPeek(n int) bool {
	// Grow to n
	//
	peekLen := p.cache.Len() - p.matchLen
	for peekLen < n {
		// Nothing to do if EOF reached already
		//
		if p.eof {
			return false
		}
		// Fetch next token from input
		//
		token, err := p.tokens.Next()
		// Process any returned token, regardless of er
		//
		if token != nil {
			p.cache.PushBack(token)
			peekLen++
		}
		// If there was an error, process it now
		//
		if err != nil {
			switch err {
			// EOF Error
			//
			case io.EOF:
				p.eof = true

			// NON-EOF Error
			//
			default:
				// For lack of a better plan, treat as EOF for now
				// TODO Think about how to handle non-EOF errors.
				// TODO Log error.
				// TODO Expose upstream?
				//
				p.eof = true
			}
		}
	}
	return true
}

// peekHead computes the peek buffer head as a function of the matchTail.
//
func (p *Parser) peekHead() *list.Element {
	// If any consumed tokens
	//
	if p.matchLen > 0 {
		// Peek buffer starts after consumed tokens
		//
		// assert(p.matchTail != nil)
		return p.matchTail.Next()
	}
	// Its ALL the peek buffer
	//
	return p.cache.Front()
}

// emit Emits an AST.
// Panics if EOF already emitted.
//
func (p *Parser) emit(ast interface{}) {
	// Nothing can be emitted after EOF
	// NOTE: This check is a fail-safe and will likely never hit as all public methods check/panic explicitly.
	//
	if p.eofOut {
		panic("Parser: No further emits allowed after EOF is emitted")
	}
	// If emitting EOF
	//
	if ast == nil {
		// Clear the peek buffer, discarding consumed tokens
		//
		p.matchTail = nil
		p.matchLen = 0
		p.cache.Init()
		// Invalidate outstanding markers manually,
		// avoiding otherwise redundant call to consume()
		//
		p.markerId++ // TODO If it ever takes 2+ commands to invalidate markers, then turn into separate method.
		// Mark EOF
		//
		p.eof = true
		p.eofOut = true
		// Emit EOF marker
		//
		p.emits.PushBack(nil)
	} else {
		p.consume()

		p.emits.PushBack(ast)
	}
}

// consume consumes the matched tokens.
// All outstanding markers are invalidated after this call.
//
func (p *Parser) consume() {
	// Discard tokens
	//
	for p.matchLen > 0 {
		p.cache.Remove(p.cache.Front())
		p.matchLen--
	}
	// Invalidate outstanding markers
	//
	p.markerId++ // Invalidate outstanding markers
}
