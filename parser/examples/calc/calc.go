package main

//
//	Input is read from STDIN
//
//	The input expression is matched against the following pattern:
//
//	input_exp:
//	( id '=' )? general_exp
//	general_exp:
//		operand ( operator operand )?
//	operand:
//		number | id | '(' general_exp ')'
//	operator:
//		'+' | '-' | '*' | '/'
//	number:
//		digit+ ( '.' digit+ )?
//	digit:
//		['0'..'9']
//	id:
//		alpha ( alpha | digit )*
//	alpha:
//		['a'..'z'] | ['A'..'Z']
//
//	Precedence is as expected, with '*' and '/' have higher precedence
//	than '+' and '-', as follows:
//
//	1 + 2 * 3 - 4 / 5  ==  1 + (2 * 3) - (4 / 5)
//

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/tekwizely/go-parsing/lexer"
	"github.com/tekwizely/go-parsing/lexer/token"
	"github.com/tekwizely/go-parsing/parser"
)

// We define our lexer tokens starting from the pre-defined EOF token
//
const (
	TId token.Type = lexer.TStart + iota
	TNumber
	TPlus
	TMinus
	TMultiply
	TDivide
	TEquals
	TOpenParen
	TCloseParen
)

// To store variables
//
var vars = map[string]float64{}

// Single-character tokens
//
var singleChars = []byte{'+', '-', '*', '/', '=', '(', ')'}

var singleTokens = []token.Type{TPlus, TMinus, TMultiply, TDivide, TEquals, TOpenParen, TCloseParen}

// main
//
func main() {
	// Create a buffered reader from STDIN
	//
	stdin := bufio.NewReader(os.Stdin)

	// Read each line of input
	//
	for input, _, err := stdin.ReadLine(); err == nil; input, _, err = stdin.ReadLine() {
		// Anything to process?
		//
		if len(input) > 0 {
			// Create a new lexer to turn the input text into tokens
			//
			tokens := lexer.LexBytes(input, lex)

			// Create a new parser that feeds off the lexer and generates expression values
			//
			values := parser.Parse(tokens, parse)

			// Loop over parser emits
			//
			for value, parseErr := values.Next(); parseErr == nil; value, parseErr = values.Next() {
				fmt.Printf("%v\n", value)
			}
		}
	}
}

// lex is the starting (and only) StateFn for lexing the input into tokens
//
func lex(l *lexer.Lexer) lexer.Fn {

	// Single-char token?
	//
	if i := bytes.IndexRune(singleChars, l.Peek(1)); i >= 0 {
		l.Next()                    // Match the rune
		l.EmitType(singleTokens[i]) // Emit just the type, discarding the matched rune
		return lex
	}

	switch {

	// Skip whitespace
	//
	case tryMatchWhitespace(l):
		l.Clear()

	// Number
	//
	case tryMatchNumber(l):
		l.EmitToken(TNumber)

	// ID
	//
	case tryMatchID(l):
		l.EmitToken(TId)

	// Unknown
	//
	default:
		r := l.Next()
		l.Clear()
		fmt.Printf("Unknown Character: '%c'\n", r)
	}

	// See you again soon!
	return lex
}

// tryMatchWhitespace
//
func tryMatchWhitespace(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); r == ' ' || r == '\t' {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchRune
//
func tryMatchRune(l *lexer.Lexer, r rune) bool {
	if l.CanPeek(1) {
		if p := l.Peek(1); r == p {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchDigit
//
func tryMatchDigit(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); r >= '0' && r <= '9' {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchAlpha
//
func tryMatchAlpha(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchAlphaNum
//
func tryMatchAlphaNum(l *lexer.Lexer) bool {
	if l.CanPeek(1) {
		if r := l.Peek(1); (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			l.Next()
			return true
		}
	}
	return false
}

// tryMatchNumber [0-9]+ ( . [0-9]+ )?
//
func tryMatchNumber(l *lexer.Lexer) bool {
	if tryMatchDigit(l) {
		for tryMatchDigit(l) {
			// Nothing to do, rune already matched
		}
		if m := l.Marker(); tryMatchRune(l, '.') && tryMatchDigit(l) {
			for tryMatchDigit(l) {

			}
		} else {
			m.Apply()
		}
		return true
	}
	return false
}

// tryMatchID [a-zA-Z] [0-9a-zA-Z]*
//
func tryMatchID(l *lexer.Lexer) bool {
	if tryMatchAlpha(l) {
		for tryMatchAlphaNum(l) {
			// Nothing to do
		}
		return true
	}
	return false
}

// parse tries to parse an expression from the lexed tokens.
// Delegates to either parseEvaluation or parseAssignment.
//
func parse(p *parser.Parser) parser.Fn {

	switch {

	// Assignment
	//
	case p.CanPeek(3) && p.PeekType(1) == TId && p.PeekType(2) == TEquals:
		return parseAssignment

	// Evaluation
	//
	default:
		return parseEvaluation
	}
}

// parseAssignment evaluates an expression and stores the results in the specified variable.
// The assignment will be in the form [ ID '=' expression ].
// Assumes "ID '='" has been peek-matched by root parser.
//
func parseAssignment(p *parser.Parser) parser.Fn {
	tID := p.Next()
	p.Next() // Skip '='
	if value, err := parseGeneralExpression(p); err == nil {
		// Should be at end of input
		//
		if !p.CanPeek(1) {
			vars[tID.Value()] = value
		} else {
			fmt.Println("Expecting Operator")
		}
	} else {
		fmt.Println(err.Error())
	}
	return nil // One pass
}

// parseEvaluation parses a general experssion and emits the computed result.
//
func parseEvaluation(p *parser.Parser) parser.Fn {
	if value, err := parseGeneralExpression(p); err == nil {
		// Should be at end of input
		//
		if !p.CanPeek(1) {
			p.Emit(value)
		} else {
			fmt.Println("Expecting Operator")
		}
	} else {
		fmt.Println(err.Error())
	}
	return nil // One pass
}

// parseGeneralExpression is the starting point for parsing a General Expression.
// It is basically a pass-through to parseAdditiveExpression, but it feels cleaner.
//
func parseGeneralExpression(p *parser.Parser) (f float64, err error) {
	return parseAdditiveExpression(p)
}

// parseAdditiveExpression parses [ expression ( ( '+' | '-' ) expression )? ].
//
func parseAdditiveExpression(p *parser.Parser) (f float64, err error) {

	var a float64
	if f, err = parseMultiplicitiveExpression(p); err == nil && p.CanPeek(1) {

		switch p.PeekType(1) {

		// Add (+)
		//
		case TPlus:
			p.Next() // Skip '+'
			if a, err = parseAdditiveExpression(p); err == nil {
				f += a
			}

		// Subtract (-)
		//
		case TMinus:
			p.Next() // Skip '-'
			if a, err = parseAdditiveExpression(p); err == nil {
				f -= a
			}
		}
	}

	return
}

// parseMultiplicitiveExpression parses [ expression ( ( '*' | '/' ) expression )? ].
//
func parseMultiplicitiveExpression(p *parser.Parser) (f float64, err error) {

	var m float64
	if f, err = parseOperand(p); err == nil && p.CanPeek(1) {

		switch p.PeekType(1) {

		// Multiply (*)
		//
		case TMultiply:
			p.Next() // Skip '*'
			if m, err = parseMultiplicitiveExpression(p); err == nil {
				f *= m
			}

		// Divide (/)
		//
		case TDivide:
			p.Next() // Skip '/'
			if m, err = parseMultiplicitiveExpression(p); err == nil {
				f /= m
			}
		}
	}

	return
}

// parseOperand parses [ id | number | '(' expression ')' ].
//
func parseOperand(p *parser.Parser) (f float64, err error) {

	// EOF
	//
	if !p.CanPeek(1) {
		return 0, errors.New("Unexpected EOF - Expecting operand")
	}

	m := p.Marker()

	switch p.PeekType(1) {

	// ID
	//
	case TId:
		var id = p.Next().Value()
		var ok bool
		if f, ok = vars[id]; !ok {
			err = fmt.Errorf("id '%s' not defined", id)
		}

	// Number
	//
	case TNumber:
		n := p.Next().Value()
		if f, err = strconv.ParseFloat(n, 64); err != nil {
			fmt.Printf("Error parsing number '%s': %s", n, err.Error())
		}

	// '(' Expresson ')'
	//
	case TOpenParen:
		p.Next() // Skip '('
		if f, err = parseGeneralExpression(p); err == nil {
			if p.CanPeek(1) && p.PeekType(1) == TCloseParen {
				p.Next() // Skip ')'
			} else {
				err = errors.New("Unbalanced Paren")
			}
		}

	// Unknown
	//
	default:
		err = errors.New("Expecting operand")
	}

	if err != nil {
		m.Apply()
	}

	return
}
