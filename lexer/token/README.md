# go-parsing / lexer / token
[![GoDoc](https://godoc.org/github.com/tekwizely/go-parsing/lexer/token?status.svg)](https://godoc.org/github.com/tekwizely/go-parsing/lexer/token)
[![MIT license](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/tekwizely/go-parsing/blob/master/LICENSE)

# Overview

Token-related types and interfaces used between the lexer and the parser.

## Using

#### Importing

```go
import "github.com/tekwizely/go-parsing/lexer/token"
```

### token.Token

```go
// Token captures the type code + optional text string emitted from the lexer.
//
type Token interface {

	// Type returns the type code of the token.
	//
	Type() Type

	// Value returns the matched rune(s) that represent the token value.
	// Can be the empty string.
	//
	Value() string
}
```

### token.Type

```go
// Type identifies the type code of tokens emitted from the lexer.
//
type Type int
```

### token.Nexter

```go
// Nexter provides methods to retrieve tokens emitted from the lexer.
// Implements a basic iterator pattern with HasNext() and Next() methods.
//
type Nexter interface {

	// HasNext confirms if there are tokens available.
	// If it returns true, you can safely call Next() to retrieve the next token.
	// If it returns false, EOF has been reached and calling Next() will generate a panic.
	//
	HasNext() bool

	// Next Retrieves the next token from the lexer.
	// See HasNext() to determine if any tokens are available.
	// Panics if HasNext() returns false.
	//
	Next() Token
}
```

## License

The `go-parsing` repo and all contained packages are released under the [MIT](https://opensource.org/licenses/MIT) License.  See `LICENSE` file.
