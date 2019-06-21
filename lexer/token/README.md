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
// Nexter provides a means of retrieving tokens (and errors) emitted from the lexer.
//
type Nexter interface {

	// Next tries to fetch the next available token, returning an error if something goes wrong.
	// Will return io.EOF to indicate end-of-file.
	// An error other than io.EOF may be recoverable and does not necessarily indicate end-of-file.
	// Even when an error is present, the returned token may still be valid and should be checked.
	// Once io.EOF is returned, any further calls will continue to return io.EOF.
	//
	Next(Token, error)
}
```

## License

The `go-parsing` repo and all contained packages are released under the [MIT](https://opensource.org/licenses/MIT) License.  See `LICENSE` file.
