package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/tekwizely/go-parsing/lexer"
)

// Usage : wordcount <filename>
//
func usage() {
	fmt.Printf("usage: %s <filename>\n", os.Args[0])
}

// We define our lexer tokens starting from the pre-defined START token
//
const (
	tSpace = lexer.TStart + iota
	tNewline
	tWord
)

// We will attempt to match 3 newline styles: [ "\n", "\r", "\r\n" ]
//
const (
	runeNewLine = '\n'
	runeReturn  = '\r'
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	var (
		file *os.File
		err  error
	)

	//  Open the file, panic on error
	//
	if file, err = os.Open(os.Args[1]); err != nil {
		panic(err)
	}

	// int inits to 0
	//
	var (
		chars  int
		words  int
		spaces int
		lines  int
	)

	// To help us track last line in file (which may not have a newline)
	//
	var emptyLine = true

	tokens := lexer.LexReader(file, lexerFn)

	// Process lexer-emitted tokens
	//
	for t, lexErr := tokens.Next(); lexErr == nil; t, lexErr = tokens.Next() { // We only emit EOF so !nil should do it
		chars += len(t.Value())

		switch t.Type() {
		case tWord:
			words++
			emptyLine = false

		case tNewline:
			lines++
			spaces += len(t.Value())
			emptyLine = true

		case tSpace:
			spaces += len(t.Value())
			emptyLine = false

		default:
			panic("Unreachable")
		}
	}

	// If last line not empty, up line count
	//
	if !emptyLine {
		lines++
	}

	fmt.Printf("%d words, %d spaces, %d lines, %d chars\n", words, spaces, lines, chars)
}

func lexerFn(l *lexer.Lexer) lexer.Fn {

	// Can skip canPeek() check on first rune, per lexer rules
	//
	switch r := l.Peek(1); {

	// Newline '\n'
	// We check this before Space to avoid hit from unicode.IsSpace() check
	//
	case r == runeNewLine:
		l.Next()
		l.EmitToken(tNewline)

	// Return '\r', optionally followed by newLine '\n'
	// We check this before Space to avoid hit from unicode.IsSpace() check
	//
	case r == runeReturn:
		l.Next()
		if l.CanPeek(1) && l.Peek(1) == runeNewLine {
			l.Next()
		}
		l.EmitToken(tNewline)

	// Space or Word
	//
	default:
		isSpace := unicode.IsSpace(r)
		// Match verified rune to avoid re-check
		//
		l.Next()
		// Match further consecutive runes of same type
		//
		for l.CanPeek(1) && unicode.IsSpace(l.Peek(1)) == isSpace {
			l.Next()
		}
		// Emit token
		//
		if isSpace {
			l.EmitToken(tSpace)
		} else {
			l.EmitToken(tWord)
		}
	}

	return lexerFn // Let's do it again
}
