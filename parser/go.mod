module github.com/tekwizely/go-parsing/parser

go 1.12

// To update:
//
// $ go get github.com/tekwizely/go-parsing/lexer@master
//
require github.com/tekwizely/go-parsing/lexer v0.0.0-20190620203355-f59e82d2158e

// For Local testing against changes that aren't upstream
//
//replace github.com/tekwizely/go-parsing/lexer => ../lexer
