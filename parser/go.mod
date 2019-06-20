module github.com/tekwizely/go-parsing/parser

go 1.12

// To update:
//
// $ go get github.com/tekwizely/go-parsing/lexer@master
//
require github.com/tekwizely/go-parsing/lexer v0.0.0-20190620064451-d1cd406dab0b

// For Local testing against changes that aren't upstream
//
//replace github.com/tekwizely/go-parsing/lexer => ../lexer
