# braceexpansion

[![Build Status](https://travis-ci.org/thomasheller/braceexpansion.svg?branch=master)](https://travis-ci.org/thomasheller/braceexpansion)
[![Go Report Card](https://goreportcard.com/badge/github.com/thomasheller/braceexpansion)](https://goreportcard.com/report/github.com/thomasheller/braceexpansion)
[![Coverage Status](https://coveralls.io/repos/github/thomasheller/braceexpansion/badge.svg?branch=master)](https://coveralls.io/github/thomasheller/braceexpansion?branch=master)

Shell brace expansion implemented in Go (golang).

Supports some specialties required by
[multigoogle](https://github.com/thomasheller/multigoogle).

Numeric ranges are currently not supported, as I didn't need them.
Feel free to send a PR.

## Build

```sh
$ cd cmd && go build -o be
```

## Usage (command line):

```sh
$ be '{a,b}{1,2}'
a1
a2
b1
b2
```

## Usage (library):

```go
import (
	be "github.com/thomasheller/braceexpansion"
	"fmt"
)

func main() {
	tree, err := be.New().Parse("{a,b}{1,2}")
	if err != nil {
		panic(err)
	}
	for _, s := range tree.Expand() {
		fmt.Println(s)
	}
	
	// Output:
	// a1
	// a2
	// b1
	// b2
}
```
