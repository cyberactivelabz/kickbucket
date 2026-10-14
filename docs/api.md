# Go library

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cyberactivelabz/kickbucket/pkg/kickbucket"
)

func main() {
	results, err := kickbucket.Scan(context.Background(), "acme", kickbucket.Options{})
	if err != nil {
		log.Fatalf("scanning: %v", err)
	}
	for _, r := range results {
		fmt.Println(r.Provider, r.Name, r.Access)
	}
}
```

::: warning
This shows the intended API shape. Update it to match the real `pkg/kickbucket` package once it exists.
:::
