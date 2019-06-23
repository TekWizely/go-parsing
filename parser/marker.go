package parser

import "container/list"

// Marker snapshots the state of the parser to allow rewinding.
//
// See the following parser functions for creating and user markers:
//
//  - Parser.Marker()
//  - Parser.CanReset()
//  - Parser.Reset()
//
type Marker struct {
	markerID  int
	matchTail *list.Element
	matchLen  int
	nextFn    ParserFn
}

// Marker returns a marker that you can use to reset the parser to a previous state.
// A marker is good up until the next Emit() or Discard() action.
// Use CanReset() to verify that a marker is still valid before using it.
// Use Reset() to reset the parser state to the marker position.
//
func (p *Parser) Marker() *Marker {
	return &Marker{markerID: p.markerID, matchTail: p.matchTail, matchLen: p.matchLen, nextFn: p.nextFn}
}

// CanReset confirms if the marker is still valid.
// If CanReset returns true, you can safely reset the parser state to the marker position.
//
func (p *Parser) CanReset(m *Marker) bool {
	// ALL markers invalid once EOF emitted
	//
	return !p.eofOut && m.markerID == p.markerID
}

// Reset resets the parser state to the marker position.
// Returns the ParserFn that was stored at the time the marker was created.
// Use `return marker.Reset()` to tell the parser to forward to the marked function.
// Use CanReset() to verify that a marker is still valid before using it.
// It is safe to reset a marker multiple times, as long as it passes CanReset().
// Panics if marker fails CanReset() check.
//
func (p *Parser) Reset(m *Marker) ParserFn {
	if p.CanReset(m) == false {
		panic("Invalid marker")
	}
	p.matchTail = m.matchTail
	p.matchLen = m.matchLen
	return p.nextFn
}
