# braceexpansion

Shell brace expansion implemented in Go (golang).

Supports some specialties required by
[multigoogle](https://github.com/thomasheller/multigoogle).

Numeric ranges are currently not supported, as I didn't need them.
Feel free to send a PR.

# Usage:

```go
import (
	be "github.com/thomasheller/braceexpansion"
	"fmt"
)

func main() {
	tree, err := be.New.Parse("{a,b}{1,2}")
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
