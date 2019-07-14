module github.com/tekwizely/go-parsing/parser

go 1.12

require (
	// To update:
	//
	// $ go get github.com/tekwizely/go-parsing/lexer@master
	//
	github.com/tekwizely/go-parsing/lexer v0.0.0-20190714043513-9514494dd58a
	github.com/tekwizely/go-parsing/lexer/token v0.0.0-20190714025745-8a1a69651c50
)

// For Local testing against changes that aren't upstream
//
//replace github.com/tekwizely/go-parsing/lexer => ../lexer

//replace github.com/tekwizely/go-parsing/lexer/token => ../lexer/token
