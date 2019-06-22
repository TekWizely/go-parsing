module github.com/tekwizely/go-parsing/lexer

go 1.12

// To update:
//
// $ go get github.com/tekwizely/go-parsing/lexer/token@master
//
require github.com/tekwizely/go-parsing/lexer/token v0.0.0-20190621000622-442f3491df4a

// For Local testing against changes that aren't upstream
//
replace github.com/tekwizely/go-parsing/lexer/token => ./token
