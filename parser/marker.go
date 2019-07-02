package parser

import "container/list"

// Marker snapshots the state of the parser to allow rewinding.
//
// See the following parser functions for creating and user markers:
//
//  - Parser.Marker()
//  - Marker.Valid()
//  - Marker.Apply()
//
type Marker struct {
	parser    *Parser
	markerID  int
	matchTail *list.Element
	matchLen  int
	nextFn    Fn
}

// Marker returns a marker that you can use to reset the parser to a previous state.
// A marker is good up until the next Emit() or Clear() action.
// Use Marker.Valid() to verify that a marker is still valid before using it.
// Use Marker.Apply() to reset the parser state to the marker position.
//
func (p *Parser) Marker() *Marker {
	return &Marker{parser: p, markerID: p.markerID, matchTail: p.matchTail, matchLen: p.matchLen, nextFn: p.nextFn}
}

// Valid confirms if the marker is still valid.
// If Valid returns true, you can safely reset the parser state to the marker position via Marker.Apply().
//
func (m *Marker) Valid() bool {
	// ALL markers invalid once EOF emitted
	//
	return !m.parser.eofOut && m.markerID == m.parser.markerID
}

// Apply resets the parser state to the marker position.
// Returns the Parser.Fn that was stored at the time the marker was created.
// Use `return marker.Apply()` to tell the parser to forward to the marked function.
// Use Valid() to verify that a marker is still valid before using it.
// It is safe to apply a marker multiple times, as long as it passes Valid().
// Panics if marker fails Valid check.
//
func (m *Marker) Apply() Fn {
	if !m.Valid() {
		panic("Invalid marker")
	}
	m.parser.matchTail = m.matchTail
	m.parser.matchLen = m.matchLen
	return m.nextFn
}
