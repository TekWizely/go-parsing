module github.com/tekwizely/go-parsing/parser

go 1.12

require (
	// To update:
	//
	// $ go get github.com/tekwizely/go-parsing/lexer@master
	//
	github.com/tekwizely/go-parsing/lexer v0.0.0-20190621065023-8ea7298baf63
	github.com/tekwizely/go-parsing/lexer/token v0.0.0-20190621065023-8ea7298baf63
)

// For Local testing against changes that aren't upstream
//
replace github.com/tekwizely/go-parsing/lexer => ../lexer

replace github.com/tekwizely/go-parsing/lexer/token => ../lexer/token
