## Textquery

---

Textquery is a lucene-like query parser and matcher. It uses `strings.Contains`
for each query component matching.

### Example

```go
package main

import (
	"fmt"
	"github.com/denisdubovitskiy/textquery-go"
)

func main() {
	tree := textquery.Parse("((a AND b) AND NOT c) AND (d AND e)")

	fmt.Println(tree.Match("a b d e")) // true
	fmt.Println(tree.Match("a b c d e")) // false
}
```