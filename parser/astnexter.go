package parser

import "io"

// ASTNexter is returned by the Parse function and provides a means of retrieving ASTs emitted from the parser.
//
type ASTNexter interface {

	// Next tries to fetch the next available AST, returning an error if something goes wrong.
	// Will return io.EOF to indicate end-of-file.
	// An error other than io.EOF may be recoverable and does not necessarily indicate end-of-file.
	// Even when an error is present, the returned AST may still be valid and should be checked.
	// Once io.EOF is returned, any further calls will continue to return io.EOF.
	//
	Next() (interface{}, error)
}

// astNexter is the internal structure that backs the parser's ASTNexter.
//
type astNexter struct {
	parser *Parser
	next   interface{}
	eof    bool
}

// Next implements ASTNexter.Next().
// We build on the previous HasNext/Next impl to keep changes minimal.
//
func (e *astNexter) Next() (interface{}, error) {
	if !e.hasNext() {
		return nil, io.EOF
	}
	tok := e.next
	e.next = nil
	return tok, nil
}

// hasNext Initiates calls to Parser.Fn functions and is the primary entry point for retrieving ASTs from the parser.
//
func (e *astNexter) hasNext() bool {
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
	for e.parser.output.Len() == 0 {
		// Anyone to call?
		// Any tokens to scan?
		//
		if e.parser.nextFn != nil && e.parser.CanPeek(1) {
			e.parser.nextFn = e.parser.nextFn(e.parser)
		} else
		// Parser Terminated, let's clean up.
		// If EOF was never emitted, then emit it now.
		//
		if e.parser.eofOut == false {
			e.parser.EmitEOF()
		}
	}
	// Consume the AST.
	// We'll either cache it or discard it.
	//
	emit := e.parser.output.Remove(e.parser.output.Front())
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
