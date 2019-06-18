package parser

// Emits is returned by the Parse function and provides to retrieve emitted ASTs from the parser.
// Implements a basic iterator pattern with HasNext() and Next() methods.
//
type Emits struct {
	parser *Parser
	next   interface{}
	eof    bool
}

// Next Retrieves the next AST from the parser.
// See HasNext() to determine if any ASTs are available.
// Panics if HasNext() returns false.
//
func (e *Emits) Next() interface{} {
	// We double check for saved next to maybe avoid the call
	//
	if e.next == nil && e.HasNext() == false {
		panic("Emits.Next: No AST available")
	}
	tok := e.next
	e.next = nil
	return tok
}

// HasNext confirms if there are ASTs available.
// This method initiates calls to ParserFn functions and is the primary entry point for retrieving ASTs from the parser.
// If it returns true, you can safely call Next() to retrieve the next AST.
// If it returns false, EOF has been reached and calling Next() will generate a panic.
//
func (e *Emits) HasNext() bool {
	// If AST previously fetched, return now
	//
	if e.next != nil {
		return true
	}
	// Nothing to do once EOF reached
	//
	if e.eof {
		return false
	}
	// If no ASTs available, try to fetch some.
	//
	for e.parser.emits.Len() == 0 {
		// Anyone to call?
		// Any tokens to scan?
		//
		if e.parser.nextFn != nil && e.parser.HasNext() {
			e.parser.nextFn = e.parser.nextFn(e.parser)
		} else {
			// Parser Terminated, let's clean up.
			// If EOF was never emitted, then emit it now.
			//
			if e.parser.eofOut == false {
				e.parser.EmitEOF()
			}
		}
	}
	// Consume the AST.
	// We'll either cache it or discard it.
	//
	emit := e.parser.emits.Remove(e.parser.emits.Front())
	// Is if EOF?
	//
	if emit == nil {
		// Mark EOF, discarding the AST
		//
		e.eof = true
		return false
	}
	// Store the AST for pickup
	//
	e.next = emit
	return true
}
