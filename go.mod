module github.com/tekwizely/go-parsing

go 1.12

require (
	github.com/tekwizely/go-parsing/lexer v0.0.0 // indirect
	github.com/tekwizely/go-parsing/lexer/token v0.0.0
	github.com/tekwizely/go-parsing/parser v0.0.0 // indirect
)

replace (
	github.com/tekwizely/go-parsing/lexer => ./lexer
	github.com/tekwizely/go-parsing/lexer/token => ./lexer/token
	github.com/tekwizely/go-parsing/parser => ./parser
)
