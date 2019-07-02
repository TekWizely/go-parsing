module github.com/tekwizely/go-parsing/parser

go 1.12

require (
	// To update:
	//
	// $ go get github.com/tekwizely/go-parsing/lexer@master
	//
	github.com/tekwizely/go-parsing/lexer v0.0.0-20190702211947-89bd4df59bc0
	github.com/tekwizely/go-parsing/lexer/token v0.0.0-20190622183031-974f82a44df9
)

// For Local testing against changes that aren't upstream
//
//replace github.com/tekwizely/go-parsing/lexer => ../lexer

//replace github.com/tekwizely/go-parsing/lexer/token => ../lexer/token
