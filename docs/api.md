# Go library

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cyberactivelabz/kickbucket/pkg/bucketmap"
)

func main() {
	results, err := bucketmap.Scan(context.Background(), "acme", bucketmap.Options{})
	if err != nil {
		log.Fatalf("scanning: %v", err)
	}
	for _, r := range results {
		fmt.Println(r.Provider, r.Name, r.Access)
	}
}
```

::: warning
This shows the intended API shape. Update it to match the real `pkg/bucketmap` package once it exists.
:::
