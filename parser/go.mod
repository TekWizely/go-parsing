module github.com/tekwizely/go-parsing/parser

go 1.12

// To update:
//
// $ go get github.com/tekwizely/go-parsing/lexer@master
//
require github.com/tekwizely/go-parsing/lexer v0.0.0-20190617061751-164d4ff03e0d

// For Local testing against changes that aren't upstream
//
//replace github.com/tekwizely/go-parsing/lexer => ../lexer
