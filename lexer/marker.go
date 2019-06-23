package lexer

import "container/list"

// Marker snapshots the state of the lexer to allow rewinding.
//
// See the following lexer functions for creating and user markers:
//
//  - Lexer.Marker()
//  - Lexer.CanReset()
//  - Lexer.Reset()
//
type Marker struct {
	markerID  int
	matchTail *list.Element
	matchLen  int
	nextFn    LexerFn
}

// Marker returns a marker that you can use to reset the lexer to a previous state.
// A marker is good up until the next Emit() or Clear() action.
// Use CanReset() to verify that a marker is still valid before using it.
// Use Reset() to reset the lexer state to the marker position.
//
func (l *Lexer) Marker() *Marker {
	return &Marker{markerID: l.markerID, matchTail: l.matchTail, matchLen: l.matchLen, nextFn: l.nextFn}
}

// CanReset confirms if the marker is still valid.
// If CanReset returns true, you can safely reset the lexer state to the marker position.
//
func (l *Lexer) CanReset(m *Marker) bool {
	// ALL markers invalid once EOF emitted
	//
	return !l.eofOut && m.markerID == l.markerID
}

// Reset resets the lexer state to the marker position.
// Returns the LexerFn that was stored at the time the marker was created.
// Use `return marker.Reset()` to tell the lexer to forward to the marked function.
// Use CanReset() to verify that a marker is still valid before using it.
// It is safe to reset a marker multiple times, as long as it passes CanReset().
// Panics if marker fails CanReset() check.
//
func (l *Lexer) Reset(m *Marker) LexerFn {
	if l.CanReset(m) == false {
		panic("Invalid marker")
	}
	l.matchTail = m.matchTail
	l.matchLen = m.matchLen
	return l.nextFn
}
