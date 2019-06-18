module github.com/tekwizely/go-parsing/parser

go 1.12

// To update:
//
// $ go get github.com/tekwizely/go-parsing/lexer@master
//
require github.com/tekwizely/go-parsing/lexer v0.0.0-20190618060741-bb9c88748c57

// For Local testing against changes that aren't upstream
//
//replace github.com/tekwizely/go-parsing/lexer => ../lexer
