package lexer

import "container/list"

// Marker snapshots the state of the lexer to allow rewinding.
//
// See the following lexer functions for creating and user markers:
//
//  - Lexer.Marker()
//  - Marker.Valid()
//  - marker.Apply()
//
type Marker struct {
	lexer     *Lexer
	markerID  int
	matchTail *list.Element
	matchLen  int
	nextFn    Fn
}

// Marker returns a marker that you can use to reset the lexer to a previous state.
// A marker is good up until the next Emit() or Clear() action.
// Use Marker.Valid() to verify that a marker is still valid before using it.
// Use Marker.Apply() to reset the lexer state to the marker position.
//
func (l *Lexer) Marker() *Marker {
	return &Marker{lexer: l, markerID: l.markerID, matchTail: l.matchTail, matchLen: l.matchLen, nextFn: l.nextFn}
}

// Valid confirms if the marker is still valid.
// If Valid returns true, you can safely reset the lexer state to the marker position via Marker.Apply().
//
func (m *Marker) Valid() bool {
	// ALL markers invalid once EOF emitted
	//
	return !m.lexer.eofOut && m.markerID == m.lexer.markerID
}

// Apply resets the lexer state to the marker position.
// Returns the Lexer.Fn that was stored at the time the marker was created.
// Use `return marker.Apply()` to tell the lexer to forward to the marked function.
// Use Valid() to verify that a marker is still valid before using it.
// It is safe to apply a marker multiple times, as long as it passes Valid().
// Panics if marker fails Valid() check.
//
func (m *Marker) Apply() Fn {
	if m.Valid() == false {
		panic("Invalid marker")
	}
	m.lexer.matchTail = m.matchTail
	m.lexer.matchLen = m.matchLen
	return m.nextFn
}
